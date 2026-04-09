package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mustafacavusoglu/hill/internal/httpclient"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "HTTP GET isteği gönder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		headers, _ := cmd.Flags().GetStringArray("header")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		req := httpclient.Request{
			Method:  "GET",
			URL:     args[0],
			Headers: parseHeaders(headers),
			Timeout: timeout,
		}

		resp, err := httpclient.Execute(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
			os.Exit(1)
		}
		httpclient.PrintResponse(resp)
		return nil
	},
}

var postCmd = &cobra.Command{
	Use:   "post <url>",
	Short: "HTTP POST isteği gönder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		headers, _ := cmd.Flags().GetStringArray("header")
		data, _ := cmd.Flags().GetString("data")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		parsedHeaders := parseHeaders(headers)
		if _, ok := parsedHeaders["Content-Type"]; !ok && data != "" {
			parsedHeaders["Content-Type"] = "application/json"
		}

		req := httpclient.Request{
			Method:  "POST",
			URL:     args[0],
			Headers: parsedHeaders,
			Body:    []byte(data),
			Timeout: timeout,
		}

		resp, err := httpclient.Execute(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
			os.Exit(1)
		}
		httpclient.PrintResponse(resp)
		return nil
	},
}

var putCmd = &cobra.Command{
	Use:   "put <url>",
	Short: "HTTP PUT isteği gönder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		headers, _ := cmd.Flags().GetStringArray("header")
		data, _ := cmd.Flags().GetString("data")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		parsedHeaders := parseHeaders(headers)
		if _, ok := parsedHeaders["Content-Type"]; !ok && data != "" {
			parsedHeaders["Content-Type"] = "application/json"
		}

		req := httpclient.Request{
			Method:  "PUT",
			URL:     args[0],
			Headers: parsedHeaders,
			Body:    []byte(data),
			Timeout: timeout,
		}

		resp, err := httpclient.Execute(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
			os.Exit(1)
		}
		httpclient.PrintResponse(resp)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <url>",
	Short: "HTTP DELETE isteği gönder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		headers, _ := cmd.Flags().GetStringArray("header")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		req := httpclient.Request{
			Method:  "DELETE",
			URL:     args[0],
			Headers: parseHeaders(headers),
			Timeout: timeout,
		}

		resp, err := httpclient.Execute(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
			os.Exit(1)
		}
		httpclient.PrintResponse(resp)
		return nil
	},
}

// parseHeaders "Key: Value" formatındaki string dizisini map'e dönüştürür
func parseHeaders(headers []string) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func init() {
	// GET
	getCmd.Flags().StringArrayP("header", "H", nil, "Header (örn: -H 'Authorization: Bearer token')")
	getCmd.Flags().DurationP("timeout", "t", 30000000000, "Timeout (örn: 10s, 1m)")
	rootCmd.AddCommand(getCmd)

	// POST
	postCmd.Flags().StringArrayP("header", "H", nil, "Header")
	postCmd.Flags().StringP("data", "d", "", "Request body (JSON)")
	postCmd.Flags().DurationP("timeout", "t", 30000000000, "Timeout")
	rootCmd.AddCommand(postCmd)

	// PUT
	putCmd.Flags().StringArrayP("header", "H", nil, "Header")
	putCmd.Flags().StringP("data", "d", "", "Request body (JSON)")
	putCmd.Flags().DurationP("timeout", "t", 30000000000, "Timeout")
	rootCmd.AddCommand(putCmd)

	// DELETE
	deleteCmd.Flags().StringArrayP("header", "H", nil, "Header")
	deleteCmd.Flags().DurationP("timeout", "t", 30000000000, "Timeout")
	rootCmd.AddCommand(deleteCmd)
}
