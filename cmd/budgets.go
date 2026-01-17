package cmd

import (
	"fmt"

	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var budgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "Manage budgets",
	Long:  `List and view budget information.`,
}

var budgetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budgets",
	Long:  `Returns a list of all budgets associated with your YNAB account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgets, err := apiClient.GetBudgets()
		if err != nil {
			return fmt.Errorf("failed to get budgets: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(budgets)
	},
}

var budgetsGetCmd = &cobra.Command{
	Use:   "get [budget-id]",
	Short: "Get budget details",
	Long: `Returns details for a specific budget.

If no budget ID is provided, uses the default budget from config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var id string
		if len(args) > 0 {
			id = args[0]
		} else {
			var err error
			id, err = getBudgetID()
			if err != nil {
				return err
			}
		}

		budget, err := apiClient.GetBudget(id)
		if err != nil {
			return fmt.Errorf("failed to get budget: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(budget)
	},
}

var budgetsSettingsCmd = &cobra.Command{
	Use:   "settings [budget-id]",
	Short: "Get budget settings",
	Long: `Returns settings for a specific budget.

If no budget ID is provided, uses the default budget from config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var id string
		if len(args) > 0 {
			id = args[0]
		} else {
			var err error
			id, err = getBudgetID()
			if err != nil {
				return err
			}
		}

		settings, err := apiClient.GetBudgetSettings(id)
		if err != nil {
			return fmt.Errorf("failed to get budget settings: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(settings)
	},
}

func init() {
	rootCmd.AddCommand(budgetsCmd)
	budgetsCmd.AddCommand(budgetsListCmd)
	budgetsCmd.AddCommand(budgetsGetCmd)
	budgetsCmd.AddCommand(budgetsSettingsCmd)
}
