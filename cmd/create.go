package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"ptc/helpers"

	"github.com/spf13/cobra"
)

var description string

var createCmd = &cobra.Command{
	Use:   "create <file> <patch>",
	Short: "Create a patch from a modified file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		patch := args[1]
		origBackup := file + ".bak"

		if _, err := os.Stat(origBackup); os.IsNotExist(err) {
			fmt.Println("No backup found. Creating backup automatically:", origBackup)
			err := helpers.CopyFile(file, origBackup)
			if err != nil {
				fmt.Println("Error creating backup:", err)
				os.Exit(1)
			}
		}

		origContent, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		origHash := helpers.SHA256Sum(origContent)

		patchAbs, err := filepath.Abs(patch)
		if err != nil {
			fmt.Println("Error resolving patch file path:", err)
			os.Exit(1)
		}
		patchDir := filepath.Dir(patchAbs)
		fileAbs, err := filepath.Abs(file)
		if err != nil {
			fmt.Println("Error resolving target file path:", err)
			os.Exit(1)
		}
		relTarget, err := filepath.Rel(patchDir, fileAbs)
		if err != nil {
			relTarget = file
		}
		relTarget = filepath.ToSlash(relTarget)

		patchFile, err := os.Create(patch)
		if err != nil {
			fmt.Println("Error creating patch file:", err)
			os.Exit(1)
		}
		defer patchFile.Close()

		// Write patch header with relative target path
		patchFile.WriteString(fmt.Sprintf("# PTC PATCH v1\nTARGET = %s\nHASH = %s\nDESCRIPTION = %s\n\n", relTarget, origHash, description))
		patchFile.WriteString("--- PATCH CONTENT ---\n")
		patchFile.Write(origContent)
	},
}

func init() {
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Description for the patch")
	rootCmd.AddCommand(createCmd)
}
