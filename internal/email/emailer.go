package email

import (
	"bytes"
	"fmt"
	"github.com/zozole1818/nanny-contract/internal/reports/view"
	"github.com/zozole1818/nanny-contract/internal/reports/zus"
	"gopkg.in/gomail.v2"
	"html/template"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type Emailer struct {
	host      string
	port      int
	from      string
	whitelist []string
}

func NewEmailer(host string, port int, from string) *Emailer {
	whitelist := os.Getenv("EMAIL_WHITELIST")
	return &Emailer{
		host:      host,
		port:      port,
		from:      from,
		whitelist: strings.Split(whitelist, ","),
	}
}

func (e *Emailer) Send(to []string, summary zus.Summary) error {
	var approvedTo []string
	for _, addr := range to {
		if slices.Contains(e.whitelist, addr) {
			approvedTo = append(approvedTo, addr)
		} else {
			slog.Warn("Email not approved", "email", addr)
		}
	}
	pwd, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !ok {
		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS env var not set")
	}

	body, err := e.prepareEmailBody(summary)
	if err != nil {
		return fmt.Errorf("error when preparing email body: %v", err)
	}
	fromAddr := e.from
	m := gomail.NewMessage()

	m.SetHeader("From", fromAddr)
	m.SetHeader("To", approvedTo...)
	m.SetHeader("Subject", fmt.Sprintf("[Nanny Contract Calculation] %s", summary.Name))
	m.SetBody("text/html", body)

	// Use app password here
	d := gomail.NewDialer(e.host, e.port, fromAddr, pwd)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error when sending email: %v", err)
	}
	return nil
}

func (e *Emailer) prepareEmailBody(summary zus.Summary) (string, error) {
	pd := view.PageDetails2{
		Title:       "Nanny Contract Calculation",
		Summary:     summary,
		GeneratePDF: false,
		SendEmail:   false,
	}
	htmlTpl := `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
</head>
<body style="margin:0; padding:0; background-color:#f0f8ff; font-family:Arial,sans-serif;">

  <table role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%">
    <tr>
      <td align="center" style="padding:20px 0;">
        <table width="600" cellpadding="20" cellspacing="0" border="0" style="background:#ffffff; border-radius:8px; box-shadow:0 2px 6px rgba(0,0,0,0.1);">
          <tr>
            <td>
              <h2 style="font-size:24px; margin:0;">{{.Title}}</h2>
              <p style="color:#888; font-size:12px; margin:0;">by Zuzua</p>

              <hr style="border:0; border-top:1px solid #ccc; margin:20px 0;">

              {{ if .Summary.Name }}
              <h3 style="font-size:20px; margin:0 0 10px;">{{ .Summary.Name }}</h3>

              <table width="100%" cellpadding="5" cellspacing="0" border="0">
                <tr>
                  <td>Kwota na umowie (brutto):</td>
                  <td align="right"><strong>{{.Summary.GrossIncome}} zł</strong></td>
                </tr>
                <tr>
                  <td>Na rękę dla niani:</td>
                  <td align="right"><strong>{{.Summary.NetIncome}} zł</strong></td>
                </tr>
                <tr>
                  <td>Całkowity koszt rodzica:</td>
                  <td align="right"><strong>{{.Summary.EmployerTotalCost}} zł</strong></td>
                </tr>
                <tr>
                  <td>Suma wszystkich składek ZUS:</td>
                  <td align="right"><strong>{{.Summary.Zus.EmployerPaid}} zł</strong></td>
                </tr>
              </table>

              <hr style="border:0; border-top:1px solid #ccc; margin:20px 0;">

              <h3 style="font-size:18px;">Składki opłacane przez rodzica:</h3>
              <table width="100%" cellpadding="5" cellspacing="0" border="0">
                {{ range $key, $value := .Summary.PaidByEmployer }}
                <tr>
                  <td>{{$value.Name}}</td>
                  <td align="right">{{$value.Value}} zł</td>
                </tr>
                  {{ range $key2, $value2 := $value.Records }}
                  <tr>
                    <td style="padding-left:20px; font-size:13px; color:#555;">{{$value2.Name}}</td>
                    <td align="right">{{$value2.Value}} zł</td>
                  </tr>
                  {{ end }}
                {{ end }}
              </table>

              <h3 style="font-size:18px; margin-top:20px;">Składki opłacane przez państwo:</h3>
              <table width="100%" cellpadding="5" cellspacing="0" border="0">
                {{ range $key, $value := .Summary.PaidByCountry }}
                <tr>
                  <td>{{$key}}</td>
                  <td align="right">{{$value.Value}} zł</td>
                </tr>
                {{ end }}
              </table>
              {{ end }}
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>

</body>
</html>

`
	var buf bytes.Buffer
	t := template.Must(template.New("emailHTML").Parse(htmlTpl))
	if err := t.Execute(&buf, pd); err != nil {
		return "", err
	}
	return buf.String(), nil
}
