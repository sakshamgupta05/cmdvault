package cmd

import (
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/sakshamgupta05/cmdvault/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commands",
	Run: func(cmd *cobra.Command, args []string) {
		collection := collectionFlag
		if collection == "" {
			collection = config.GetDefaultCollection()
		}

		ui.ListCommands(collection)
	},
}

func init() {
	listCmd.Flags().StringVarP(&collectionFlag, "collection", "c", "", "Collection to list")
}
