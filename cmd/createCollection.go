package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sakshamgupta05/cmdvault/internal/config"
	"github.com/spf13/cobra"
)

var createCollectionCmd = &cobra.Command{
	Use:   "create-collection [name]",
	Short: "Create a new collection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := config.AddCollection(name); err != nil {
			fmt.Printf("Error creating collection: %v\n", err)
			return
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Collection \"%s\" created successfully!\n", green("✓"), name)
	},
}
