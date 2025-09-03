package zus

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

type Dra struct {
	SettlementPeriod time.Month
	Rcas             []Rca

	contributionFundByEmployee map[ContributionName]decimal.Decimal
	contributionFundByEmployer map[ContributionName]decimal.Decimal
	contributionFundByCountry  map[ContributionName]decimal.Decimal

	pensionSum  decimal.Decimal
	rentalSum   decimal.Decimal
	sickSum     decimal.Decimal
	accidentSum decimal.Decimal

	zusPaidByEmployerSum decimal.Decimal
	zusPaidByCountrySum  decimal.Decimal

	calculated bool
}

func NewDra(rcas []Rca) *Dra {
	month := time.January
	if len(rcas) == 0 {
		return &Dra{
			SettlementPeriod: month,
			Rcas:             rcas,
		}
	}
	return &Dra{
		SettlementPeriod: rcas[0].SettlementPeriod,
		Rcas:             rcas,
	}
}

func (d *Dra) Validate() error {
	if len(d.Rcas) == 0 {
		return nil
	}
	var month = d.Rcas[0].SettlementPeriod
	for _, rca := range d.Rcas {
		if rca.SettlementPeriod != month {
			return fmt.Errorf("settlement period must be the same for all rcas")
		}
	}
	return nil
}

//func (d *Dra) GetGrossIncome() decimal.Decimal {
//	sum := decimal.Zero
//	for _, rca := range d.Rcas {
//		sum = sum.Add(rca.Base)
//	}
//	return sum
//}

func (d *Dra) Calculated() bool {
	return d.calculated
}

func (d *Dra) Calculate() Dra {
	fbEmployee := make(map[ContributionName]decimal.Decimal, 4)
	fbEmployer := make(map[ContributionName]decimal.Decimal, 4)
	fbCountry := make(map[ContributionName]decimal.Decimal, 4)

	copyRcas := make([]Rca, len(d.Rcas))
	copy(copyRcas, d.Rcas)
	for i, rca := range copyRcas {
		if !rca.Calculated() {
			rca = rca.Calculate()
			d.Rcas[i] = rca
		}

		fbEmployee[Pension] = fbEmployee[Pension].Add(rca.contributionPaidByEmployee[Pension])
		fbEmployee[Rental] = fbEmployee[Rental].Add(rca.contributionPaidByEmployee[Rental])
		fbEmployee[Sick] = fbEmployee[Sick].Add(rca.contributionPaidByEmployee[Sick])
		fbEmployee[Accident] = fbEmployee[Accident].Add(rca.contributionPaidByEmployee[Accident])

		fbEmployer[Pension] = fbEmployer[Pension].Add(rca.contributionPaidByEmployer[Pension])
		fbEmployer[Rental] = fbEmployer[Rental].Add(rca.contributionPaidByEmployer[Rental])
		fbEmployer[Sick] = fbEmployer[Sick].Add(rca.contributionPaidByEmployer[Sick])
		fbEmployer[Accident] = fbEmployer[Accident].Add(rca.contributionPaidByEmployer[Accident])

		fbCountry[Pension] = fbCountry[Pension].Add(rca.contributionPaidByCountry[Pension])
		fbCountry[Rental] = fbCountry[Rental].Add(rca.contributionPaidByCountry[Rental])
		fbCountry[Sick] = fbCountry[Sick].Add(rca.contributionPaidByCountry[Sick])
		fbCountry[Accident] = fbCountry[Accident].Add(rca.contributionPaidByCountry[Accident])

		fbEmployee[Health] = fbEmployee[Health].Add(rca.contributionPaidByEmployee[Health])
		fbEmployer[Health] = fbEmployer[Health].Add(rca.contributionPaidByEmployer[Health]) // should be always 0.00
		fbCountry[Health] = fbCountry[Health].Add(rca.contributionPaidByCountry[Health])
	}

	d.contributionFundByEmployee = fbEmployee
	d.contributionFundByEmployer = fbEmployer
	d.contributionFundByCountry = fbCountry

	d.pensionSum = fbEmployee[Pension].Add(fbEmployer[Pension]).Add(fbCountry[Pension])
	d.rentalSum = fbEmployee[Rental].Add(fbEmployer[Rental]).Add(fbCountry[Rental])
	d.sickSum = fbEmployee[Sick].Add(fbEmployer[Sick]).Add(fbCountry[Sick])
	d.accidentSum = fbCountry[Accident].Add(fbEmployer[Accident]).Add(fbCountry[Accident])

	for _, v := range fbEmployee {
		d.zusPaidByEmployerSum = d.zusPaidByEmployerSum.Add(v)
	}
	for _, v := range fbEmployer {
		d.zusPaidByEmployerSum = d.zusPaidByEmployerSum.Add(v)
	}
	for _, v := range fbCountry {
		d.zusPaidByCountrySum = d.zusPaidByCountrySum.Add(v)
	}

	d.calculated = true
	return Dra{
		Rcas:                       d.Rcas,
		contributionFundByEmployee: d.contributionFundByEmployee,
		contributionFundByEmployer: d.contributionFundByEmployer,
		contributionFundByCountry:  d.contributionFundByCountry,

		pensionSum:  d.pensionSum,
		rentalSum:   d.rentalSum,
		sickSum:     d.sickSum,
		accidentSum: d.accidentSum,

		zusPaidByEmployerSum: d.zusPaidByEmployerSum,
		zusPaidByCountrySum:  d.zusPaidByCountrySum,
		calculated:           d.calculated,
	}
}

func (d *Dra) CreateSummary() Summary {
	return NewSummary(*d)
}

func (d *Dra) Print() {

	for _, rca := range d.Rcas {
		fmt.Println("RCA")
		rca.Print()
	}

	fmt.Println("DRA")
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "#"},
			{Align: simpletable.AlignRight, Text: "ubz. emerytalne"},
			{Align: simpletable.AlignRight, Text: "ubz. rentowe"},
			{Align: simpletable.AlignRight, Text: "ubz. chorobowe"},
			{Align: simpletable.AlignRight, Text: "ubz. wypadkowe"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "suma")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.pensionSum.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.rentalSum.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.sickSum.String())},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.accidentSum.String())},
	})

	rows := map[string]map[ContributionName]decimal.Decimal{
		"fin. przez ubezpieczonego": d.contributionFundByEmployee,
		"fin. przez płatnika":       d.contributionFundByEmployer,
		"fin. przez państwo":        d.contributionFundByCountry,
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
		if k != "fin. przez państwo" {
			sum = sum.Add(rows[k][Pension]).Add(rows[k][Rental]).Add(rows[k][Sick]).Add(rows[k][Accident])
		}
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{},
			{Align: simpletable.AlignRight, Text: "Łącznie do zapłaty przez płatnika"},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", sum.String())},
		},
	}

	fmt.Println(table.String())

	healthTable := simpletable.New()
	healthTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("#")},
			{Align: simpletable.AlignRight, Text: "ubz. zdrowotne"},
		},
	}

	healthTable.Body.Cells = append(healthTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "fin. przez ubezpieczonego")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.contributionFundByEmployee[Health].String())},
	})

	healthTable.Body.Cells = append(healthTable.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "fin. przez państwo")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", d.contributionFundByCountry[Health].String())},
	})

	//healthTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println(healthTable.String())
}
