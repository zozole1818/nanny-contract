package zus

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/shopspring/decimal"
	"sort"
)

type Summary struct {
	Name              string
	GrossIncome       decimal.Decimal
	NetIncome         decimal.Decimal
	EmployerTotalCost decimal.Decimal
	//ZUSAmount         decimal.Decimal
	Zus Zus

	PaidByEmployer map[string]Record
	PaidByCountry  map[string]Record
}

type Zus struct {
	EmployeeFund decimal.Decimal
	EmployerFund decimal.Decimal

	EmployerPaid decimal.Decimal
	CountryPaid  decimal.Decimal
}

type Record struct {
	Name        string
	Value       decimal.Decimal
	Records     map[string]Record
	Description string
}

func NewSummary(dra Dra) Summary { // todo: polish names by default - add other languages
	summary := Summary{}
	if !dra.Calculated() {
		dra = dra.Calculate()
	}

	grossIncome := decimal.Zero
	for _, rca := range dra.Rcas {
		grossIncome = grossIncome.Add(rca.Base)
	}

	netIncome := grossIncome
	zusEmployeeFund := decimal.Zero
	for _, v := range dra.contributionFundByEmployee {
		netIncome = netIncome.Sub(v)
		zusEmployeeFund = zusEmployeeFund.Add(v)
	}

	paidByEmployer := make(map[string]Record)
	for k, v := range dra.contributionFundByEmployee {
		key := paidByPolMap[employeePaidBy]
		paidByEmployer[string(k)] = Record{
			Name:  contributionNamePolMap[k],
			Value: v,
			Records: map[string]Record{
				key: {
					Name:        key,
					Value:       v,
					Description: "",
				},
			},
			Description: "",
		}
	}
	zusEmployerFund := decimal.Zero
	for k, v := range dra.contributionFundByEmployer {
		zusEmployerFund = zusEmployerFund.Add(v)
		tmp, ok := paidByEmployer[string(k)]
		if ok {
			key := paidByPolMap[employerPaidBy]
			tmp.Records[key] = Record{
				Name:        key,
				Value:       v,
				Description: "",
			}
			paidByEmployer[string(k)] = Record{
				Name:        contributionNamePolMap[k],
				Value:       tmp.Value.Add(v),
				Records:     tmp.Records,
				Description: "",
			}

		}
	}

	paidByCountry := make(map[string]Record, len(dra.contributionFundByCountry))
	for k, v := range dra.contributionFundByCountry {
		p1Key := contributionNamePolMap[k]
		p2Key := paidByPolMap[countryPaidBy]
		key := p1Key + " " + p2Key
		paidByCountry[p1Key] = Record{
			Name:        key,
			Value:       v,
			Description: "",
		}
	}

	summary.PaidByCountry = paidByCountry
	summary.PaidByEmployer = paidByEmployer

	summary.Name = fmt.Sprintf("Podsumowanie za miesiąc: %s", monthsPolMap[dra.SettlementPeriod.String()])
	summary.GrossIncome = grossIncome
	summary.NetIncome = netIncome
	summary.EmployerTotalCost = netIncome.Add(dra.zusPaidByEmployerSum)
	summary.Zus = Zus{
		EmployeeFund: zusEmployeeFund,
		EmployerFund: zusEmployerFund,

		EmployerPaid: dra.zusPaidByEmployerSum,
		CountryPaid:  dra.zusPaidByCountrySum,
	}

	return summary
}

func (s Summary) Print() {
	fmt.Println("Summary")
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "NAME"},
			{Align: simpletable.AlignRight, Text: "VALUE"},
		},
	}
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "kwota brutto")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.GrossIncome.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "kwota netto")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.NetIncome.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "całkowity koszt płatnika")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.EmployerTotalCost.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "ZUS finansowany przez nianie")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.Zus.EmployeeFund.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "ZUS finansowany przez rodzica")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.Zus.EmployerFund.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "kwota składek ZUS (do zapłaty przez płatnika)")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.Zus.EmployerPaid.StringFixed(2))},
	})
	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", "kwota składek ZUS (do zapłaty przez państwo)")},
		{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.Zus.CountryPaid.StringFixed(2))},
	})

	keys := make([]string, 0, len(s.PaidByEmployer))
	for k := range s.PaidByEmployer {
		keys = append(keys, k)
	}

	// Sort keys
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		value := s.PaidByEmployer[k]
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", value.Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", value.Value.StringFixed(2))},
		})
		if value.Records != nil {
			for _, v2 := range value.Records {
				table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
					{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", v2.Name)},
					{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", v2.Value.StringFixed(2))},
				})
			}
		}
	}

	keys2 := make([]string, 0, len(s.PaidByCountry))
	for k := range s.PaidByCountry {
		keys2 = append(keys2, k)
	}

	// Sort keys
	sort.Slice(keys2, func(i, j int) bool {
		return keys2[i] > keys2[j]
	})
	for _, k := range keys2 {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", k)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%s", s.PaidByCountry[k].Value.StringFixed(2))},
		})
	}
	fmt.Println(table.String())
}
