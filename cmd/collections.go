package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/spf13/cobra"
)

var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "List all collections",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("%s Available collections:\n\n", bold("•"))
		for _, collection := range cfg.Collections {
			if collection == cfg.DefaultCollection {
				fmt.Printf("  %s (default)\n", collection)
			} else {
				fmt.Printf("  %s\n", collection)
			}
		}
	},
}
