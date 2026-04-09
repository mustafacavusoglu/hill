package cmd

import (
	"os"

	"github.com/mustafacavusoglu/hill/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hill",
	Short: "hill — Terminal HTTP client, load tester ve network checker",
	Long: `hill: Postman + hey karışımı terminal aracı.

  hill                         → TUI modu
  hill get <url>               → HTTP GET
  hill post <url> -d '{...}'   → HTTP POST
  hill benchmark -n 1000 -c 50 → Yük testi
  hill check <ip/host>         → Ağ kontrolü`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
