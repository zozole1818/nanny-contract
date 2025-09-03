package pdf

import (
	"github.com/zozole1818/nanny-contract/internal/reports"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"os"
	"testing"
)

func TestRCACreator_Generate(t *testing.T) {
	svc := reports.NewReportSvc(model.JSON{})
	rq := model.ReportRequest{
		GrossIncome:          3700,
		SicknessContribution: true,
	}
	resp, err := svc.CreateReport(rq)
	if err != nil {
		t.Fatal(err)
	}

	pdfGenerator := NewRCACreator()
	b, err := pdfGenerator.Generate(resp)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("C:\\cygwin64\\home\\zuzanna.tomaszewska\\return-to-work\\nanny-contract\\backend\\cmd\\test_rca.pdf", b, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
