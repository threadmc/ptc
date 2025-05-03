package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ptc",
	Short:   "A lightweight patch manager for submodules",
	Long:    "PTC creates and applies patches inside submodules without Git commits.",
	Version: "1.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
