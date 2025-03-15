package fcs

import (
	"fmt"

	"github.com/immunoconductor/cyto/fcs"
	"github.com/spf13/cobra"
)

var InputFile string
var OutputFile string
var ShortNames bool

func init() {
	FcsCmd.Flags().StringVarP(&InputFile, "input", "i", "", "input file to read from")
	FcsCmd.Flags().StringVarP(&OutputFile, "output", "o", "", "output file to write to")
	FcsCmd.MarkFlagsRequiredTogether("input", "output")
	FcsCmd.Flags().BoolVarP(&ShortNames, "shortnames", "s", false, "shortnames output")
	FcsCmd.AddCommand(VersionCmd)
}

var FcsCmd = &cobra.Command{
	Use:   "fcs",
	Short: "Read fcs files to csv",
	Long:  `Read fcs files to csv (as defined by the flow cytometry data file standard)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fcs is a Cyto library for reading Flow Cytometry Standard (FCS) data")

		inputFile, err := cmd.Flags().GetString("input")
		if err != nil {
			return
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return
		}
		shortnames, err := cmd.Flags().GetBool("shortnames")
		if err != nil {
			return
		}
		fcs, err := fcs.NewFCS(inputFile)
		if err != nil {
			return
		}
		if shortnames {
			fcs.ToShortNameCSV(outputFile)
		} else {
			fcs.ToCSV(outputFile)
		}
	},
}
