package fcs

import (
	"fmt"
	"log"
	"os"

	"github.com/immunoconductor/cyto/fcs"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	InputFile  string
	OutputFile string
	Names      bool
	ShortNames bool
	Transform  bool
)

func init() {
	FcsCmd.Flags().StringVarP(&InputFile, "input", "i", "", "input file to read from")
	FcsCmd.Flags().StringVarP(&OutputFile, "output", "o", "", "output file to write to")
	FcsCmd.Flags().BoolVarP(&ShortNames, "shortnames", "s", false, "whether the output file should contain names or friendly names (shortnames) as fields, to be used with input and output flags")
	FcsCmd.Flags().BoolVarP(&Transform, "transform", "t", false, "whether to apply asinh transformation to the data (cofactor of 5 is used)")
	FcsCmd.Flags().BoolVarP(&Names, "names", "n", false, "names")
}

var FcsCmd = &cobra.Command{
	Use:   "fcs",
	Short: "Convert fcs files to csv",
	Long:  `Convert fcs files to csv (as defined by the flow cytometry data file standard)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fcs is a Cyto library for reading Flow Cytometry Standard (FCS) data")
		bar := progressbar.Default(int64(1))

		inputFile, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatal(err)
		}
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}
		shortnames, err := cmd.Flags().GetBool("shortnames")
		if err != nil {
			log.Fatal(err)
		}
		transform, err := cmd.Flags().GetBool("transform")
		if err != nil {
			log.Fatal(err)
		}
		names, err := cmd.Flags().GetBool("names")
		if err != nil {
			log.Fatal(err)
		}

		fcsData, err := fcs.Read(inputFile, transform)
		if err != nil {
			log.Fatal(err)
		}

		if names {
			fmt.Println(fcsData.Names())
		}
		if shortnames {
			fmt.Println(fcsData.ShortNames())
		}

		if outputFile != "" {
			if shortnames {
				fcsData.ToShortNameCSV(outputFile)
			} else {
				fcsData.ToCSV(outputFile)
			}

			fileInfo, err := os.Stat(outputFile)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Output location: %s (%v)\n", outputFile, humanReadableSize(fileInfo.Size()))
		}

		_ = bar.Add(1)
	},
}

func humanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	return fmt.Sprintf("%.2f %s", float64(size)/float64(div), units[exp])
}
