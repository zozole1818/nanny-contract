package model

var MinimumWageInPoland int64 = 4666
var BaseForCalculationOfContributionsPaidByCountry = MinimumWageInPoland / 2

var PensionContribution = 19.52          // % składki emerytalnej
var SocialPensionContribution = 8.0      // % składki rentowej
var AccidentInsuranceContribution = 1.67 // % składki wypadkowej
var HealthInsuranceContribution = 9.0    // % składki zdrowotnej
var SicknessContribution = 2.45          // % składki zdrowotnej

// paid by family
var (
	PensionContributionPaidByEmployer           = PensionContribution / 2 // %
	SocialPensionContributionPaidByEmployer     = 6.5                     // %
	AccidentInsuranceContributionPaidByEmployer = 1.67                    // %
	HealthInsuranceContributionPaidByEmployer   = 0.0                     // %
)

// paid by nanny
var (
	PensionContributionPaidByEmployee           = PensionContribution / 2 // %
	SocialPensionContributionPaidByEmployee     = 1.5                     // %
	AccidentInsuranceContributionPaidByEmployee = 0.0                     // %
	HealthInsuranceContributionPaidByEmployee   = 9.0                     // %
)

//var (
//	pensionContribution           = "składka emerytalna"
//	socialPensionContribution     = "składka rentowa"
//	accidentInsuranceContribution = "składka wypadkowa"
//	healthInsuranceContribution   = "składka zdrowotna"
//	sicknessContribution          = "składka chorobowa"
//)

type ContributionName int

const (
	//paid by country
	countryPensionContribution ContributionName = iota
	countrySocialPensionContribution
	countryAccidentInsuranceContribution
	countryHealthInsuranceContribution
	countrySicknessContribution

	//paid by employer
	employerPensionContribution
	employerSocialPensionContribution
	employerAccidentInsuranceContribution
	employerHealthInsuranceContribution
	employerSicknessContribution

	//paid by employee
	employeePensionContribution
	employeeSocialPensionContribution
	employeeAccidentInsuranceContribution
	employeeHealthInsuranceContribution
	employeeSicknessContribution
)

func (n ContributionName) name() string {
	return [...]string{
		"składka emerytalna",
		"składka rentowa",
		"składka wypadkowa",
		"składka zdrowotna",
		"składka chorobowa",

		"składka emerytalna",
		"składka rentowa",
		"składka wypadkowa",
		"składka zdrowotna",
		"składka chorobowa",

		"składka emerytalna",
		"składka rentowa",
		"składka wypadkowa",
		"składka zdrowotna",
		"składka chorobowa",
	}[n]
}

//var countryContributionList = []string{
//	pensionContribution,
//	socialPensionContribution,
//	accidentInsuranceContribution,
//	healthInsuranceContribution,
//}

var contributionsMap = map[ContributionName]float64{
	countryPensionContribution:           19.52,
	countrySocialPensionContribution:     8.0,
	countryAccidentInsuranceContribution: 1.67,
	countryHealthInsuranceContribution:   9.0,
	countrySicknessContribution:          0.0,

	employerPensionContribution:           9.76,
	employerSocialPensionContribution:     6.5,
	employerAccidentInsuranceContribution: 1.67,
	employerHealthInsuranceContribution:   0.0,
	employerSicknessContribution:          0.0,

	employeePensionContribution:           9.76,
	employeeSocialPensionContribution:     1.5,
	employeeAccidentInsuranceContribution: 0.0,
	employeeHealthInsuranceContribution:   9.0,
	employeeSicknessContribution:          2.45,
}

//var contributionPaidByCountryMap = map[string]float64{
//	pensionContribution:           PensionContribution,
//	socialPensionContribution:     SocialPensionContribution,
//	accidentInsuranceContribution: AccidentInsuranceContribution,
//	healthInsuranceContribution:   HealthInsuranceContribution,
//}
//
//var contributionPaidByEmployerMap = map[string]float64{
//	pensionContribution:           PensionContributionPaidByEmployer,
//	socialPensionContribution:     SocialPensionContributionPaidByEmployer,
//	accidentInsuranceContribution: AccidentInsuranceContributionPaidByEmployer,
//	healthInsuranceContribution:   HealthInsuranceContributionPaidByEmployer,
//}
//
//var contributionPaidByEmployeeMap = map[string]float64{
//	pensionContribution:           PensionContributionPaidByEmployee,
//	socialPensionContribution:     SocialPensionContributionPaidByEmployee,
//	accidentInsuranceContribution: AccidentInsuranceContributionPaidByEmployee,
//	healthInsuranceContribution:   HealthInsuranceContributionPaidByEmployee,
//	sicknessContribution:          SicknessContribution,
//}
