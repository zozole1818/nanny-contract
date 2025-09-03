package model

import (
	"fmt"
)

type ReportRequest struct {
	MinimumWageInPoland float64 `xml:"minimum_wage_in_poland" json:"minimum_wage_in_poland,omitempty"`
	// NumberOfWorkingHours int     `xml:"number_of_working_hours" json:"number_of_working_hours,omitempty"`
	// GrossPerHour         float64 `xml:"gross_per_hour" json:"gross_per_hour,omitempty"`
	GrossIncome          float64 `xml:"gross_income" json:"gross_income,omitempty"`
	SicknessContribution bool    `xml:"sickness_contribution" json:"sickness_contribution,omitempty"`
}

func (r ReportRequest) Validate() error {
	if r.GrossIncome == 0 {
		// if r.NumberOfWorkingHours == 0 || r.GrossPerHour == 0 {
		// 	return fmt.Errorf("cannot calculate gross income")
		// }
		return fmt.Errorf("missing gross income")
	}
	return nil
}

type Contribution struct {
	Name       string  `json:"name,omitempty"`
	Value      float64 `json:"value"`
	Base       float64 `json:"base,omitempty"`
	Percentage float64 `json:"percentage"`
}

type ReportResponse struct {
	Name                                            string         `json:"name,omitempty"`
	NumberOfWorkingHours                            int64          `json:"number_of_working_hours,omitempty"`
	GrossPerHour                                    float64        `json:"gross_per_hour,omitempty"`
	GrossIncome                                     float64        `json:"gross_income,omitempty"`
	NetIncome                                       float64        `json:"net_income,omitempty"`
	SicknessContribution                            bool           `json:"sickness_contribution,omitempty"`
	BaseForCalculationOfContributionsPaidByCountry  float64        `json:"base_for_calculation_of_contributions_paid_by_country,omitempty"`
	BaseForCalculationOfContributionsPaidByEmployer float64        `json:"base_for_calculation_of_contributions_paid_by_employer,omitempty"`
	ContributionPaidByEmployer                      []Contribution `json:"contribution_paid_by_employer,omitempty"`
	ContributionPaidByEmployee                      []Contribution `json:"contribution_paid_by_employee,omitempty"`
	ContributionPaidByCountry                       []Contribution `json:"contribution_paid_by_country,omitempty"`
	EmployerTotalCost                               float64        `json:"employer_total_cost,omitempty"`
	ZUSAmount                                       float64        `json:"zus_amount,omitempty"`
}

type ReportResponseView struct {
	Name                  string                  `json:"name,omitempty"`
	GrossIncome           float64                 `json:"gross_income,omitempty"`
	NetIncome             float64                 `json:"net_income,omitempty"`
	SicknessContribution  bool                    `json:"sickness_contribution,omitempty"`
	PaidByEmployerCompact []PaidByEmployerCompact `json:"paid_by_employer_compact,omitempty"`
	//ContributionPaidByEmployer []Contribution          `json:"contribution_paid_by_employer,omitempty"`
	//ContributionPaidByEmployee []Contribution          `json:"contribution_paid_by_employee,omitempty"`
	ContributionPaidByCountry []Contribution `json:"contribution_paid_by_country,omitempty"`
	EmployerTotalCost         float64        `json:"employer_total_cost,omitempty"`
	ZUSAmount                 float64        `json:"zus_amount,omitempty"`
}

type PaidByEmployerCompact struct {
	Name                  string  `json:"name,omitempty"`
	Value                 float64 `json:"value"`
	ToolTip               string  `json:"tooltip,omitempty"`
	PaidByEmployer        float64 `json:"paid_by_employer"`
	PaidByEmployerToolTip string  `json:"paid_by_employer_tooltip,omitempty"`
	PaidByEmployee        float64 `json:"paid_by_employee"`
	PaidByEmployeeToolTip string  `json:"paid_by_employee_tooltip,omitempty"`
}

func (r ReportResponse) ToReportResponseView() ReportResponseView {
	var compact []PaidByEmployerCompact
	larger, smaller := func(arr1 []Contribution, arr2 []Contribution) ([]Contribution, []Contribution) {
		if len(arr1) > len(arr2) {
			return arr1, arr2
		}
		return arr2, arr1
	}(r.ContributionPaidByEmployee, r.ContributionPaidByEmployer)
	for _, c1 := range larger {
		tmp := PaidByEmployerCompact{}
		tmp.Name = c1.Name
		added := false
		for _, c2 := range smaller {
			if c1.Name == c2.Name {
				tmp.PaidByEmployee = c1.Value
				tmp.PaidByEmployeeToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł", c1.Percentage, c1.Base, c1.Value)

				tmp.PaidByEmployer = c2.Value
				tmp.PaidByEmployerToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł", c2.Percentage, c2.Base, c2.Value)

				tmp.Value = c1.Value + c2.Value
				tmp.ToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł", c1.Percentage+c2.Percentage, c2.Base, c1.Value+c2.Value)

				compact = append(compact, tmp)
				added = true
			}
		}
		if !added { // means that we did not find this contribution
			if tmp.Name == "składka chorobowa" {
				tmp.Value = c1.Value
				tmp.ToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł (w całości opłacana przez nianie)", c1.Percentage, c1.Base, c1.Value)

				tmp.PaidByEmployee = c1.Value
				tmp.PaidByEmployeeToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł", c1.Percentage, c1.Base, c1.Value)

				tmp.PaidByEmployer = 0.0
				//tmp.PaidByEmployerToolTip = fmt.Sprintf("%0.2f %% z %0.2fzł = %0.2fzł", 0.0, c2.Base, c2.Value)

				compact = append(compact, tmp)
			}
		}
	}

	return ReportResponseView{
		Name:                      r.Name,
		GrossIncome:               r.GrossIncome,
		NetIncome:                 r.NetIncome,
		SicknessContribution:      r.SicknessContribution,
		PaidByEmployerCompact:     compact,
		ContributionPaidByCountry: r.ContributionPaidByCountry,
		EmployerTotalCost:         r.EmployerTotalCost,
		ZUSAmount:                 r.ZUSAmount,
	}
}
