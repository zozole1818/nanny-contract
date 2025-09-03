package zus

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestRca_Print(t *testing.T) {
	type fields struct {
		Code InsuranceCode
		Base decimal.Decimal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "#2333",
			fields: fields{
				Code: Code043000,
				Base: decimal.NewFromInt(2333),
			},
		},
		{
			name: "#1367",
			fields: fields{
				Code: Code043100,
				Base: decimal.NewFromInt(1367),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rca := NewRca(time.August, tt.fields.Code, tt.fields.Base)
			rca.Calculate()
			rca.Print()
		})
	}
}
