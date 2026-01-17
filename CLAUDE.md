# ynabctl - Claude Instructions

## Testing

**IMPORTANT**: Always use the **Test** budget for testing API calls, never the production budget.

- Test budget ID: `bea83a82-bf56-40ea-a482-817cdf84d546`
- Test budget name: "Test"

Use `--budget bea83a82-bf56-40ea-a482-817cdf84d546` when testing commands.

## Budgets

| Budget | ID | Usage |
|--------|-----|-------|
| Langtind v4 | `a26490bd-8540-48be-b9af-635d1fa2c223` | Production - DO NOT use for testing |
| Test | `bea83a82-bf56-40ea-a482-817cdf84d546` | Testing only |

## YNAB API Notes

- API uses `date` for SaveScheduledTransaction, not `date_first`
- `date_first` is only in response objects
- Amounts are in milliunits (1000 = $1.00)
