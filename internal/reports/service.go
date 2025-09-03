package reports

import (
	"context"
	"github.com/zozole1818/nanny-contract/internal/misc"
	"github.com/zozole1818/nanny-contract/internal/reports/cache"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"log/slog"
)

type ReportService interface {
	CreateReport(request model.ReportRequest) (model.ReportResponse, error)
}

type report struct {
	cache cache.ReportCache
}

func NewReportSvc(transformer model.Transformer) ReportService {
	return report{
		cache: cache.NewReportCache("", transformer),
	}
}

func (r report) CreateReport(request model.ReportRequest) (model.ReportResponse, error) {
	ctx := context.Background()
	if err := request.Validate(); err != nil {
		return model.ReportResponse{}, err
	}

	resp, err := r.cache.GetReport(ctx, request.GrossIncome, request.SicknessContribution)
	if err == nil {
		return resp, nil
	} else {
		slog.Warn("cache error", "error", err)
	}

	var (
		baseForCountryContributions  float64
		baseForEmployerContributions float64

		paidByCountry  = []model.Contribution{}
		paidByEmployer = []model.Contribution{}
		paidByEmployee = []model.Contribution{}

		sicPercentage float64
	)

	if request.MinimumWageInPoland == 0 {
		baseForCountryContributions = min(float64(model.BaseForCalculationOfContributionsPaidByCountry), request.GrossIncome)
	} else {
		baseForCountryContributions = min(request.MinimumWageInPoland/2, request.GrossIncome)
	}

	if request.SicknessContribution {
		sicPercentage = 2.45
	}

	paidByCountry = model.NewPayer(model.Country, baseForCountryContributions, 19.52, 8.0, 1.67, 9.0, sicPercentage).
		CalculateContributions(baseForCountryContributions)

	if request.GrossIncome > baseForCountryContributions {
		baseForEmployerContributions := request.GrossIncome - baseForCountryContributions
		paidByEmployer = model.NewPayer(model.Employer, baseForEmployerContributions, 9.76, 6.5, 1.67, 0.0, 0.0).
			CalculateContributions(baseForEmployerContributions)
		paidByEmployee = model.NewPayer(model.Employee, baseForEmployerContributions, 9.76, 1.5, 0.0, 9.0, sicPercentage).
			CalculateContributions(request.GrossIncome)
	} else {
		if request.SicknessContribution {
			paidByEmployee = append(paidByEmployee, model.NewPayer(model.Employee, baseForEmployerContributions, 9.76, 1.5, 0.0, 9.0, sicPercentage).CalulateSicknessContribution(request.GrossIncome))
		}
	}

	netIncome := misc.Round(request.GrossIncome - sum(paidByEmployee))
	zusAmout := misc.Round(sum(paidByEmployer) + sum(paidByEmployee))
	employerTotalCost := misc.Round(netIncome + zusAmout)

	resp = model.ReportResponse{
		Name:                       "Umowa Aktywacyjna dla niani",
		GrossIncome:                request.GrossIncome,
		NetIncome:                  netIncome,
		SicknessContribution:       request.SicknessContribution,
		ContributionPaidByCountry:  paidByCountry,
		ContributionPaidByEmployer: paidByEmployer,
		ContributionPaidByEmployee: paidByEmployee,
		EmployerTotalCost:          employerTotalCost,
		ZUSAmount:                  zusAmout,
	}

	err = r.cache.AddReport(ctx, resp)
	if err != nil {
		slog.Warn("cache error:", "error", err)
	}

	return resp, nil
}

func sum(s []model.Contribution) float64 {
	result := 0.0
	for _, contribution := range s {
		result = result + contribution.Value
	}
	return result
}
