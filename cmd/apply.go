package cmd

import (
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

		content, err := os.ReadFile(patchFile)
		if err != nil {
			fmt.Println("Error reading patch file:", err)
			os.Exit(1)
		}

		var targetFile, expectedHash string
		var patchContent []byte
		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			if strings.HasPrefix(line, "TARGET") {
				parts := strings.SplitN(line, "=", 2)
				targetFile = strings.TrimSpace(parts[1])
			}
			if strings.HasPrefix(line, "HASH") {
				parts := strings.SplitN(line, "=", 2)
				expectedHash = strings.TrimSpace(parts[1])
			}
			if strings.HasPrefix(line, "---") {
				patchContent = append(patchContent, []byte(strings.Join(lines[1:], "\n"))...)
			}
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
