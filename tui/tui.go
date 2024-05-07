package tui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lnk00/prosp/models"
)

type model struct {
	keys  keyMap
	help  help.Model
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	tableView := baseStyle.Render(m.table.View())
	helpView := m.help.View(m.keys)
	return tableView + "\n" + helpView + "\n"
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func buildTable(jobs []models.Job) ([]table.Column, []table.Row) {
	rows := []table.Row{}
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Title", Width: 50},
		{Title: "Location", Width: 50},
		{Title: "Link", Width: 50},
		{Title: "Status", Width: 20},
	}

	for id, job := range jobs {
		rows = append(rows, table.Row{
			strconv.Itoa(id),
			job.Title,
			job.Location,
			job.Link,
			string(job.Status),
		})
	}

	return columns, rows
}

func Render(jobs []models.Job) {
	columns, rows := buildTable(jobs)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m := model{keys, help.New(), t}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		log.Fatalf("failed to render: %v", err)
	}

}
