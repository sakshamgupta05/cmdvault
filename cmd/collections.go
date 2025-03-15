package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/store"
	"github.com/spf13/cobra"
)

var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "List all collections",
	Run: func(cmd *cobra.Command, args []string) {
		collections, err := store.ListCollections()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("%s Available collections:\n\n", bold("•"))
		for _, collection := range collections {
			fmt.Printf("  %s\n", collection)
		}
	},
}
