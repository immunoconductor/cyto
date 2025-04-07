package cyto

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the cyto library",
	Long:  `Print the version number of the cyto library`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cyto version cyto1.0.0")
	},
}
