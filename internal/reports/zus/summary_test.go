package zus

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestSummary_Print(t *testing.T) {
	type fields struct {
		Rca []Rca
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "#basic",
			fields: fields{
				Rca: []Rca{
					*NewRca(time.August, Code043000, decimal.NewFromInt(2333)),
					*NewRca(time.August, Code043100, decimal.NewFromInt(100)),
				},
			},
		},
		{
			name: "#mine",
			fields: fields{
				Rca: []Rca{
					*NewRca(time.August, Code043000, decimal.NewFromInt(2333)),
					*NewRca(time.August, Code043100, decimal.NewFromInt(1367)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSummary(*NewDra(tt.fields.Rca))
			s.Print()
		})
	}
}
