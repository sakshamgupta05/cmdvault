package cmd

import (
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/sakshamgupta05/cmdvault/internal/ui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [term]",
	Short: "Search for commands",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTerm := ""
		if len(args) > 0 {
			searchTerm = args[0]
		}

		collection := collectionFlag
		if collection == "" {
			collection = config.GetDefaultCollection()
		}

		ui.InteractiveSearch(searchTerm)
	},
}

func init() {
	searchCmd.Flags().StringVarP(&collectionFlag, "collection", "c", "", "Collection to search in")
}
