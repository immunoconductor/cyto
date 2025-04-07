package cmd

import (
	"github.com/immunoconductor/cyto/cmd/cyto"
	"github.com/immunoconductor/cyto/cmd/fcs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cyto",
	Short: "cyto",
	Long:  `Cyto is a CLI library for the analysis of cytometry data.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(fcs.FcsCmd)
	rootCmd.AddCommand(cyto.VersionCmd)
}
