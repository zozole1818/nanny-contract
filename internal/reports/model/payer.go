package model

import "github.com/zozole1818/nanny-contract/internal/misc"

type PayerType int

const (
	Country  PayerType = 0
	Employer PayerType = 1
	Employee PayerType = 2
)

type Payer struct {
	Type                          PayerType
	Base                          float64
	PensionContribution           float64
	SocialPensionContribution     float64
	AccidentInsuranceContribution float64
	HealthInsuranceContribution   float64
	SicknessContribution          float64
}

func NewPayer(pType PayerType, base float64, pensionContribution float64, socialPensionContribution float64, accidentInsuranceContribution float64, healthInsuranceContribution float64, sicknessContribution float64) Payer {
	return Payer{
		Type:                          pType,
		Base:                          base,
		PensionContribution:           pensionContribution,
		SocialPensionContribution:     socialPensionContribution,
		AccidentInsuranceContribution: accidentInsuranceContribution,
		HealthInsuranceContribution:   healthInsuranceContribution,
		SicknessContribution:          sicknessContribution,
	}
}

func (p Payer) calulatePensionContribution() Contribution {
	return Contribution{
		Name:       "składka emerytalna",
		Value:      misc.CalculatePercentage(p.Base, p.PensionContribution),
		Base:       p.Base,
		Percentage: p.PensionContribution,
	}
}

func (p Payer) calulateSocialPensionContribution() Contribution {
	return Contribution{
		Name:       "składka rentowa",
		Value:      misc.CalculatePercentage(p.Base, p.SocialPensionContribution),
		Base:       misc.Round(p.Base),
		Percentage: p.SocialPensionContribution,
	}
}

func (p Payer) calulateAccidentInsuranceContribution() Contribution {
	return Contribution{
		Name:       "składka wypadkowa",
		Value:      misc.CalculatePercentage(p.Base, p.AccidentInsuranceContribution),
		Base:       misc.Round(p.Base),
		Percentage: p.AccidentInsuranceContribution,
	}
}

func (p Payer) calulateHealthInsuranceContribution() Contribution {
	sicknessContr := p.calulateSicknessContribution()
	base := p.Base
	switch p.Type {
	case Country:
		base = p.Base - sicknessContr.Value
	case Employer: // contribition is 0% so it doesn't really matter
		base = p.Base
	case Employee:
		pension := p.calulatePensionContribution()
		social := p.calulateSocialPensionContribution()
		base = p.Base - pension.Value - social.Value - sicknessContr.Value
	}

	return Contribution{
		Name:       "składka zdrowotna",
		Value:      misc.CalculatePercentage(base, p.HealthInsuranceContribution),
		Base:       misc.Round(base),
		Percentage: p.HealthInsuranceContribution,
	}
}

func (p Payer) calulateSicknessContribution() Contribution {
	return Contribution{
		Name:       "składka chorobowa",
		Value:      misc.CalculatePercentage(p.Base, p.SicknessContribution),
		Base:       misc.Round(p.Base),
		Percentage: p.SicknessContribution,
	}
}

func (p Payer) CalulateSicknessContribution(baseForSicknesscontribution float64) Contribution {
	return Contribution{
		Name:       "składka chorobowa",
		Value:      misc.CalculatePercentage(baseForSicknesscontribution, p.SicknessContribution),
		Base:       misc.Round(baseForSicknesscontribution),
		Percentage: p.SicknessContribution,
	}
}

func (p Payer) CalculateContributions(baseForSicknessContribution float64) []Contribution {
	result := make([]Contribution, 0)
	result = append(result, p.calulatePensionContribution())
	result = append(result, p.calulateSocialPensionContribution())
	result = append(result, p.calulateAccidentInsuranceContribution())
	result = append(result, p.calulateHealthInsuranceContribution())
	// we calculate sickness contribution here because we use grossIncome as a base
	sic := p.CalulateSicknessContribution(baseForSicknessContribution)
	if p.Type != Country && sic.Value > 0 {
		result = append(result, sic)
	}

	return result
}
