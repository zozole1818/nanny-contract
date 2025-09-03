package view

import (
	"github.com/labstack/echo/v4"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"github.com/zozole1818/nanny-contract/internal/reports/zus"
	"html/template"
	"io"
)

type Template struct {
	Templates *template.Template
}

func (t Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

var FuncMap = template.FuncMap{
	"add": func(a, b float64) float64 {
		return a + b
	},
}

var Title = "Umowa Aktywacyjna dla Niani"

type PageDetails struct {
	Title          string
	ReportRequest  model.ReportRequest
	ReportResponse model.ReportResponseView
	GeneratePDF    bool
}

type PageDetails2 struct {
	Title       string
	GrossIncome string
	Summary     zus.Summary
	GeneratePDF bool
	SendEmail   bool
}
