package cmd

import (
	"fmt"
	"musescoreDownloader/cmd/dl"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "musescoreDownloader",
	Short: "Download musescore sheets as pdf",
	Long: `Ever wanted to see your musescore sheet music offline?

	Musescore allows you to see pdfs online as much as you want for free
	But downloading isn't allowed
	So here you go!`,
	// Run: func (cmd *cobra.Command, args []string){
	// 	cmd.Help()
	// },
}

func Execute(){
	if err := rootCmd.Execute();err!=nil {
		fmt.Println("cmd execution failed")
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(dl.DlCmd)
}
