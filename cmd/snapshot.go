package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/langtind/ynabctl/internal/client"
	"github.com/langtind/ynabctl/internal/period"
	"github.com/spf13/cobra"
)

var (
	snapshotPeriod   string
	snapshotSpecific string
	snapshotOut      string
)

type snapshot struct {
	Period       period.Range                  `json:"period"`
	FetchedAt    string                        `json:"fetched_at"`
	BudgetID     string                        `json:"budget_id"`
	Accounts     []client.Account              `json:"accounts"`
	Categories   []client.CategoryGroup        `json:"categories"`
	Payees       []client.Payee                `json:"payees"`
	Months       []client.Month                `json:"months"`
	Transactions []client.Transaction          `json:"transactions"`
	Scheduled    []client.ScheduledTransaction `json:"scheduled_transactions"`
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Fetch a combined data snapshot for a period",
	Long: `Fetch accounts, categories, payees, months, transactions, and
scheduled transactions for a given period and emit a single JSON document.

Transactions are filtered by since_date = period start. Writes to stdout
unless --out is given.`,
	Example: `  ynabctl snapshot --period month
  ynabctl snapshot --period quarter --specific 2026-Q1
  ynabctl snapshot --period week --out data/raw/current-week.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bID, err := getBudgetID()
		if err != nil {
			return err
		}
		p, err := period.Compute(snapshotPeriod, snapshotSpecific)
		if err != nil {
			return err
		}

		accounts, err := apiClient.GetAccounts(bID)
		if err != nil {
			return fmt.Errorf("accounts: %w", err)
		}
		cats, err := apiClient.GetCategories(bID)
		if err != nil {
			return fmt.Errorf("categories: %w", err)
		}
		payees, err := apiClient.GetPayees(bID)
		if err != nil {
			return fmt.Errorf("payees: %w", err)
		}
		months, err := apiClient.GetMonths(bID)
		if err != nil {
			return fmt.Errorf("months: %w", err)
		}
		txns, err := apiClient.GetTransactions(bID, &client.TransactionFilter{SinceDate: p.StartDate})
		if err != nil {
			return fmt.Errorf("transactions: %w", err)
		}
		sched, err := apiClient.GetScheduledTransactions(bID)
		if err != nil {
			return fmt.Errorf("scheduled: %w", err)
		}

		snap := snapshot{
			Period:       p,
			FetchedAt:    time.Now().UTC().Format(time.RFC3339),
			BudgetID:     bID,
			Accounts:     accounts,
			Categories:   cats,
			Payees:       payees,
			Months:       months,
			Transactions: txns,
			Scheduled:    sched,
		}

		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return err
		}

		if snapshotOut != "" {
			if err := os.WriteFile(snapshotOut, data, 0o644); err != nil {
				return fmt.Errorf("write %s: %w", snapshotOut, err)
			}
			fmt.Fprintf(os.Stderr, "snapshot saved: %s (%s → %s)\n", snapshotOut, p.StartDate, p.EndDate)
			return nil
		}
		_, err = os.Stdout.Write(data)
		if err == nil {
			_, _ = os.Stdout.WriteString("\n")
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().StringVar(&snapshotPeriod, "period", "", "Period kind: week|month|quarter|year (required)")
	snapshotCmd.Flags().StringVar(&snapshotSpecific, "specific", "", "Specific period (e.g. 2026-03, 2026-W15, 2026-Q1, 2026)")
	snapshotCmd.Flags().StringVar(&snapshotOut, "out", "", "Write JSON to file instead of stdout")
	_ = snapshotCmd.MarkFlagRequired("period")
}
