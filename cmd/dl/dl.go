package dl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DlCmd = &cobra.Command{
	Use:   "dl",
	Short: "dl is a palette that contains download functionalities",
	Long:  "Downloads a pdf from a url",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		remList := args
		rem := make(chan string, 10)

		for _, url := range args {
			go DownloadFromUrl(url, &rem)
		}

		fmt.Println("\nRemaining: ", len(remList))

		for i := 0; i < len(args); i++ {
			completedUrl := <-rem
			for i, v := range remList {
				if v == completedUrl {
					remList[i] = remList[len(remList)-1]
					remList = remList[:len(remList)-1]
				}
			}
			fmt.Println("\nRemaining: ", len(remList), "\n", remList)
		}
	},
}

func init() {
	// DlCmd.Flags().String
}
