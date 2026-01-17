package cmd

import (
	"fmt"

	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var payeesCmd = &cobra.Command{
	Use:   "payees",
	Short: "Manage payees",
	Long:  `List, view, and update payees.`,
}

var payeesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all payees",
	Long:  `Returns a list of all payees for the budget.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		payees, err := apiClient.GetPayees(budgetID)
		if err != nil {
			return fmt.Errorf("failed to get payees: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(payees)
	},
}

var payeesGetCmd = &cobra.Command{
	Use:   "get <payee-id>",
	Short: "Get payee details",
	Long:  `Returns details for a specific payee.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		payee, err := apiClient.GetPayee(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get payee: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(payee)
	},
}

var payeeNewName string

var payeesUpdateCmd = &cobra.Command{
	Use:   "update <payee-id>",
	Short: "Update a payee",
	Long:  `Update a payee's name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		if payeeNewName == "" {
			return fmt.Errorf("new name is required (--name)")
		}

		payee, err := apiClient.UpdatePayee(budgetID, args[0], payeeNewName)
		if err != nil {
			return fmt.Errorf("failed to update payee: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(payee)
	},
}

func init() {
	rootCmd.AddCommand(payeesCmd)
	payeesCmd.AddCommand(payeesListCmd)
	payeesCmd.AddCommand(payeesGetCmd)
	payeesCmd.AddCommand(payeesUpdateCmd)

	payeesUpdateCmd.Flags().StringVar(&payeeNewName, "name", "", "New payee name (required)")
}
