package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <directory>",
	Short: "List all patch files in a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			os.Exit(1)
		}

		count := 0
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".ptc" {
				fmt.Println(file.Name())
				count++
			}
		}

		if count == 0 {
			fmt.Println("No .ptc patch files found in", dir)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
