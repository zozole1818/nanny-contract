package zus

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/shopspring/decimal"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"sort"
	"time"
)

type InsuranceCode string

var (
	Code043000 InsuranceCode = "043000" // calculation for half of minimal wage
	Code043100 InsuranceCode = "043100" // calculation for above half of minimal wage
)

type NewResponse struct {
	Name                 string
	GrossIncome          float64
	NetIncome            float64
	SicknessContribution bool
	EmployerTotalCost    float64
	ZUSAmount            float64
}

type ContributionName string

var (
	Rental   ContributionName = "rental"
	Pension  ContributionName = "pension"
	Sick     ContributionName = "sick"
	Accident ContributionName = "accident"
	Health   ContributionName = "health"
)

type Rca struct {
	SettlementPeriod time.Month
	InsuranceCode    InsuranceCode
	Base             decimal.Decimal

	baseForPension  decimal.Decimal
	baseForRental   decimal.Decimal
	baseForSick     decimal.Decimal
	baseForAccident decimal.Decimal
	baseForHealth   decimal.Decimal

	contributionPaidByEmployee map[ContributionName]decimal.Decimal
	contributionPaidByEmployer map[ContributionName]decimal.Decimal
	contributionPaidByCountry  map[ContributionName]decimal.Decimal

	calculated bool
}

func NewRca(settlementPeriod time.Month, insuranceCode InsuranceCode, base decimal.Decimal) *Rca {
	return &Rca{
		SettlementPeriod: settlementPeriod,
		InsuranceCode:    insuranceCode,
		Base:             base,

		baseForPension:  base,
		baseForRental:   base,
		baseForSick:     base,
		baseForAccident: base,
	}
}

func (r *Rca) Calculate() Rca {
	pbEmployee := make(map[ContributionName]decimal.Decimal, 4)
	pbEmployer := make(map[ContributionName]decimal.Decimal, 4)
	pbCountry := make(map[ContributionName]decimal.Decimal, 4)
	switch r.InsuranceCode {
	case Code043000:
		pbEmployee[Pension] = decimal.Zero
		pbEmployee[Rental] = decimal.Zero
		pbEmployee[Sick] = calculatePercentage(r.baseForSick, model.SicknessContribution)
		pbEmployee[Accident] = decimal.Zero

		pbEmployer[Pension] = decimal.Zero
		pbEmployer[Rental] = decimal.Zero
		pbEmployer[Sick] = decimal.Zero
		pbEmployer[Accident] = decimal.Zero

		pbCountry[Pension] = calculatePercentage(r.baseForPension, model.PensionContribution)
		pbCountry[Rental] = calculatePercentage(r.baseForRental, model.SocialPensionContribution)
		pbCountry[Sick] = decimal.Zero
		pbCountry[Accident] = calculatePercentage(r.baseForAccident, model.AccidentInsuranceContribution)
	case Code043100:
		pbEmployee[Pension] = calculatePercentage(r.baseForPension, model.PensionContributionPaidByEmployee)
		pbEmployee[Rental] = calculatePercentage(r.baseForRental, model.SocialPensionContributionPaidByEmployee)
		pbEmployee[Sick] = calculatePercentage(r.baseForSick, model.SicknessContribution)
		pbEmployee[Accident] = calculatePercentage(r.baseForAccident, model.AccidentInsuranceContributionPaidByEmployee)

		pbEmployer[Pension] = calculatePercentage(r.baseForPension, model.PensionContributionPaidByEmployer)
		pbEmployer[Rental] = calculatePercentage(r.baseForRental, model.SocialPensionContributionPaidByEmployer)
		pbEmployer[Sick] = decimal.Zero
		pbEmployer[Accident] = calculatePercentage(r.baseForAccident, model.AccidentInsuranceContributionPaidByEmployer)

		pbCountry[Pension] = decimal.Zero
		pbCountry[Rental] = decimal.Zero
		pbCountry[Sick] = decimal.Zero
		pbCountry[Accident] = decimal.Zero

	}

	r.baseForHealth = r.baseForPension.Sub(pbEmployee[Pension]).Sub(pbEmployee[Rental]).Sub(pbEmployee[Sick]).Sub(pbEmployee[Accident])
	switch r.InsuranceCode {
	case Code043000:
		pbEmployee[Health] = decimal.Zero
		pbEmployer[Health] = decimal.Zero
		pbCountry[Health] = calculatePercentage(r.baseForHealth, model.HealthInsuranceContribution)
	case Code043100:
		pbEmployee[Health] = calculatePercentage(r.baseForHealth, model.HealthInsuranceContribution)
		pbEmployer[Health] = decimal.Zero
		pbCountry[Health] = decimal.Zero
	}

	r.contributionPaidByEmployee = pbEmployee
	r.contributionPaidByEmployer = pbEmployer
	r.contributionPaidByCountry = pbCountry

	r.calculated = true
	return *r
}

func (r *Rca) Calculated() bool {
	return r.calculated
}

func calculatePercentage(base decimal.Decimal, percentage float64) decimal.Decimal {
	return base.Mul(decimal.NewFromFloat(percentage)).DivRound(decimal.NewFromInt(100), 2)
}

func (r *Rca) Print() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("code[%s]", r.InsuranceCode)},
			{Align: simpletable.AlignRight, Text: "ubz. emerytalne"},
			{Align: simpletable.AlignRight, Text: "ubz. rentowe"},
			{Align: simpletable.AlignRight, Text: "ubz. chorobowe"},
			{Align: simpletable.AlignRight, Text: "ubz. wypadkowe"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "podstawa")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", r.baseForPension.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", r.baseForRental.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", r.baseForSick.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", r.baseForAccident.String())},
	})

	rows := map[string]map[ContributionName]decimal.Decimal{
		"fin. przez ubezpieczonego": r.contributionPaidByEmployee,
		"fin. przez płatnika":       r.contributionPaidByEmployer,
		"fin. przez państwo":        r.contributionPaidByCountry,
	}

	keys := make([]string, 0, len(rows))
	for k := range rows {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	sum := decimal.Zero
	for _, k := range keys {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", k)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", rows[k][Pension].String())},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", rows[k][Rental].String())},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", rows[k][Sick].String())},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", rows[k][Accident].String())},
		}

		table.Body.Cells = append(table.Body.Cells, r)
		sum = sum.Add(rows[k][Pension]).Add(rows[k][Rental]).Add(rows[k][Sick]).Add(rows[k][Accident])
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{},
			{Align: simpletable.AlignRight, Text: "Łącznie"},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", sum.String())},
		},
	}

	//table.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(table.String())

	healthTable := simpletable.New()
	healthTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("code[%s]", r.InsuranceCode)},
			{Align: simpletable.AlignRight, Text: "ubz. zdrowotne"},
		},
	}

	healthTable.Body.Cells = append(healthTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "podstawa")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", r.baseForHealth.String())},
	})
	title := ""
	healthContributionValue := decimal.Decimal{}
	switch r.InsuranceCode {
	case Code043000:
		title = "fin. przez państwo"
		healthContributionValue = r.contributionPaidByCountry[Health]
	case Code043100:
		title = "fin. przez ubezpieczonego"
		healthContributionValue = r.contributionPaidByEmployee[Health]
	}

	healthTable.Body.Cells = append(healthTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", title)},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", healthContributionValue.String())},
	})

	//healthTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(healthTable.String())

}
