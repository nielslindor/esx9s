package tui

import (
	"context"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nielslindor/esx9s/internal/domain"
)

type viewID int

const (
	hostsView viewID = iota
	vmsView
	datastoresView
	eventsView
)

type viewSpec struct {
	id       viewID
	key      string
	title    string
	subtitle string
	columns  []string
	rows     [][]string
}

var defaultViews = []viewSpec{
	{
		id:       hostsView,
		key:      "1",
		title:    "Hosts",
		subtitle: "ESXi host inventory placeholder",
		columns:  []string{"Name", "State", "CPU", "Memory"},
		rows: [][]string{
			{"esx-lab-01", "connected", "24 cores", "128 GiB"},
			{"esx-lab-02", "maintenance", "16 cores", "96 GiB"},
			{"esx-edge-01", "connected", "12 cores", "64 GiB"},
		},
	},
	{
		id:       vmsView,
		key:      "2",
		title:    "VMs",
		subtitle: "Virtual machine list placeholder",
		columns:  []string{"Name", "Power", "Host", "IP"},
		rows: [][]string{
			{"vcsa", "on", "esx-lab-01", "10.0.20.10"},
			{"dns-01", "on", "esx-lab-02", "10.0.20.53"},
			{"backup-proxy", "off", "esx-edge-01", "-"},
		},
	},
	{
		id:       datastoresView,
		key:      "3",
		title:    "Datastores",
		subtitle: "Storage capacity placeholder",
		columns:  []string{"Name", "Type", "Free", "Capacity"},
		rows: [][]string{
			{"shared-nvme", "VMFS", "5.8 TiB", "8.0 TiB"},
			{"iso-library", "NFS", "820 GiB", "1.0 TiB"},
			{"vsan-default", "vSAN", "12.4 TiB", "18.0 TiB"},
		},
	},
	{
		id:       eventsView,
		key:      "4",
		title:    "Events/Audit",
		subtitle: "Recent activity placeholder",
		columns:  []string{"Time", "Severity", "Actor", "Message"},
		rows: [][]string{
			{"10:42", "info", "system", "Host inventory refresh queued"},
			{"10:39", "warn", "admin", "Snapshot policy drift detected"},
			{"10:31", "info", "operator", "Datastore view opened"},
		},
	},
}

type Styles struct {
	app      lipgloss.Style
	title    lipgloss.Style
	tab      lipgloss.Style
	active   lipgloss.Style
	panel    lipgloss.Style
	header   lipgloss.Style
	row      lipgloss.Style
	selected lipgloss.Style
	muted    lipgloss.Style
	footer   lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		app:      lipgloss.NewStyle().Padding(1, 2),
		title:    lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")),
		tab:      lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("8")),
		active:   lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("10")),
		panel:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).Padding(1, 2),
		header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14")),
		row:      lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		selected: lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("12")),
		muted:    lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		footer:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
}

type Model struct {
	width    int
	height   int
	active   int
	cursor   map[viewID]int
	views    []viewSpec
	styles   Styles
	mode     string
	status   string
	showHelp bool
	quitting bool
}

func NewModel() Model {
	return Model{
		active: 0,
		views:  defaultViews,
		cursor: map[viewID]int{
			hostsView:      0,
			vmsView:        0,
			datastoresView: 0,
			eventsView:     0,
		},
		styles: defaultStyles(),
		mode:   "mock",
		status: "Ready",
	}
}

func NewModelFromInventory(inventory domain.Inventory) Model {
	model := NewModel()
	model.views = viewsFromInventory(inventory)
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "right", "l", "tab":
			m.active = (m.active + 1) % len(m.views)
			m.status = fmt.Sprintf("Showing %s", m.ActiveView())
		case "left", "shift+tab":
			m.active = (m.active + len(m.views) - 1) % len(m.views)
			m.status = fmt.Sprintf("Showing %s", m.ActiveView())
		case "down", "j":
			m.moveCursor(1)
		case "up", "k":
			m.moveCursor(-1)
		case "1", "2", "3", "4", "h", "v", "d", "t":
			m.selectView(msg.String())
			m.status = fmt.Sprintf("Showing %s", m.ActiveView())
		case "?":
			m.showHelp = !m.showHelp
		case "r":
			m.status = fmt.Sprintf("Refreshed %s from %s provider", m.ActiveView(), m.mode)
		case "enter":
			m.status = fmt.Sprintf("Details placeholder for %s selection", m.ActiveView())
		case "p":
			if m.views[m.active].id == vmsView {
				m.status = "VM power action placeholder: typed confirmation and audit are required"
			}
		case "s":
			if m.views[m.active].id == vmsView {
				m.status = "VM snapshot action placeholder: confirmation and audit are required"
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "esx9s closed\n"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderTitle(),
		m.renderTabs(),
		m.renderPanel(),
		m.renderHelp(),
		m.renderFooter(),
	)

	return m.styles.app.Render(content)
}

func (m Model) ActiveView() string {
	return m.views[m.active].title
}

func (m *Model) moveCursor(delta int) {
	view := m.views[m.active]
	next := m.cursor[view.id] + delta
	if next < 0 {
		next = 0
	}
	if next >= len(view.rows) {
		next = len(view.rows) - 1
	}
	m.cursor[view.id] = next
}

func (m *Model) selectView(key string) {
	aliases := map[string]string{
		"h": "1",
		"v": "2",
		"d": "3",
		"t": "4",
	}
	if alias, ok := aliases[key]; ok {
		key = alias
	}

	for i, view := range m.views {
		if view.key == key {
			m.active = i
			return
		}
	}
}

func (m Model) renderTitle() string {
	title := m.styles.title.Render("esx9s")
	status := m.styles.muted.Render(fmt.Sprintf("mode:%s  active:%s", m.mode, m.ActiveView()))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", status)
}

func (m Model) renderTabs() string {
	tabs := make([]string, 0, len(m.views))
	for i, view := range m.views {
		label := fmt.Sprintf("%s %s", view.key, view.title)
		if i == m.active {
			tabs = append(tabs, m.styles.active.Render(label))
			continue
		}
		tabs = append(tabs, m.styles.tab.Render(label))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m Model) renderPanel() string {
	view := m.views[m.active]
	body := []string{
		m.styles.header.Render(view.title),
		m.styles.muted.Render(view.subtitle),
		"",
		m.renderTable(view),
	}

	panelWidth := m.width - 6
	if panelWidth < 64 {
		panelWidth = 64
	}

	return m.styles.panel.Width(panelWidth).Render(strings.Join(body, "\n"))
}

func (m Model) renderTable(view viewSpec) string {
	widths := columnWidths(view)
	lines := []string{m.styles.header.Render(formatRow(view.columns, widths))}

	for i, row := range view.rows {
		prefix := "  "
		style := m.styles.row
		if i == m.cursor[view.id] {
			prefix = "> "
			style = m.styles.selected
		}
		lines = append(lines, style.Render(prefix+formatRow(row, widths)))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderFooter() string {
	lines := []string{
		m.styles.footer.Render("1/h hosts  2/v vms  3/d datastores  4/t tasks  ? help  r refresh  enter details  q quit"),
		m.styles.muted.Render("status: " + m.status),
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderHelp() string {
	if !m.showHelp {
		return ""
	}

	help := []string{
		"Keyboard",
		"1/h Hosts   2/v VMs   3/d Datastores   4/t Tasks/Events",
		"j/k Move selection   r Refresh   enter Details placeholder",
		"p VM power placeholder   s VM snapshot placeholder",
		"All destructive actions require confirmation and audit before real provider execution.",
	}

	return m.styles.panel.Width(64).Render(strings.Join(help, "\n"))
}

func columnWidths(view viewSpec) []int {
	widths := make([]int, len(view.columns))
	for i, column := range view.columns {
		widths[i] = len(column)
	}
	for _, row := range view.rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	for i := range widths {
		widths[i] += 2
	}
	return widths
}

func formatRow(cells []string, widths []int) string {
	parts := make([]string, len(cells))
	for i, cell := range cells {
		parts[i] = fmt.Sprintf("%-*s", widths[i], cell)
	}
	return strings.TrimRight(strings.Join(parts, ""), " ")
}

func viewsFromInventory(inventory domain.Inventory) []viewSpec {
	return []viewSpec{
		{
			id:       hostsView,
			key:      "1",
			title:    "Hosts",
			subtitle: "Mock ESXi host inventory",
			columns:  []string{"Name", "State", "CPU", "Memory", "VMs"},
			rows:     hostRows(inventory.Hosts),
		},
		{
			id:       vmsView,
			key:      "2",
			title:    "VMs",
			subtitle: "Mock virtual machine inventory",
			columns:  []string{"Name", "Power", "Host", "CPU", "Memory", "IP"},
			rows:     vmRows(inventory.VMs),
		},
		{
			id:       datastoresView,
			key:      "3",
			title:    "Datastores",
			subtitle: "Mock storage capacity",
			columns:  []string{"Name", "Host", "Type", "Free", "Capacity"},
			rows:     datastoreRows(inventory.Datastores),
		},
		{
			id:       eventsView,
			key:      "4",
			title:    "Events/Audit",
			subtitle: "Mock recent provider tasks",
			columns:  []string{"Started", "Status", "Host", "Operation", "Target"},
			rows:     taskRows(inventory.Tasks),
		},
	}
}

func hostRows(hosts []domain.Host) [][]string {
	if len(hosts) == 0 {
		return [][]string{{"-", "no hosts", "-", "-", "-"}}
	}

	rows := make([][]string, 0, len(hosts))
	for _, host := range hosts {
		rows = append(rows, []string{
			host.Name,
			string(host.Status),
			fmt.Sprintf("%s/%s", formatMHz(host.CPUUsedMHz), formatMHz(host.CPUCapacityMHz)),
			fmt.Sprintf("%s/%s", formatBytes(host.MemoryUsedBytes), formatBytes(host.MemoryBytes)),
			fmt.Sprintf("%d", host.VMCount),
		})
	}
	return rows
}

func vmRows(vms []domain.VM) [][]string {
	if len(vms) == 0 {
		return [][]string{{"-", "no VMs", "-", "-", "-", "-"}}
	}

	rows := make([][]string, 0, len(vms))
	for _, vm := range vms {
		ipAddress := vm.IPAddress
		if ipAddress == "" {
			ipAddress = "-"
		}
		rows = append(rows, []string{
			vm.Name,
			string(vm.PowerState),
			vm.HostName,
			fmt.Sprintf("%d", vm.CPUCount),
			formatBytes(vm.MemoryBytes),
			ipAddress,
		})
	}
	return rows
}

func datastoreRows(datastores []domain.Datastore) [][]string {
	if len(datastores) == 0 {
		return [][]string{{"-", "-", "no datastores", "-", "-"}}
	}

	rows := make([][]string, 0, len(datastores))
	for _, datastore := range datastores {
		rows = append(rows, []string{
			datastore.Name,
			datastore.HostName,
			datastore.Type,
			formatBytes(datastore.FreeBytes),
			formatBytes(datastore.CapacityBytes),
		})
	}
	return rows
}

func taskRows(tasks []domain.Task) [][]string {
	if len(tasks) == 0 {
		return [][]string{{"-", "no tasks", "-", "-", "-"}}
	}

	rows := make([][]string, 0, len(tasks))
	for _, task := range tasks {
		target := task.TargetName
		if target == "" {
			target = "-"
		}
		rows = append(rows, []string{
			task.StartedAt.Format("15:04"),
			string(task.Status),
			task.HostName,
			task.Operation,
			target,
		})
	}
	return rows
}

func formatMHz(mhz int64) string {
	if mhz >= 1000 {
		return fmt.Sprintf("%.1f GHz", float64(mhz)/1000)
	}
	return fmt.Sprintf("%d MHz", mhz)
}

func formatBytes(bytes int64) string {
	const gib = float64(1024 * 1024 * 1024)
	if bytes == 0 {
		return "0 GiB"
	}
	return fmt.Sprintf("%.0f GiB", float64(bytes)/gib)
}

func Run(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
	return RunModel(ctx, stdin, stdout, NewModel())
}

func RunModel(ctx context.Context, stdin io.Reader, stdout io.Writer, model Model) error {
	program := tea.NewProgram(
		model,
		tea.WithInput(stdin),
		tea.WithOutput(stdout),
		tea.WithContext(ctx),
	)
	_, err := program.Run()
	return err
}
