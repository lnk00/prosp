package tui

import (
	"log"
	"strconv"

	"github.com/lnk00/prosp/db"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lnk00/prosp/models"
	"github.com/pkg/browser"
)

type model struct {
	keys  keyMap
	help  help.Model
	table table.Model
	db    db.Database
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) UpdatePrevJobStatus() {
	idx := m.table.Cursor()
	rows := m.table.Rows()
	rows[idx][3] = string(models.JobStatus(rows[idx][3]).GetPrevStatus())
	m.db.UpdateJobStatus(rows[idx][4], models.JobStatus(rows[idx][3]))
	m.table.SetRows(rows)
}

func (m *model) UpdateNextJobStatus() {
	idx := m.table.Cursor()
	rows := m.table.Rows()
	rows[idx][3] = string(models.JobStatus(rows[idx][3]).GetNextStatus())
	m.db.UpdateJobStatus(rows[idx][4], models.JobStatus(rows[idx][3]))
	m.table.SetRows(rows)
}

func (m model) OpenUrl() {
	idx := m.table.Cursor()
	rows := m.table.Rows()
	browser.OpenURL(rows[idx][4])
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Left):
			m.UpdatePrevJobStatus()
		case key.Matches(msg, m.keys.Right):
			m.UpdateNextJobStatus()
		case key.Matches(msg, m.keys.Return):
			m.OpenUrl()
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
		{Title: "Status", Width: 20},
		{Title: "Link", Width: 50},
	}

	for id, job := range jobs {
		rows = append(rows, table.Row{
			strconv.Itoa(id),
			job.Title,
			job.Location,
			string(job.Status),
			job.Link,
		})
	}

	return columns, rows
}

func Render(db db.Database) {
	columns, rows := buildTable(db.GetJobs())

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(40),
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
	m := model{keys, help.New(), t, db}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		log.Fatalf("failed to render: %v", err)
	}

}
