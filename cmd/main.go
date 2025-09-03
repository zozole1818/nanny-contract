package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/zozole1818/nanny-contract/internal/email"
	"github.com/zozole1818/nanny-contract/internal/reports"
	"github.com/zozole1818/nanny-contract/internal/reports/model"
	"github.com/zozole1818/nanny-contract/internal/reports/pdf"
	"github.com/zozole1818/nanny-contract/internal/reports/view"
	"github.com/zozole1818/nanny-contract/internal/reports/zus"
	"gopkg.in/yaml.v3"
	"html/template"
	"log/slog"

	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := loadEnvs(".env")
	if err != nil {
		slog.Error("Error when loading env file", "error", err)
		return
	}

	transformer := model.JSON{}
	svc := reports.NewReportSvc(transformer)
	pdfGen := pdf.NewGenerator()
	emailer := email.NewEmailer("smtp.gmail.com", 587, "zuzanna.go.mail@gmail.com")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Warn("Interrupt signal detected. Closing down...")
		cancel()
	}()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("public/static"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	t := &view.Template{
		Templates: template.Must(template.New("").Funcs(view.FuncMap).ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/views/reports")
	})

	e.GET("/views/reports", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", view.PageDetails{Title: view.Title, GeneratePDF: false})
	})

	e.GET("/views/reports-new", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index2.html", view.PageDetails2{Title: view.Title, GeneratePDF: false, SendEmail: false})
	})

	e.POST("/api/v1/reports", func(c echo.Context) error {
		request := model.ReportRequest{}
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		resp, err := svc.CreateReport(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		return c.JSON(http.StatusCreated, resp)
	})

	e.POST("/views/reports", func(c echo.Context) error {
		grossIncome, _ := strconv.ParseFloat(c.FormValue("grossIncome"), 64)
		sicCont := c.FormValue("sicknessContribution") == "on"
		action := c.FormValue("action")
		request := model.ReportRequest{GrossIncome: grossIncome, SicknessContribution: sicCont}
		resp, err := svc.CreateReport(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		pd := view.PageDetails{
			Title:          view.Title,
			ReportRequest:  request,
			ReportResponse: resp.ToReportResponseView(),
		}
		if action == "download-pdf" {
			pd.GeneratePDF = true
		}
		return c.Render(http.StatusOK, "index.html", pd)
	})

	e.POST("/views/reports-new", func(c echo.Context) error {
		grossIncome, err := decimal.NewFromString(c.FormValue("grossIncome"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		sicCont := c.FormValue("sicknessContribution") == "on"
		action := c.FormValue("action")
		//request := model.ReportRequest{GrossIncome: grossIncome.F, SicknessContribution: sicCont}
		dra, err := zus.Generator{}.NannyDRA(time.August, decimal.NewFromInt(model.MinimumWageInPoland), grossIncome, sicCont)
		//resp, err := svc.CreateReport(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		pd := view.PageDetails2{
			Title:       view.Title,
			GrossIncome: grossIncome.String(),
			Summary:     dra.CreateSummary(),
			GeneratePDF: false,
			SendEmail:   false,
		}
		switch action {
		case "download-pdf":
			pd.GeneratePDF = true
		case "send-email":
			pd.SendEmail = true
		}
		return c.Render(http.StatusOK, "index2.html", pd)
	})

	e.GET("/download-pdf", func(c echo.Context) error {
		grossIncome, _ := strconv.ParseFloat(c.QueryParam("grossIncome"), 64)
		sicCont := c.QueryParam("sicknessContribution") == "true"
		request := model.ReportRequest{GrossIncome: grossIncome, SicknessContribution: sicCont}
		resp, err := svc.CreateReport(request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		b, err := pdfGen.Generate(resp)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		c.Response().Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
		return c.Blob(http.StatusOK, "application/pdf", b)
	})

	e.GET("/send-email", func(c echo.Context) error {
		grossIncome, err := decimal.NewFromString(c.QueryParam("grossIncome"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		dra, err := zus.Generator{}.NannyDRA(time.August, decimal.NewFromInt(model.MinimumWageInPoland), grossIncome, true)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		err = emailer.Send([]string{"zuzua.tel@gmail.com"}, dra.CreateSummary())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err))
		}
		return c.JSON(http.StatusNoContent, nil)
	})

	go func() {
		// Start server
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	<-ctx.Done()

	// Shutdown
	shutdownContext, shutdownContextCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownContextCancel()
	if err := e.Shutdown(shutdownContext); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}
	<-shutdownContext.Done()

	slog.Info("Server shutdown complete.")
}

func loadEnvs(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error when reading %s file: %v", path, err)
	}
	result := make(map[string]string)
	err = yaml.Unmarshal(b, &result)
	if err != nil {
		return fmt.Errorf("error when Unmarshal %s file: %v", path, err)
	}
	var errors []string
	for k, v := range result {
		err = os.Setenv(k, v)
		if err != nil {
			errors = append(errors, fmt.Errorf("error when setting %s env var: %v", k, err).Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("errors when setting env vars: %s", strings.Join(errors, ", "))
	}
	return nil
}
