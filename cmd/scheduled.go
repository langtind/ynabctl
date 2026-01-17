package cmd

import (
	"fmt"
	"time"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/output"
	"github.com/spf13/cobra"
)

var scheduledCmd = &cobra.Command{
	Use:     "scheduled",
	Aliases: []string{"scheduled-transactions"},
	Short:   "Manage scheduled transactions",
	Long:    `List, view, create, update, and delete scheduled transactions.`,
}

var scheduledListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scheduled transactions",
	Long:  `Returns a list of all scheduled transactions for the budget.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		transactions, err := apiClient.GetScheduledTransactions(budgetID)
		if err != nil {
			return fmt.Errorf("failed to get scheduled transactions: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transactions)
	},
}

var scheduledGetCmd = &cobra.Command{
	Use:   "get <scheduled-transaction-id>",
	Short: "Get scheduled transaction details",
	Long:  `Returns details for a specific scheduled transaction.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		transaction, err := apiClient.GetScheduledTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get scheduled transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var (
	schedAccountID  string
	schedDate       string
	schedFrequency  string
	schedAmount     float64
	schedPayeeID    string
	schedPayeeName  string
	schedCategoryID string
	schedMemo       string
	schedFlagColor  string
)

var scheduledCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a scheduled transaction",
	Long: `Create a new scheduled transaction.

Required flags:
  --account: Account ID
  --date: First occurrence date (YYYY-MM-DD)
  --frequency: Recurrence frequency
  --amount: Transaction amount

Frequency options:
  never, daily, weekly, everyOtherWeek, twiceAMonth,
  every4Weeks, monthly, everyOtherMonth, every3Months,
  every4Months, twiceAYear, yearly, everyOtherYear`,
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		if schedAccountID == "" {
			return fmt.Errorf("account ID is required (--account)")
		}
		if schedFrequency == "" {
			return fmt.Errorf("frequency is required (--frequency)")
		}

		date := schedDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		st := client.SaveScheduledTransaction{
			AccountID:  schedAccountID,
			Date:       date,
			Frequency:  schedFrequency,
			Amount:     client.AmountToMilliunits(schedAmount),
			PayeeID:    schedPayeeID,
			PayeeName:  schedPayeeName,
			CategoryID: schedCategoryID,
			Memo:       schedMemo,
			FlagColor:  schedFlagColor,
		}

		transaction, err := apiClient.CreateScheduledTransaction(budgetID, st)
		if err != nil {
			return fmt.Errorf("failed to create scheduled transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var scheduledUpdateCmd = &cobra.Command{
	Use:   "update <scheduled-transaction-id>",
	Short: "Update a scheduled transaction",
	Long:  `Update an existing scheduled transaction.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		// Get existing scheduled transaction
		existing, err := apiClient.GetScheduledTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get existing scheduled transaction: %w", err)
		}

		st := client.SaveScheduledTransaction{
			AccountID:  existing.AccountID,
			Date:       existing.DateFirst,
			Frequency:  existing.Frequency,
			Amount:     existing.Amount,
			PayeeID:    existing.PayeeID,
			CategoryID: existing.CategoryID,
			Memo:       existing.Memo,
			FlagColor:  existing.FlagColor,
		}

		if cmd.Flags().Changed("account") {
			st.AccountID = schedAccountID
		}
		if cmd.Flags().Changed("date") {
			st.Date = schedDate
		}
		if cmd.Flags().Changed("frequency") {
			st.Frequency = schedFrequency
		}
		if cmd.Flags().Changed("amount") {
			st.Amount = client.AmountToMilliunits(schedAmount)
		}
		if cmd.Flags().Changed("payee-id") {
			st.PayeeID = schedPayeeID
		}
		if cmd.Flags().Changed("payee-name") {
			st.PayeeName = schedPayeeName
		}
		if cmd.Flags().Changed("category") {
			st.CategoryID = schedCategoryID
		}
		if cmd.Flags().Changed("memo") {
			st.Memo = schedMemo
		}
		if cmd.Flags().Changed("flag") {
			st.FlagColor = schedFlagColor
		}

		transaction, err := apiClient.UpdateScheduledTransaction(budgetID, args[0], st)
		if err != nil {
			return fmt.Errorf("failed to update scheduled transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

var scheduledDeleteCmd = &cobra.Command{
	Use:   "delete <scheduled-transaction-id>",
	Short: "Delete a scheduled transaction",
	Long:  `Delete a scheduled transaction.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		budgetID, err := getBudgetID()
		if err != nil {
			return err
		}

		transaction, err := apiClient.DeleteScheduledTransaction(budgetID, args[0])
		if err != nil {
			return fmt.Errorf("failed to delete scheduled transaction: %w", err)
		}

		formatter := output.New(getOutputFormat())
		return formatter.Print(transaction)
	},
}

func init() {
	rootCmd.AddCommand(scheduledCmd)
	scheduledCmd.AddCommand(scheduledListCmd)
	scheduledCmd.AddCommand(scheduledGetCmd)
	scheduledCmd.AddCommand(scheduledCreateCmd)
	scheduledCmd.AddCommand(scheduledUpdateCmd)
	scheduledCmd.AddCommand(scheduledDeleteCmd)

	// Create flags
	scheduledCreateCmd.Flags().StringVar(&schedAccountID, "account", "", "Account ID (required)")
	scheduledCreateCmd.Flags().StringVar(&schedDate, "date", "", "First occurrence date (YYYY-MM-DD)")
	scheduledCreateCmd.Flags().StringVar(&schedFrequency, "frequency", "", "Recurrence frequency (required)")
	scheduledCreateCmd.Flags().Float64Var(&schedAmount, "amount", 0, "Amount")
	scheduledCreateCmd.Flags().StringVar(&schedPayeeID, "payee-id", "", "Payee ID")
	scheduledCreateCmd.Flags().StringVar(&schedPayeeName, "payee-name", "", "Payee name")
	scheduledCreateCmd.Flags().StringVar(&schedCategoryID, "category", "", "Category ID")
	scheduledCreateCmd.Flags().StringVar(&schedMemo, "memo", "", "Memo")
	scheduledCreateCmd.Flags().StringVar(&schedFlagColor, "flag", "", "Flag color")

	// Update flags
	scheduledUpdateCmd.Flags().StringVar(&schedAccountID, "account", "", "Account ID")
	scheduledUpdateCmd.Flags().StringVar(&schedDate, "date", "", "Date (YYYY-MM-DD)")
	scheduledUpdateCmd.Flags().StringVar(&schedFrequency, "frequency", "", "Recurrence frequency")
	scheduledUpdateCmd.Flags().Float64Var(&schedAmount, "amount", 0, "Amount")
	scheduledUpdateCmd.Flags().StringVar(&schedPayeeID, "payee-id", "", "Payee ID")
	scheduledUpdateCmd.Flags().StringVar(&schedPayeeName, "payee-name", "", "Payee name")
	scheduledUpdateCmd.Flags().StringVar(&schedCategoryID, "category", "", "Category ID")
	scheduledUpdateCmd.Flags().StringVar(&schedMemo, "memo", "", "Memo")
	scheduledUpdateCmd.Flags().StringVar(&schedFlagColor, "flag", "", "Flag color")
}
