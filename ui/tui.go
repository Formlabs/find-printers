package ui

import (
	"find-printers/borg"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.design/x/clipboard"
)

const (
	columnKeyName     = "serial"
	columnKeyType     = "type"
	columnKeyIP       = "ip"
	columnKeyFirmware = "firmware"
)

const (
	sortDirectionAsc  = "asc"
	sortDirectionDesc = "desc"
)

var (
	styleHighlight = lipgloss.NewStyle().Foreground(lipgloss.Color("#5EC6D4"))
)

type Model struct {
	scrollableTable table.Model
	filterTextInput textinput.Model
	columnSortKey   string
	sortDirection   string
	copyBuffer      string
}

func NewModel() Model {
	return Model{
		scrollableTable: table.New([]table.Column{
			table.NewColumn(columnKeyName, "SERIAL", 30).WithFiltered(true),
			table.NewColumn(columnKeyType, "TYPE", 20).WithFiltered(true),
			table.NewColumn(columnKeyIP, "IP", 20).WithFiltered(true),
			table.NewColumn(columnKeyFirmware, "FIRMWARE", 80).WithFiltered(true),
		}).
			Filtered(true).
			WithHorizontalFreezeColumnCount(1).
			WithFooterVisibility(false).
			WithPageSize(40).
			Focused(true).
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("#337D87")).
					Foreground(lipgloss.Color("#5EC6D4")).
					Align(lipgloss.Left),
			).WithMultiline(true),
		filterTextInput: textinput.New(),
		columnSortKey:   columnKeyName,
		sortDirection:   sortDirectionAsc,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) AddDevices(devices []borg.Device) tea.Model {
	rows := []table.Row{}
	for _, device := range devices {
		row := table.NewRow(table.RowData{
			columnKeyName:     device.Serial,
			columnKeyType:     device.MachineTypeID,
			columnKeyIP:       device.IPAddress,
			columnKeyFirmware: device.FirmwareVersion,
		})
		rows = append(rows, row)
	}

	m.scrollableTable = m.scrollableTable.WithRows(rows)
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// global
		if msg.String() == "ctrl+c" {
			cmds = append(cmds, tea.Quit)

			return m, tea.Batch(cmds...)
		}
		// event to filter
		if m.filterTextInput.Focused() {
			if msg.String() == "enter" {
				m.filterTextInput.Blur()
			} else {
				m.filterTextInput, _ = m.filterTextInput.Update(msg)
			}
			m.scrollableTable = m.scrollableTable.WithFilterInput(m.filterTextInput)

			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "/":
			m.filterTextInput.Focus()

		case "q":
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)

		// Sort commands
		case "S":
			m.columnSortKey = columnKeyName
			model := resortBySelectedColumn(m)
			return model, nil
		case "T":
			m.columnSortKey = columnKeyType
			model := resortBySelectedColumn(m)
			return model, nil
		case "I":
			m.columnSortKey = columnKeyIP
			model := resortBySelectedColumn(m)
			return model, nil
		case "F":
			m.columnSortKey = columnKeyFirmware
			model := resortBySelectedColumn(m)
			return model, nil

		// Copy commands
		case "s":
			copydata := m.scrollableTable.HighlightedRow().Data[columnKeyName]
			m.copyBuffer = copydata.(string)
			clipboard.Write(clipboard.FmtText, []byte(copydata.(string)))
		case "t":
			copydata := m.scrollableTable.HighlightedRow().Data[columnKeyType]
			m.copyBuffer = copydata.(string)
			clipboard.Write(clipboard.FmtText, []byte(copydata.(string)))
		case "i":
			copydata := m.scrollableTable.HighlightedRow().Data[columnKeyIP]
			m.copyBuffer = copydata.(string)
			clipboard.Write(clipboard.FmtText, []byte(copydata.(string)))
		case "f":
			copydata := m.scrollableTable.HighlightedRow().Data[columnKeyFirmware]
			m.copyBuffer = copydata.(string)
			clipboard.Write(clipboard.FmtText, []byte(copydata.(string)))

		default:
			m.scrollableTable, cmd = m.scrollableTable.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func resortBySelectedColumn(m Model) Model {
	switch m.sortDirection {
	case sortDirectionDesc:
		m.sortDirection = sortDirectionAsc
		m.scrollableTable = m.scrollableTable.SortByAsc(m.columnSortKey)
		return m
	case sortDirectionAsc:
		m.sortDirection = sortDirectionDesc
		m.scrollableTable = m.scrollableTable.SortByDesc(m.columnSortKey)
		return m
	}
	return m
}

func sortDirectionToAscii(d string) string {
	if d == sortDirectionAsc {
		return "↓"
	} else {
		return "↑"
	}
}

func (m Model) View() string {
	body := strings.Builder{}

	h := func(s string) string { return styleHighlight.Render(s) }
	body.WriteString("Press " + h("q") + " or " + h("ctrl+c") + " to quit")
	body.WriteString("\nPress " + h("/") + " to start filtering")
	body.WriteString("\nSort by (" + h("S") + ")erial, (" + h("T") + ")ype, (" + h("I") + ")p or (" + h("F") + ")irmware	| Sorted by " + h(m.columnSortKey) + " " + h(sortDirectionToAscii(m.sortDirection)))
	body.WriteString("\nPress (" + h("s") + ")erial, (" + h("t") + ")ype, (" + h("i") + ")p or (" + h("f") + ")irmware	| Buffer: " + h(m.copyBuffer))

	body.WriteString("\n")

	body.WriteString(m.filterTextInput.View() + "\n")
	body.WriteString(m.scrollableTable.View())

	return body.String()
}

func Run(m Model) {
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
