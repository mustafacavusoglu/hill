package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap tüm keyboard kısayollarını tanımlar
type KeyMap struct {
	// Global
	Quit      key.Binding
	Help      key.Binding
	FocusNext key.Binding
	FocusPrev key.Binding

	// Panel odağı
	PanelRequest  key.Binding
	PanelResponse key.Binding
	PanelHistory  key.Binding

	// Request panel
	Send         key.Binding
	ChangeMethod key.Binding
	AddHeader    key.Binding
	NextField    key.Binding
	PrevField    key.Binding

	// Response panel
	CopyBody   key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding

	// History
	SelectEntry key.Binding
	DeleteEntry key.Binding
}

// DefaultKeyMap varsayılan kısayolları döner
var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "çıkış"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "yardım"),
	),
	FocusNext: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "sonraki alan"),
	),
	FocusPrev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "önceki alan"),
	),
	PanelRequest: key.NewBinding(
		key.WithKeys("f1"),
		key.WithHelp("F1", "istek"),
	),
	PanelResponse: key.NewBinding(
		key.WithKeys("f2"),
		key.WithHelp("F2", "yanıt"),
	),
	PanelHistory: key.NewBinding(
		key.WithKeys("f3"),
		key.WithHelp("F3", "geçmiş"),
	),
	Send: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "gönder"),
	),
	ChangeMethod: key.NewBinding(
		key.WithKeys("ctrl+m"),
		key.WithHelp("ctrl+m", "method değiştir"),
	),
	AddHeader: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "header ekle"),
	),
	NextField: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "sonraki alan"),
	),
	PrevField: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "önceki alan"),
	),
	CopyBody: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "body kopyala"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "yukarı"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "aşağı"),
	),
	SelectEntry: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "yükle"),
	),
	DeleteEntry: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "sil"),
	),
}
