package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mustafacavusoglu/hill/internal/httpclient"
	"github.com/mustafacavusoglu/hill/internal/tui/panels"
)

type activePanel int

const (
	panelRequest activePanel = iota
	panelResponse
	panelHistory
)

// Mesaj tipleri
type responseMsg struct {
	response *httpclient.Response
	err      error
}

type loadHistoryMsg struct {
	entry panels.HistoryEntry
}

// Model root bubbletea model
type Model struct {
	width  int
	height int
	active activePanel
	keys   KeyMap

	request  panels.RequestModel
	response panels.ResponseModel
	history  panels.HistoryModel

	spinner spinner.Model
	loading bool
	err     error
}

func NewModel() Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	return Model{
		active:   panelRequest,
		keys:     DefaultKeyMap,
		request:  panels.NewRequestModel(),
		response: panels.NewResponseModel(),
		history:  panels.NewHistoryModel(),
		spinner:  sp,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateSizes()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			// body editörü dışında çıkış
			if m.active != panelRequest {
				return m, tea.Quit
			}
		case "f1":
			m.active = panelRequest
			m.history = m.history.Blur()
		case "f2":
			m.active = panelResponse
			m.history = m.history.Blur()
		case "f3":
			m.active = panelHistory
			m.history = m.history.Focus()
		case "ctrl+r":
			if !m.loading {
				req := m.request.ToRequest()
				if req.URL != "" {
					m.loading = true
					cmds = append(cmds, sendRequest(req))
					cmds = append(cmds, m.spinner.Tick)
				}
			}
		}

		// Panel mesajlarını ilet
		if m.active == panelRequest {
			var cmd tea.Cmd
			m.request, cmd = m.request.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.active == panelResponse {
			var cmd tea.Cmd
			m.response, cmd = m.response.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.active == panelHistory {
			switch msg.String() {
			case "enter":
				if entry := m.history.Selected(); entry != nil {
					m.request.SetURL(entry.Request.URL)
					m.request.SetMethod(entry.Request.Method)
					body := ""
					if len(entry.Request.Body) > 0 {
						body = string(entry.Request.Body)
					}
					m.request.SetBody(body)
					m.active = panelRequest
					m.history = m.history.Blur()
				}
			default:
				var cmd tea.Cmd
				m.history, cmd = m.history.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case responseMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.response = m.response.SetResponse(msg.response)
			// History'e ekle
			entry := panels.HistoryEntry{
				Request:   m.request.ToRequest(),
				Response:  msg.response,
				Timestamp: time.Now(),
			}
			m.history = m.history.Add(entry)
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Yükleniyor..."
	}

	// Header / tab bar
	header := m.renderHeader()

	// Üst alan yüksekliği (%70)
	topHeight := int(float64(m.height-5) * 0.70)
	if topHeight < 10 {
		topHeight = 10
	}
	// Alt alan (history)
	bottomHeight := m.height - topHeight - 5
	if bottomHeight < 4 {
		bottomHeight = 4
	}

	halfW := m.width / 2

	// Request panel
	reqStyle := StyleInactivePanel
	if m.active == panelRequest {
		reqStyle = StyleActivePanel
	}
	reqContent := m.request.View()
	reqPanel := reqStyle.
		Width(halfW - 2).
		Height(topHeight).
		Render(reqContent)

	// Response panel
	respStyle := StyleInactivePanel
	if m.active == panelResponse {
		respStyle = StyleActivePanel
	}
	respContent := m.response.View()
	if m.loading {
		respContent = fmt.Sprintf("%s  İstek gönderiliyor...", m.spinner.View())
	} else if m.err != nil {
		respContent = StyleStatusErr.Render("Hata: " + m.err.Error())
	}
	respPanel := respStyle.
		Width(m.width - halfW - 2).
		Height(topHeight).
		Render(respContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, reqPanel, respPanel)

	// History panel
	histStyle := StyleInactivePanel
	if m.active == panelHistory {
		histStyle = StyleActivePanel
	}
	histTitle := StyleDim.Render("HISTORY")
	histContent := histTitle + "\n" + m.history.View()
	histPanel := histStyle.
		Width(m.width - 2).
		Height(bottomHeight).
		Render(histContent)

	// Footer
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		topRow,
		histPanel,
		footer,
	)
}

func (m Model) renderHeader() string {
	title := StyleTitle.Padding(0, 1).Render("hill")

	tab := func(label string, panel activePanel, key string) string {
		s := fmt.Sprintf("[%s] %s", key, label)
		if m.active == panel {
			return StyleTabActive.Render(s)
		}
		return StyleTabInactive.Render(s)
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Left,
		title,
		tab("Request", panelRequest, "F1"),
		tab("Response", panelResponse, "F2"),
		tab("History", panelHistory, "F3"),
		StyleTabInactive.Render("[?] Yardım"),
	)

	return lipgloss.NewStyle().
		Width(m.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("8")).
		Render(tabs)
}

func (m Model) renderFooter() string {
	help := []string{
		"[ctrl+r] Gönder",
		"[ctrl+m] Method",
		"[F1-F3] Panel",
		"[q] Çıkış",
	}
	helpStr := StyleHelp.Render(strings.Join(help, "  "))
	return lipgloss.NewStyle().
		Width(m.width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(lipgloss.Color("8")).
		Render(helpStr)
}

func (m Model) updateSizes() Model {
	topHeight := int(float64(m.height-5) * 0.70)
	if topHeight < 10 {
		topHeight = 10
	}
	bottomHeight := m.height - topHeight - 5
	if bottomHeight < 4 {
		bottomHeight = 4
	}
	halfW := m.width / 2

	m.request = m.request.SetSize(halfW-4, topHeight-2)
	m.response = m.response.SetSize(m.width-halfW-4, topHeight-2)
	m.history = m.history.SetSize(m.width-4, bottomHeight-2)
	return m
}

// sendRequest HTTP isteğini goroutine'de çalıştırır
func sendRequest(req httpclient.Request) tea.Cmd {
	return func() tea.Msg {
		resp, err := httpclient.Execute(req)
		return responseMsg{response: resp, err: err}
	}
}

// Start TUI uygulamasını başlatır
func Start() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
