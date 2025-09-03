package zus

var contributionNamePolMap = map[ContributionName]string{
	Pension:  "składka emerytalna",
	Rental:   "składka rentowa",
	Sick:     "składka chorobowa",
	Accident: "składka wypadkowa",
	Health:   "składka zdrowotna",
}

type paidBy string

var (
	employeePaidBy paidBy = "employee"
	employerPaidBy paidBy = "employer"
	countryPaidBy  paidBy = "country"
)

var paidByPolMap = map[paidBy]string{
	employeePaidBy: "opłacana przez nianie",
	employerPaidBy: "opłacana przez rodzica",
	countryPaidBy:  "opłacana przez państwo",
}

var monthsPolMap = map[string]string{
	"January":   "Styczeń",
	"February":  "Luty",
	"March":     "Marzec",
	"April":     "Kwiecień",
	"May":       "Maj",
	"June":      "Czerwiec",
	"July":      "Lipiec",
	"August":    "Sierpień",
	"September": "Wrzesień",
	"October":   "Październik",
	"November":  "Listopad",
	"December":  "Grudzień",
}
