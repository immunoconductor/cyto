package fcs

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the fcs library",
	Long:  `Print the version number of the fcs library`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fcs: Flow Cytometry Standard (FCS) library v0.1")
	},
}
