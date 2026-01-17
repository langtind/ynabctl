package cmd

import (
	"fmt"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage accounts",
	Long:  `List, view, and create accounts within a budget.`,
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  `Returns a list of all accounts for the specified budget.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := getBudgetID()
		if err != nil {
			return err
		}

		accounts, err := apiClient.GetAccounts(id)
		if err != nil {
			return fmt.Errorf("failed to get accounts: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(accounts)
	},
}

var accountsGetCmd = &cobra.Command{
	Use:   "get <account-id>",
	Short: "Get account details",
	Long:  `Returns details for a specific account.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		account, err := apiClient.GetAccount(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(account)
	},
}

var (
	accountName    string
	accountType    string
	accountBalance float64
)

var accountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	Long: `Create a new account in the budget.

Account types: checking, savings, cash, creditCard, lineOfCredit,
otherAsset, otherLiability, mortgage, autoLoan, studentLoan,
personalLoan, medicalDebt, otherDebt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		if accountName == "" {
			return fmt.Errorf("account name is required (--name)")
		}
		if accountType == "" {
			return fmt.Errorf("account type is required (--type)")
		}

		// Convert balance to milliunits
		balance := client.AmountToMilliunits(accountBalance)

		account, err := apiClient.CreateAccount(budgetID, accountName, accountType, balance)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(account)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsGetCmd)
	accountsCmd.AddCommand(accountsCreateCmd)

	accountsCreateCmd.Flags().StringVar(&accountName, "name", "", "Account name (required)")
	accountsCreateCmd.Flags().StringVar(&accountType, "type", "", "Account type (required)")
	accountsCreateCmd.Flags().Float64Var(&accountBalance, "balance", 0, "Starting balance")
}
