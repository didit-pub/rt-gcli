package commands

import (
	"fmt"

	"github.com/didit-pub/rt-gcli/pkg/updater"
	"github.com/didit-pub/rt-gcli/pkg/version"
	"github.com/spf13/cobra"
)

func newUpdateMeCmd() *cobra.Command {
	var checkOnly bool

	cmd := &cobra.Command{
		Use:   "updateme",
		Short: "Check and apply updates",
		Run: func(cmd *cobra.Command, args []string) {
			checkOnly, _ := cmd.Flags().GetBool("check")

			fmt.Printf("Current version: %s\n", version.GetVersion())
			fmt.Println("Checking for updates...")

			release, hasUpdate, err := updater.CheckForUpdates()
			if err != nil {
				fmt.Printf("Error checking for updates: %v\n", err)
				return
			}

			if !hasUpdate {
				fmt.Println("You already have the latest version.")
				return
			}

			fmt.Printf("New version available: %s\n", release.TagName)

			if checkOnly {
				return
			}

			fmt.Println("Starting update...")
			if err := updater.DoSelfUpdate(release); err != nil {
				fmt.Printf("Error during update: %v\n", err)
				return
			}

			fmt.Println("Update completed successfully!")
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check", false, "Check for updates without applying them")

	return cmd
}
