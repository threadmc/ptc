package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <patch>",
	Short: "Show the diff between the current file and the patch content",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		patchFile := args[0]

		f, err := os.Open(patchFile)
		if err != nil {
			fmt.Println("Error reading patch file:", err)
			os.Exit(1)
		}
		defer f.Close()

		var targetFile string
		var inContent bool
		var patchContent []byte

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "TARGET") {
				parts := strings.SplitN(line, "=", 2)
				targetFile = strings.TrimSpace(parts[1])
			}
			if strings.HasPrefix(line, "--- PATCH CONTENT ---") {
				inContent = true
				continue
			}
			if inContent {
				patchContent = append(patchContent, []byte(line+"\n")...)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading patch file:", err)
			os.Exit(1)
		}

		if targetFile == "" {
			fmt.Println("Malformed patch file: missing TARGET")
			os.Exit(1)
		}

		currentContent, err := os.ReadFile(targetFile)
		if err != nil {
			fmt.Println("Error reading target file:", err)
			os.Exit(1)
		}

		// Simple line-by-line diff
		currentLines := bytes.Split(currentContent, []byte("\n"))
		patchLines := bytes.Split(patchContent, []byte("\n"))

		max := len(currentLines)
		if len(patchLines) > max {
			max = len(patchLines)
		}
		for i := 0; i < max; i++ {
			var cur, pat []byte
			if i < len(currentLines) {
				cur = currentLines[i]
			}
			if i < len(patchLines) {
				pat = patchLines[i]
			}
			if !bytes.Equal(cur, pat) {
				fmt.Printf("-%s\n+%s\n", cur, pat)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
