package pdf

import (
	"github.com/signintech/gopdf"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
)

type RCACreator struct {
}

func NewRCACreator() Generator {
	return &RCACreator{}
}

func (R RCACreator) Generate(response model.ReportResponse) ([]byte, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	// Import the ZUS blank form as background
	//tpl := pdf.ImportPage("C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\internal\\reports\\pdf\\ZUS_RCA.pdf", 1, "/MediaBox")
	//pdf.UseImportedTemplate(tpl, 0, 0, 210, 297)

	err := pdf.Image("C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\internal\\reports\\pdf\\ZUS_RCA.pdf", 0, 0, &gopdf.Rect{W: 210, H: 297})
	if err != nil {
		return nil, err
	}

	err = pdf.AddTTFFont("arial", "C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\cmd\\arial.ttf") // you need the .ttf file
	if err != nil {
		return nil, err
	}
	pdf.SetXY(45.45, 142.72)
	err = pdf.Cell(nil, "78442536156")
	if err != nil {
		return nil, err
	}

	err = pdf.WritePdf("test_rca.pdf")
	if err != nil {
		return nil, err
	}
	return nil, err
}
