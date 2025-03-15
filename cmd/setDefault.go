package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/spf13/cobra"
)

var setDefaultCmd = &cobra.Command{
	Use:   "set-default [collection]",
	Short: "Set the default collection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collection := args[0]

		// Check if collection exists
		cfg := config.GetConfig()
		exists := false
		for _, c := range cfg.Collections {
			if c == collection {
				exists = true
				break
			}
		}

		if !exists {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("%s Collection \"%s\" does not exist. Create it first with create-collection.\n", yellow("!"), collection)
			return
		}

		if err := config.SetDefaultCollection(collection); err != nil {
			fmt.Printf("Error setting default collection: %v\n", err)
			return
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Default collection set to \"%s\".\n", green("✓"), collection)
	},
}
