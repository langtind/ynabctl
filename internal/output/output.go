package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/langtind/ynabctl/internal/client"
)

// Formatter handles output formatting
type Formatter struct {
	format string
	writer io.Writer
}

// New creates a new output formatter
func New(format string) *Formatter {
	return &Formatter{
		format: format,
		writer: os.Stdout,
	}
}

// Print outputs data in the configured format
func (f *Formatter) Print(data interface{}) error {
	if f.format == "table" {
		return f.printTable(data)
	}
	return f.printJSON(data)
}

// printJSON outputs data as pretty-printed JSON
func (f *Formatter) printJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTable outputs data in tabular format
func (f *Formatter) printTable(data interface{}) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	switch v := data.(type) {
	case *client.User:
		fmt.Fprintln(w, "ID")
		fmt.Fprintf(w, "%s\n", v.ID)

	case []client.Budget:
		fmt.Fprintln(w, "ID\tNAME\tLAST MODIFIED")
		for _, b := range v {
			fmt.Fprintf(w, "%s\t%s\t%s\n", b.ID, b.Name, b.LastModifiedOn)
		}

	case *client.Budget:
		fmt.Fprintln(w, "ID\tNAME\tFIRST MONTH\tLAST MONTH")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", v.ID, v.Name, v.FirstMonth, v.LastMonth)

	case *client.BudgetSettings:
		fmt.Fprintln(w, "SETTING\tVALUE")
		fmt.Fprintf(w, "Date Format\t%s\n", v.DateFormat.Format)
		fmt.Fprintf(w, "Currency\t%s\n", v.CurrencyFormat.ISOCode)
		fmt.Fprintf(w, "Currency Symbol\t%s\n", v.CurrencyFormat.CurrencySymbol)
		fmt.Fprintf(w, "Decimal Digits\t%d\n", v.CurrencyFormat.DecimalDigits)

	case []client.Account:
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tBALANCE\tON BUDGET\tCLOSED")
		for _, a := range v {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.2f\t%t\t%t\n",
				a.ID, a.Name, a.Type,
				client.MilliunitsToAmount(a.Balance),
				a.OnBudget, a.Closed)
		}

	case *client.Account:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "ID\t%s\n", v.ID)
		fmt.Fprintf(w, "Name\t%s\n", v.Name)
		fmt.Fprintf(w, "Type\t%s\n", v.Type)
		fmt.Fprintf(w, "Balance\t%.2f\n", client.MilliunitsToAmount(v.Balance))
		fmt.Fprintf(w, "Cleared Balance\t%.2f\n", client.MilliunitsToAmount(v.ClearedBalance))
		fmt.Fprintf(w, "Uncleared Balance\t%.2f\n", client.MilliunitsToAmount(v.UnclearedBalance))
		fmt.Fprintf(w, "On Budget\t%t\n", v.OnBudget)
		fmt.Fprintf(w, "Closed\t%t\n", v.Closed)
		if v.Note != "" {
			fmt.Fprintf(w, "Note\t%s\n", v.Note)
		}

	case []client.CategoryGroup:
		fmt.Fprintln(w, "GROUP\tCATEGORY\tBUDGETED\tACTIVITY\tBALANCE")
		for _, g := range v {
			if g.Deleted || g.Hidden {
				continue
			}
			for _, c := range g.Categories {
				if c.Deleted || c.Hidden {
					continue
				}
				fmt.Fprintf(w, "%s\t%s\t%.2f\t%.2f\t%.2f\n",
					g.Name, c.Name,
					client.MilliunitsToAmount(c.Budgeted),
					client.MilliunitsToAmount(c.Activity),
					client.MilliunitsToAmount(c.Balance))
			}
		}

	case *client.Category:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "ID\t%s\n", v.ID)
		fmt.Fprintf(w, "Name\t%s\n", v.Name)
		fmt.Fprintf(w, "Group\t%s\n", v.CategoryGroupName)
		fmt.Fprintf(w, "Budgeted\t%.2f\n", client.MilliunitsToAmount(v.Budgeted))
		fmt.Fprintf(w, "Activity\t%.2f\n", client.MilliunitsToAmount(v.Activity))
		fmt.Fprintf(w, "Balance\t%.2f\n", client.MilliunitsToAmount(v.Balance))
		if v.GoalType != "" {
			fmt.Fprintf(w, "Goal Type\t%s\n", v.GoalType)
			fmt.Fprintf(w, "Goal Target\t%.2f\n", client.MilliunitsToAmount(v.GoalTarget))
		}
		if v.Note != "" {
			fmt.Fprintf(w, "Note\t%s\n", v.Note)
		}

	case []client.Transaction:
		fmt.Fprintln(w, "DATE\tPAYEE\tCATEGORY\tMEMO\tAMOUNT\tCLEARED")
		for _, t := range v {
			if t.Deleted {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.2f\t%s\n",
				t.Date, t.PayeeName, t.CategoryName,
				truncate(t.Memo, 30),
				client.MilliunitsToAmount(t.Amount), t.Cleared)
		}

	case *client.Transaction:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "ID\t%s\n", v.ID)
		fmt.Fprintf(w, "Date\t%s\n", v.Date)
		fmt.Fprintf(w, "Amount\t%.2f\n", client.MilliunitsToAmount(v.Amount))
		fmt.Fprintf(w, "Payee\t%s\n", v.PayeeName)
		fmt.Fprintf(w, "Category\t%s\n", v.CategoryName)
		fmt.Fprintf(w, "Account\t%s\n", v.AccountName)
		fmt.Fprintf(w, "Cleared\t%s\n", v.Cleared)
		fmt.Fprintf(w, "Approved\t%t\n", v.Approved)
		if v.Memo != "" {
			fmt.Fprintf(w, "Memo\t%s\n", v.Memo)
		}
		if v.FlagColor != "" {
			fmt.Fprintf(w, "Flag\t%s\n", v.FlagColor)
		}

	case []client.Payee:
		fmt.Fprintln(w, "ID\tNAME\tTRANSFER ACCOUNT")
		for _, p := range v {
			if p.Deleted {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.ID, p.Name, p.TransferAccountID)
		}

	case *client.Payee:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "ID\t%s\n", v.ID)
		fmt.Fprintf(w, "Name\t%s\n", v.Name)
		if v.TransferAccountID != "" {
			fmt.Fprintf(w, "Transfer Account ID\t%s\n", v.TransferAccountID)
		}

	case []client.ScheduledTransaction:
		fmt.Fprintln(w, "DATE NEXT\tFREQUENCY\tPAYEE\tCATEGORY\tAMOUNT")
		for _, st := range v {
			if st.Deleted {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.2f\n",
				st.DateNext, st.Frequency, st.PayeeName, st.CategoryName,
				client.MilliunitsToAmount(st.Amount))
		}

	case *client.ScheduledTransaction:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "ID\t%s\n", v.ID)
		fmt.Fprintf(w, "Date First\t%s\n", v.DateFirst)
		fmt.Fprintf(w, "Date Next\t%s\n", v.DateNext)
		fmt.Fprintf(w, "Frequency\t%s\n", v.Frequency)
		fmt.Fprintf(w, "Amount\t%.2f\n", client.MilliunitsToAmount(v.Amount))
		fmt.Fprintf(w, "Payee\t%s\n", v.PayeeName)
		fmt.Fprintf(w, "Category\t%s\n", v.CategoryName)
		fmt.Fprintf(w, "Account\t%s\n", v.AccountName)
		if v.Memo != "" {
			fmt.Fprintf(w, "Memo\t%s\n", v.Memo)
		}

	case []client.Month:
		fmt.Fprintln(w, "MONTH\tINCOME\tBUDGETED\tACTIVITY\tTO BE BUDGETED")
		for _, m := range v {
			if m.Deleted {
				continue
			}
			fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\n",
				m.Month,
				client.MilliunitsToAmount(m.Income),
				client.MilliunitsToAmount(m.Budgeted),
				client.MilliunitsToAmount(m.Activity),
				client.MilliunitsToAmount(m.ToBeBudgeted))
		}

	case *client.Month:
		fmt.Fprintln(w, "FIELD\tVALUE")
		fmt.Fprintf(w, "Month\t%s\n", v.Month)
		fmt.Fprintf(w, "Income\t%.2f\n", client.MilliunitsToAmount(v.Income))
		fmt.Fprintf(w, "Budgeted\t%.2f\n", client.MilliunitsToAmount(v.Budgeted))
		fmt.Fprintf(w, "Activity\t%.2f\n", client.MilliunitsToAmount(v.Activity))
		fmt.Fprintf(w, "To Be Budgeted\t%.2f\n", client.MilliunitsToAmount(v.ToBeBudgeted))
		if v.AgeOfMoney > 0 {
			fmt.Fprintf(w, "Age of Money\t%d days\n", v.AgeOfMoney)
		}
		if v.Note != "" {
			fmt.Fprintf(w, "Note\t%s\n", v.Note)
		}

	default:
		// Fall back to JSON for unknown types
		return f.printJSON(data)
	}

	return nil
}

// truncate shortens a string to the given length
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
