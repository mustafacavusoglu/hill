package cmd

import (
	"fmt"
	"os"

	"github.com/mustafacavusoglu/hill/internal/checker"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <ip-veya-host>",
	Short: "Ağ bağlantısı ve IP kontrolü",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")

		result, err := checker.Run(args[0], port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Check hatası: %v\n", err)
			os.Exit(1)
		}
		checker.PrintResult(result)
		return nil
	},
}

func init() {
	checkCmd.Flags().IntP("port", "p", 0, "Özel port kontrolü (varsayılan: 80 ve 443)")
	rootCmd.AddCommand(checkCmd)
}
