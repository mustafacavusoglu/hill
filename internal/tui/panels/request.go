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
	fieldMethod requestField = iota
	fieldURL
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
		focused:   fieldMethod, // method seçici ile başla
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
		case "tab":
			switch m.focused {
			case fieldMethod:
				m.urlInput.Focus()
				m.focused = fieldURL
			case fieldURL:
				m.urlInput.Blur()
				m.bodyInput.Focus()
				m.focused = fieldBody
			case fieldBody:
				m.bodyInput.Blur()
				m.focused = fieldMethod
			}
			return m, nil

		case "shift+tab":
			switch m.focused {
			case fieldMethod:
				m.urlInput.Blur()
				m.bodyInput.Focus()
				m.focused = fieldBody
			case fieldURL:
				m.urlInput.Blur()
				m.focused = fieldMethod
			case fieldBody:
				m.bodyInput.Blur()
				m.urlInput.Focus()
				m.focused = fieldURL
			}
			return m, nil
		}

		// Method seçici odaklanmışken: ok tuşları veya space ile döngü
		if m.focused == fieldMethod {
			switch msg.String() {
			case "right", "l", " ":
				m.methodIdx = (m.methodIdx + 1) % len(methods)
				return m, nil
			case "left", "h":
				m.methodIdx = (m.methodIdx - 1 + len(methods)) % len(methods)
				return m, nil
			case "enter":
				// enter ile URL alanına geç
				m.urlInput.Focus()
				m.focused = fieldURL
				return m, nil
			}
			return m, nil
		}
	}

	if m.focused == fieldURL {
		var cmd tea.Cmd
		m.urlInput, cmd = m.urlInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focused == fieldBody {
		var cmd tea.Cmd
		m.bodyInput, cmd = m.bodyInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RequestModel) View() string {
	var sb strings.Builder

	// Method seçici
	methodStr := methods[m.methodIdx]
	methodColor := lipgloss.Color("2") // GET=yeşil
	switch methodStr {
	case "POST":
		methodColor = lipgloss.Color("12")
	case "PUT":
		methodColor = lipgloss.Color("3")
	case "DELETE":
		methodColor = lipgloss.Color("1")
	case "PATCH":
		methodColor = lipgloss.Color("14")
	case "HEAD":
		methodColor = lipgloss.Color("8")
	}

	var methodWidget string
	if m.focused == fieldMethod {
		// Odaklanmış: ok işaretleriyle çerçeve
		methodWidget = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(methodColor).
			Padding(0, 1).
			Render("◀ " + methodStr + " ▶")
	} else {
		methodWidget = lipgloss.NewStyle().
			Bold(true).
			Foreground(methodColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1).
			Render(methodStr)
	}

	urlLine := lipgloss.JoinHorizontal(lipgloss.Left,
		methodWidget,
		"  ",
		m.urlInput.View(),
	)
	sb.WriteString(urlLine)
	sb.WriteString("\n\n")

	// Body
	bodyLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("Body (JSON):")
	if m.focused == fieldBody {
		bodyLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("Body (JSON):")
	}
	sb.WriteString(bodyLabel)
	sb.WriteString("\n")
	sb.WriteString(m.bodyInput.View())
	sb.WriteString("\n\n")

	// Kısayollar — aktif alana göre ipucu
	var helpText string
	switch m.focused {
	case fieldMethod:
		helpText = "[←/→] Method değiştir  [enter/tab] URL'e geç"
	case fieldURL:
		helpText = "[ctrl+r] Gönder  [tab] Body'e geç  [shift+tab] Method'a geç"
	case fieldBody:
		helpText = "[ctrl+r] Gönder  [shift+tab] URL'e geç"
	}
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true).Render(helpText))

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
