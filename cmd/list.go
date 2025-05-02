package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

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

		type patchInfo struct {
			File        string
			Target      string
			Description string
		}
		var patches []patchInfo

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".ptc" {
				patchPath := filepath.Join(dir, file.Name())
				f, err := os.Open(patchPath)
				if err != nil {
					continue
				}
				defer f.Close()

				var target, desc string
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "TARGET") {
						parts := strings.SplitN(line, "=", 2)
						target = strings.TrimSpace(parts[1])
					}
					if strings.HasPrefix(line, "DESCRIPTION") {
						parts := strings.SplitN(line, "=", 2)
						desc = strings.TrimSpace(parts[1])
					}
					if target != "" && desc != "" {
						break
					}
				}
				patches = append(patches, patchInfo{
					File:        file.Name(),
					Target:      target,
					Description: desc,
				})
			}
		}

		if len(patches) == 0 {
			fmt.Println("No .ptc patch files found in", dir)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PATCH FILE\tTARGET FILE\tDESCRIPTION")
		for _, p := range patches {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.File, p.Target, p.Description)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
