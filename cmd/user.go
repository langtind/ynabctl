package cmd

import (
	"fmt"

	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get authenticated user information",
	Long:  `Returns information about the authenticated user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := apiClient.GetUser()
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(user)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
