package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"ptc/helpers"
)

var applyCmd = &cobra.Command{
	Use:   "apply <patch>",
	Short: "Apply a patch to the target file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		patchFile := args[0]

		f, err := os.Open(patchFile)
		if err != nil {
			fmt.Println("Error reading patch file:", err)
			os.Exit(1)
		}
		defer f.Close()

		var targetFile, expectedHash string
		var inContent bool
		var patchContent []byte

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "TARGET") {
				parts := strings.SplitN(line, "=", 2)
				targetFile = strings.TrimSpace(parts[1])
			}
			if strings.HasPrefix(line, "HASH") {
				parts := strings.SplitN(line, "=", 2)
				expectedHash = strings.TrimSpace(parts[1])
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

		if targetFile == "" || expectedHash == "" {
			fmt.Println("Malformed patch file: missing TARGET or HASH")
			os.Exit(1)
		}

		currentContent, err := os.ReadFile(targetFile)
		if err != nil {
			fmt.Println("Error reading target file:", err)
			os.Exit(1)
		}
		currentHash := helpers.SHA256Sum(currentContent)

		if currentHash != expectedHash {
			fmt.Println("ERROR: File hash mismatch! Patch cannot be safely applied.")
			fmt.Println("Expected:", expectedHash)
			fmt.Println("Found   :", currentHash)
			os.Exit(1)
		}

		// Remove trailing newline if present (optional)
		if len(patchContent) > 0 && patchContent[len(patchContent)-1] == '\n' {
			patchContent = patchContent[:len(patchContent)-1]
		}

		err = os.WriteFile(targetFile, patchContent, 0644)
		if err != nil {
			fmt.Println("Error applying patch:", err)
			os.Exit(1)
		}

		fmt.Println("Patch successfully applied.")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
