package panels

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mustafacavusoglu/hill/internal/httpclient"
)

// ResponseModel yanıt paneli state'ini yönetir
type ResponseModel struct {
	response *httpclient.Response
	viewport viewport.Model
	width    int
	height   int
}

func NewResponseModel() ResponseModel {
	vp := viewport.New(60, 20)
	return ResponseModel{
		viewport: vp,
	}
}

func (m ResponseModel) SetSize(w, h int) ResponseModel {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 8
	return m
}

func (m ResponseModel) SetResponse(resp *httpclient.Response) ResponseModel {
	m.response = resp
	if resp != nil {
		content := buildResponseContent(resp)
		m.viewport.SetContent(content)
		m.viewport.GotoTop()
	}
	return m
}

func (m ResponseModel) Update(msg tea.Msg) (ResponseModel, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ResponseModel) View() string {
	if m.response == nil {
		placeholder := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true).
			Render("İstek göndermek için ctrl+r'ye basın...")
		return placeholder
	}

	var sb strings.Builder

	// Status satırı
	status := m.response.Status
	statusStyle := lipgloss.NewStyle().Bold(true)
	switch {
	case m.response.StatusCode >= 500:
		statusStyle = statusStyle.Foreground(lipgloss.Color("1"))
	case m.response.StatusCode >= 400:
		statusStyle = statusStyle.Foreground(lipgloss.Color("3"))
	case m.response.StatusCode >= 300:
		statusStyle = statusStyle.Foreground(lipgloss.Color("14"))
	default:
		statusStyle = statusStyle.Foreground(lipgloss.Color("2"))
	}

	meta := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(
		fmt.Sprintf("(%s)  %s  %d bytes",
			m.response.Duration.Round(1000000),
			m.response.Proto,
			m.response.BodySize,
		),
	)
	sb.WriteString(statusStyle.Render(status))
	sb.WriteString("  ")
	sb.WriteString(meta)
	sb.WriteString("\n\n")
	sb.WriteString(m.viewport.View())
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true).Render(
		"[j/k] Scroll  [c] Kopyala",
	))

	return sb.String()
}

func buildResponseContent(resp *httpclient.Response) string {
	var sb strings.Builder

	// Headers
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	for k, vals := range resp.Headers {
		sb.WriteString(headerStyle.Render(k))
		sb.WriteString(": ")
		sb.WriteString(strings.Join(vals, ", "))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Body
	contentType := resp.Headers.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		sb.WriteString(httpclient.PrettyJSON(resp.Body))
	} else {
		sb.WriteString(string(resp.Body))
	}

	return sb.String()
}

func (m ResponseModel) Response() *httpclient.Response {
	return m.response
}
