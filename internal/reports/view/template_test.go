package view

import (
	"bytes"
	"fmt"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"html/template"
	"testing"
)

func TestRender(t *testing.T) {
	// Arrange
	tmpl := &Template{
		Templates: template.Must(template.New("").Funcs(FuncMap).ParseGlob("*.html")),
	}
	data := PageDetails{
		Title:          Title,
		ReportRequest:  model.ReportRequest{},
		ReportResponse: model.ReportResponse{Name: ""}.ToReportResponseView(),
	}
	// Act
	buf := new(bytes.Buffer)
	err := tmpl.Render(buf, "index.html", data, nil)
	if err != nil {
		t.Errorf("Render returned an error: %v", err)
	}

	// Assert
	actualOutput := buf.String()
	fmt.Print(actualOutput)
}
