package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("12")  // mavi
	ColorSuccess   = lipgloss.Color("2")   // yeşil
	ColorWarning   = lipgloss.Color("3")   // sarı
	ColorError     = lipgloss.Color("1")   // kırmızı
	ColorDim       = lipgloss.Color("8")   // gri
	ColorHighlight = lipgloss.Color("14")  // açık mavi
	ColorWhite     = lipgloss.Color("15")  // beyaz
	ColorBg        = lipgloss.Color("0")   // siyah

	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleActivePanel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary)

	StyleInactivePanel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorDim)

	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Underline(true).
			Padding(0, 1)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorDim).
				Padding(0, 1)

	StyleStatusOK = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	StyleStatusErr = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError)

	StyleStatusWarn = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorHighlight)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleBold = lipgloss.NewStyle().
			Bold(true)

	StyleHelp = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true)

	StyleMethod = map[string]lipgloss.Style{
		"GET":    lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess),
		"POST":   lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary),
		"PUT":    lipgloss.NewStyle().Bold(true).Foreground(ColorWarning),
		"DELETE": lipgloss.NewStyle().Bold(true).Foreground(ColorError),
		"PATCH":  lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight),
	}
)

func MethodStyle(method string) lipgloss.Style {
	if s, ok := StyleMethod[method]; ok {
		return s
	}
	return StyleBold
}
