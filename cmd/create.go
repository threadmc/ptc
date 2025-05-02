package cmd

import (
	"fmt"
	"os"
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

		patchFile, err := os.Create(patch)
		if err != nil {
			fmt.Println("Error creating patch file:", err)
			os.Exit(1)
		}
		defer patchFile.Close()

		patchFile.WriteString(fmt.Sprintf("# PTC PATCH v1\nTARGET = %s\nHASH = %s\n\n", file, origHash))

		patchFile.WriteString("--- ORIGINAL\n")
		patchFile.Write(origContent)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
