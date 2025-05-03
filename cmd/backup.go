package cmd

import (
	"fmt"
	"os"
	"ptc/helpers"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup <file>",
	Short: "Create a backup of a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		backup := file + ".bak"
		err := helpers.CopyFile(file, backup)
		if err != nil {
			fmt.Println("Error creating backup:", err)
			os.Exit(1)
		}
		fmt.Println("Backup created:", backup)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
