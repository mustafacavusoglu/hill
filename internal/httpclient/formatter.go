package httpclient

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleSuccess = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))   // yeşil
	styleRedirect = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))  // sarı
	styleClientErr = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1")) // kırmızı
	styleServerErr = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")) // parlak kırmızı
	styleHeader    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))            // mavi
	styleMeta      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))             // gri
)

func PrintResponse(resp *Response) {
	// Status satırı
	statusStr := resp.Status
	var styledStatus string
	switch {
	case resp.StatusCode >= 500:
		styledStatus = styleServerErr.Render(statusStr)
	case resp.StatusCode >= 400:
		styledStatus = styleClientErr.Render(statusStr)
	case resp.StatusCode >= 300:
		styledStatus = styleRedirect.Render(statusStr)
	default:
		styledStatus = styleSuccess.Render(statusStr)
	}

	meta := styleMeta.Render(fmt.Sprintf("(%s)  %s  %d bytes", resp.Duration.Round(1000000), resp.Proto, resp.BodySize))
	fmt.Printf("%s  %s\n", styledStatus, meta)

	// Headers
	fmt.Println()
	for k, vals := range resp.Headers {
		fmt.Printf("%s: %s\n", styleHeader.Render(k), strings.Join(vals, ", "))
	}

	// Body
	fmt.Println()
	contentType := resp.Headers.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		fmt.Println(prettyJSON(resp.Body))
	} else {
		fmt.Println(string(resp.Body))
	}
}

func prettyJSON(data []byte) string {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(out)
}

// PrettyJSON dışarıdan da kullanılabilsin (TUI için)
func PrettyJSON(data []byte) string {
	return prettyJSON(data)
}
