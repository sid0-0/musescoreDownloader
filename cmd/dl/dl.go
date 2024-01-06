package dl

import (

	"github.com/spf13/cobra"
)

var DlCmd = &cobra.Command{
	Use: "dl",
	Short: "dl is a palette that contains download functionalities",
	Long: "Downloads a pdf from a url",
	Run: func(cmd *cobra.Command, args []string){
		if len( args ) == 0 {
			cmd.Help()
			return
		}
		for _,url := range args {
			DownloadFromUrl(url)
		}
	},
}

func init(){
	// DlCmd.Flags().String
}
