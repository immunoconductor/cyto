package fcs

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	FcsCmd.AddCommand(VersionCmd)
}

var FcsCmd = &cobra.Command{
	Use:   "fcs",
	Short: "Manipulate fcs files",
	Long:  `Manipulate fcs files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fcs is a Cyto library for the analysis of Flow Cytometry Standard (FCS) data")
	},
}
