package pdf

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
)

type Generator interface {
	Generate(response model.ReportResponse) ([]byte, error)
}

func NewGenerator() Generator {
	return &generator{}
}

type generator struct {
}

func (g *generator) Generate(response model.ReportResponse) ([]byte, error) {

	view := response.ToReportResponseView()

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("arial", "C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\cmd\\arial.ttf") // you need the .ttf file
	if err != nil {
		return nil, err
	}

	//title
	err = printCell(pdf, 120, 25, 25, view.Name)
	if err != nil {
		return nil, err
	}

	yAnchor := 80.0

	err = printCell(pdf, 50, yAnchor, 21, "Podsumowanie")
	if err != nil {
		return nil, err
	}

	tbl := table{
		pdf:      pdf,
		x:        50,
		y:        yAnchor + 30,
		fontSize: 12,
		cellW:    230,
		cellH:    20,
		cells: [][]string{
			{"Kwota na umowie (brutto):", fmt.Sprintf("%.2f zł", view.GrossIncome)},
			{"Na rękę dla niani:", fmt.Sprintf("%.2f zł", view.NetIncome)},
			{"Całkowity koszt rodzica:", fmt.Sprintf("%.2f zł", view.EmployerTotalCost)},
			{"Suma wszystkich składek ZUS:", fmt.Sprintf("%.2f zł", view.ZUSAmount)},
		},
	}

	err = tbl.printTable()
	if err != nil {
		return nil, err
	}

	err = line(pdf, tbl.getTableEnd()+10)
	if err != nil {
		return nil, err
	}

	detailsAnchor := tbl.getTableEnd() + 20
	err = printCell(pdf, 50, detailsAnchor, 15, "Szczegóły")
	if err != nil {
		return nil, err
	}

	err = printCell(pdf, 50, detailsAnchor+30, 13, "Składki opłacane przez rodzica:")
	if err != nil {
		return nil, err
	}

	tbl = table{
		pdf:      pdf,
		x:        50,
		y:        detailsAnchor + 30 + 30,
		fontSize: 12,
		cellW:    230,
		cellH:    20,
	}

	for _, item := range view.PaidByEmployerCompact {
		err = tbl.printRow(0, 12, item.Name, fmt.Sprintf("%.2f zł", item.Value))
		if err != nil {
			return nil, err
		}
		err = tbl.printRow(10, 10, "opłacana przez rodzica", fmt.Sprintf("%.2f zł", item.PaidByEmployer))
		if err != nil {
			return nil, err
		}
		err = tbl.printRow(10, 10, "opłacana przez nianie", fmt.Sprintf("%.2f zł", item.PaidByEmployee))
		if err != nil {
			return nil, err
		}

	}

	anchor := tbl.getTableEnd() + 15

	err = printCell(pdf, 50, anchor, 13, "Składki opłacane przez państwo:")
	if err != nil {
		return nil, err
	}

	tbl = table{
		pdf:      pdf,
		x:        50,
		y:        anchor + 30,
		fontSize: 12,
		cellW:    230,
		cellH:    20,
		cells:    toCells(view.ContributionPaidByCountry),
	}

	err = tbl.printTable()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func printCell(pdf *gopdf.GoPdf, x, y float64, size int, text string) error {
	err := pdf.SetFont("arial", "", size)
	if err != nil {
		return err
	}
	pdf.SetXY(x, y)
	return pdf.Cell(nil, text)
}

type table struct {
	pdf      *gopdf.GoPdf
	x, y     float64
	fontSize int
	cellW    int
	cellH    int
	rowCount int
	cells    [][]string
	printed  bool
}

func (t *table) getTableEnd() float64 {
	return float64(t.rowCount*t.cellH) + t.y
}

func (t *table) printRow(intend float64, size int, c1, c2 string) error {
	defer func() { t.rowCount++ }()
	y := t.y + float64(t.rowCount*t.cellH)
	err := printCell(t.pdf, intend+t.x, y, size, c1)
	if err != nil {
		return err
	}
	return printCell(t.pdf, intend+t.x+float64(t.cellW), y, size, c2)
}

func (t *table) printTable() error {
	defer func() { t.printed = true }()
	for _, row := range t.cells {
		err := t.printRow(0, t.fontSize, row[0], row[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func toCells(arr []model.Contribution) [][]string {
	result := make([][]string, 0)
	for _, item := range arr {
		result = append(result, []string{
			item.Name,
			fmt.Sprintf("%.2f zł", item.Value),
		})
	}
	return result
}

func line(pdf *gopdf.GoPdf, y float64) error {
	err := pdf.SetTransparency(gopdf.Transparency{Alpha: 0.5, BlendModeType: gopdf.ColorBurn})
	if err != nil {
		return err
	}
	//pdf.SetLineType("dotted")
	pdf.SetStrokeColor(211, 211, 211)
	pdf.SetLineWidth(1.5)
	pdf.Line(30, y, 570, y)
	pdf.ClearTransparency()
	return nil
}
