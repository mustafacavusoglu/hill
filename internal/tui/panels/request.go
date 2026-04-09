package panels

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mustafacavusoglu/hill/internal/httpclient"
)

var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}

type requestField int

const (
	fieldURL requestField = iota
	fieldBody
)

type HeaderRow struct {
	Key   textinput.Model
	Value textinput.Model
}

// RequestModel istek paneli state'ini yönetir
type RequestModel struct {
	methodIdx int
	urlInput  textinput.Model
	bodyInput textarea.Model
	headers   []HeaderRow
	focused   requestField
	width     int
	height    int
}

func NewRequestModel() RequestModel {
	url := textinput.New()
	url.Placeholder = "https://api.example.com/endpoint"
	url.Focus()
	url.Width = 60

	body := textarea.New()
	body.Placeholder = "{\n  \"key\": \"value\"\n}"
	body.ShowLineNumbers = false
	body.SetWidth(60)
	body.SetHeight(8)

	return RequestModel{
		methodIdx: 0,
		urlInput:  url,
		bodyInput: body,
		focused:   fieldURL,
	}
}

func (m RequestModel) SetSize(w, h int) RequestModel {
	m.width = w
	m.height = h
	m.urlInput.Width = w - 20
	m.bodyInput.SetWidth(w - 4)
	m.bodyInput.SetHeight(h - 12)
	return m
}

func (m RequestModel) Update(msg tea.Msg) (RequestModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+m":
			m.methodIdx = (m.methodIdx + 1) % len(methods)
			return m, nil

		case "tab":
			if m.focused == fieldURL {
				m.urlInput.Blur()
				m.bodyInput.Focus()
				m.focused = fieldBody
			} else {
				m.bodyInput.Blur()
				m.urlInput.Focus()
				m.focused = fieldURL
			}
			return m, nil

		case "shift+tab":
			if m.focused == fieldBody {
				m.bodyInput.Blur()
				m.urlInput.Focus()
				m.focused = fieldURL
			} else {
				m.urlInput.Blur()
				m.bodyInput.Focus()
				m.focused = fieldBody
			}
			return m, nil
		}
	}

	if m.focused == fieldURL {
		var cmd tea.Cmd
		m.urlInput, cmd = m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.bodyInput, cmd = m.bodyInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RequestModel) View() string {
	var sb strings.Builder

	// Method + URL satırı
	methodStr := methods[m.methodIdx]
	methodStyle := lipgloss.NewStyle().Bold(true).Padding(0, 1)
	switch methodStr {
	case "GET":
		methodStyle = methodStyle.Foreground(lipgloss.Color("2"))
	case "POST":
		methodStyle = methodStyle.Foreground(lipgloss.Color("12"))
	case "PUT":
		methodStyle = methodStyle.Foreground(lipgloss.Color("3"))
	case "DELETE":
		methodStyle = methodStyle.Foreground(lipgloss.Color("1"))
	case "PATCH":
		methodStyle = methodStyle.Foreground(lipgloss.Color("14"))
	}

	urlLine := lipgloss.JoinHorizontal(lipgloss.Left,
		methodStyle.Render(methodStr),
		"  ",
		m.urlInput.View(),
	)
	sb.WriteString(urlLine)
	sb.WriteString("\n\n")

	// Body
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("Body (JSON):"))
	sb.WriteString("\n")
	sb.WriteString(m.bodyInput.View())
	sb.WriteString("\n\n")

	// Kısayollar
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true).Render(
		"[ctrl+r] Gönder  [ctrl+m] Method  [tab] Alan geçiş",
	)
	sb.WriteString(help)

	return sb.String()
}

func (m RequestModel) ToRequest() httpclient.Request {
	headers := make(map[string]string)
	for _, h := range m.headers {
		k := h.Key.Value()
		v := h.Value.Value()
		if k != "" {
			headers[k] = v
		}
	}

	return httpclient.Request{
		Method:  methods[m.methodIdx],
		URL:     m.urlInput.Value(),
		Headers: headers,
		Body:    []byte(m.bodyInput.Value()),
	}
}

func (m *RequestModel) SetURL(url string) {
	m.urlInput.SetValue(url)
}

func (m *RequestModel) SetMethod(method string) {
	for i, meth := range methods {
		if meth == method {
			m.methodIdx = i
			return
		}
	}
}

func (m *RequestModel) SetBody(body string) {
	m.bodyInput.SetValue(body)
}

func (m RequestModel) URL() string {
	return m.urlInput.Value()
}

func (m RequestModel) FocusURL() RequestModel {
	m.bodyInput.Blur()
	m.urlInput.Focus()
	m.focused = fieldURL
	return m
}
