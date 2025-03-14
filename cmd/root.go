package cmd

import (
	"github.com/immunoconductor/cyto/cmd/fcs"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	// cfgFile     string
	// userLicense string

	rootCmd = &cobra.Command{
		Use:   "cyto",
		Short: "cyto",
		Long:  `Cyto is a CLI library for the analysis of CyTOF data.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(fcs.FcsCmd)
}
