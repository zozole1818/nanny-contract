package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/zozole1818/nanny-contract/internal/reports/zus"
	"log/slog"
	"time"
)

var rcaCmd = &cobra.Command{
	Use:   "rca",
	Short: "Help to calculate ZUS RCA",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		month, code, base, err := parseRCAFlags(cmd)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		rca := zus.NewRca(month, code, base)
		rca.Calculate()
		rca.Print()
	},
}

func init() {
	rootCmd.AddCommand(rcaCmd)

	rcaCmd.Flags().IntP("month", "m", 1, "Month for which to calculate RCA. Numeric value from 1 to 12.")
	rcaCmd.Flags().StringP("code", "c", "043000", "One of the ZUS codes [043000 or 043100].")
	rcaCmd.Flags().IntP("base", "b", 100, "Base [zł] for which to calculate RCA.")
}

func parseRCAFlags(cmd *cobra.Command) (time.Month, zus.InsuranceCode, decimal.Decimal, error) {
	month, err := cmd.Flags().GetInt("month")
	if err != nil {
		return 0, "", decimal.Zero, fmt.Errorf("error when getting month: %v", err)
	}
	if month < 1 || month > 12 {
		return 0, "", decimal.Zero, fmt.Errorf("month must be between 1 and 12")
	}
	code, err := cmd.Flags().GetString("code")
	if err != nil {
		return 0, "", decimal.Zero, fmt.Errorf("error when getting code: %v", err)
	}
	if code != "043000" && code != "043100" {
		return 0, "", decimal.Zero, fmt.Errorf("code must be 043000 or 043100")
	}
	base, err := cmd.Flags().GetInt("base")
	if err != nil {
		return 0, "", decimal.Zero, fmt.Errorf("error when getting base: %v", err)
	}
	if base <= 0 {
		return 0, "", decimal.Zero, fmt.Errorf("base must be greater than 0")
	}
	return time.Month(month), zus.InsuranceCode(code), decimal.NewFromInt(int64(base)), nil
}
