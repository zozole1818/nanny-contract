package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/zozole1818/nanny-contract/internal/email"
	"github.com/zozole1818/nanny-contract/internal/misc"
	"github.com/zozole1818/nanny-contract/internal/reports/zus"
	"log/slog"
	"time"
)

var draCmd = &cobra.Command{
	Use:   "dra",
	Short: "Help to calculate ZUS DRA",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		month, minimalWage, base, createSummary, emailAddr, err := parseDRAFlags(cmd)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		dra, err := zus.Generator{}.NannyDRA(month, minimalWage, base, true)
		dra.Print()
		var summary zus.Summary
		if createSummary {
			summary = dra.CreateSummary()
			summary.Print()
		}
		if emailAddr != "" {
			err := misc.LoadEnvs(".env")
			if err != nil {
				slog.Error("Error when loading env file. Won't send email.", "error", err)
				return
			}
			if summary.Name == "" {
				summary = dra.CreateSummary()
			}
			emailer := email.NewEmailer("smtp.gmail.com", 587, "zuzanna.go.mail@gmail.com")
			err = emailer.Send([]string{emailAddr}, summary)
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(draCmd)

	draCmd.Flags().IntP("month", "m", 1, "Month for which to calculate DRA. Numeric value from 1 to 12.")
	draCmd.Flags().IntP("base", "b", 100, "Base [zł] for which to calculate DRA.")
	draCmd.Flags().IntP("minimalWage", "w", 4666, "Minimal wage in Poland [zł].")
	draCmd.Flags().BoolP("summary", "s", false, "Print summary.")
	draCmd.Flags().StringP("email", "e", "", "Send email to this address with summary.")
}

func parseDRAFlags(cmd *cobra.Command) (time.Month, decimal.Decimal, decimal.Decimal, bool, string, error) {
	month, err := cmd.Flags().GetInt("month")
	if err != nil {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("error when getting month: %v", err)
	}
	if month < 1 || month > 12 {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("month must be between 1 and 12")
	}
	base, err := cmd.Flags().GetInt("base")
	if err != nil {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("error when getting base: %v", err)
	}
	if base <= 0 {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("base must be greater than 0")
	}
	minimalWage, err := cmd.Flags().GetInt("minimalWage")
	if err != nil {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("error when getting minimalWage: %v", err)
	}
	if minimalWage <= 0 {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("minimalWage must be greater than 0")
	}
	summary, err := cmd.Flags().GetBool("summary")
	if err != nil {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("error when getting summary: %v", err)
	}
	emailAddress, err := cmd.Flags().GetString("email")
	if err != nil {
		return 0, decimal.Zero, decimal.Zero, false, "", fmt.Errorf("error when getting email address: %v", err)
	}
	return time.Month(month), decimal.NewFromInt(int64(minimalWage)), decimal.NewFromInt(int64(base)), summary, emailAddress, nil
}
