package panels

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mustafacavusoglu/hill/internal/httpclient"
)

// HistoryEntry geçmişte yapılmış bir isteği temsil eder
type HistoryEntry struct {
	Request   httpclient.Request
	Response  *httpclient.Response
	Timestamp time.Time
}

// HistoryModel geçmiş panel state'ini yönetir
type HistoryModel struct {
	entries []HistoryEntry
	table   table.Model
	width   int
	height  int
}

func NewHistoryModel() HistoryModel {
	cols := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "URL", Width: 50},
		{Title: "Status", Width: 8},
		{Title: "Süre", Width: 10},
		{Title: "Zaman", Width: 20},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("12")).
		Bold(true)
	t.SetStyles(s)

	return HistoryModel{table: t}
}

func (m HistoryModel) SetSize(w, h int) HistoryModel {
	m.width = w
	m.height = h
	m.table.SetHeight(h - 4)

	cols := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "URL", Width: w - 60},
		{Title: "Status", Width: 8},
		{Title: "Süre", Width: 10},
		{Title: "Zaman", Width: 20},
	}
	m.table.SetColumns(cols)
	return m
}

func (m HistoryModel) Add(entry HistoryEntry) HistoryModel {
	m.entries = append([]HistoryEntry{entry}, m.entries...)
	m.table.SetRows(buildRows(m.entries))
	return m
}

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m HistoryModel) View() string {
	if len(m.entries) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true).
			Render("  Henüz istek gönderilmedi...")
	}
	return m.table.View()
}

func (m HistoryModel) Selected() *HistoryEntry {
	if len(m.entries) == 0 {
		return nil
	}
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(m.entries) {
		return nil
	}
	entry := m.entries[idx]
	return &entry
}

func (m HistoryModel) Focus() HistoryModel {
	m.table.Focus()
	return m
}

func (m HistoryModel) Blur() HistoryModel {
	m.table.Blur()
	return m
}

func buildRows(entries []HistoryEntry) []table.Row {
	rows := make([]table.Row, len(entries))
	for i, e := range entries {
		status := "—"
		duration := "—"
		if e.Response != nil {
			status = fmt.Sprintf("%d", e.Response.StatusCode)
			duration = e.Response.Duration.Round(1000000).String()
		}
		url := e.Request.URL
		if len(url) > 60 {
			url = url[:57] + "..."
		}
		rows[i] = table.Row{
			e.Request.Method,
			url,
			status,
			duration,
			e.Timestamp.Format("15:04:05"),
		}
	}
	return rows
}
