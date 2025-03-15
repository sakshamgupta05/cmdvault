package cmd

import (
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/sakshamgupta05/cmdvault/internal/ui"
	"github.com/spf13/cobra"
)

var collectionFlag string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new command",
	Run: func(cmd *cobra.Command, args []string) {
		collection := collectionFlag
		if collection == "" {
			collection = config.GetDefaultCollection()
		}
		ui.AddCommandPrompt(collection)
	},
}

func init() {
	addCmd.Flags().StringVarP(&collectionFlag, "collection", "c", "", "Collection to add to")
}
