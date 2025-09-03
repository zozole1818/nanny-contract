package zus

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

//type Generator interface {
//	NannyDRA(minimalWage decimal.Decimal, grossIncome decimal.Decimal, sicknessContribution bool) Dra
//}

type Generator struct {
}

func (g Generator) NannyDRA(month time.Month, minimalWage decimal.Decimal, grossIncome decimal.Decimal, sicknessContribution bool) (Dra, error) {
	if minimalWage.LessThanOrEqual(decimal.Zero) {
		return Dra{}, fmt.Errorf("minimal wage must be greater than zero")
	}
	if !sicknessContribution {
		return Dra{}, fmt.Errorf("sickness contribution false is not supported yet")
	}
	countryBase := minimalWage.DivRound(decimal.NewFromInt(2), 2)
	employerBase := grossIncome.Sub(countryBase)
	var dra *Dra
	if employerBase.LessThanOrEqual(decimal.Zero) {
		dra = NewDra([]Rca{
			*NewRca(month, Code043000, countryBase),
		})
	} else {
		dra = NewDra([]Rca{
			*NewRca(month, Code043000, countryBase),
			*NewRca(month, Code043100, employerBase),
		})
	}
	dra.Calculate()
	return *dra, nil
}
