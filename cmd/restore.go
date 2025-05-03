package cmd

import (
	"fmt"
	"os"
	"ptc/helpers"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore <file>",
	Short: "Restore a file from its backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		backup := file + ".bak"
		if _, err := os.Stat(backup); os.IsNotExist(err) {
			fmt.Println("No backup found for", file)
			os.Exit(1)
		}
		err := helpers.CopyFile(backup, file)
		if err != nil {
			fmt.Println("Error restoring file:", err)
			os.Exit(1)
		}
		fmt.Println("File restored from backup:", backup)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
