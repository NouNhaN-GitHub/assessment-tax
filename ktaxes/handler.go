package ktaxes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Allowances() ([]Allowance, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

func (h *Handler) AllowanceHandler(c echo.Context) error {
	allowances, err := h.store.Allowances()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, allowances)
}

func (h *Handler) TaxCalculationsHandler(c echo.Context) error {
	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	tax, taxLevels := taxCalculate(taxRequest.TotalIncome, taxRequest.Wht, taxRequest.Allowances)
	res := TaxResponse{}
	res.TaxLevels = taxLevels
	res.Tax = tax
	if tax < 0 {
		res.Tax = 0
		res.TaxRefund = 0 - tax
	}
	return c.JSON(http.StatusOK, res)
}

func taxCalculate(totalIncome float64, wht float64, allowances []Allowance) (float64, []TaxLevel) {
	netIncome := totalIncome - 60000

	amountDonation := 0.0
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			amountDonation = allowance.Amount
			if allowance.Amount > 100000 {
				amountDonation = 100000
			}

		}
	}
	netIncome = netIncome - amountDonation

	displayTaxLevels := []TaxLevel{
		{"0-150,000", 0.0},
		{"150,001-500,000", 0.0},
		{"500,001-1,000,000", 0.0},
		{"1,000,001-2,000,000", 0.0},
		{"2,000,001 ขึ้นไป", 0.0},
	}

	taxLevels := []struct {
		incomeDiff float64
		incomeTax  float64
	}{
		{150000, 0},
		{350000, 0.1},
		{500000, 0.15},
		{1000000, 0.2},
		{-1, 0.35},
	}

	if netIncome <= 150000 {
		return 0, displayTaxLevels
	}

	tax := 0.0
	for idx, taxLevel := range taxLevels {
		if netIncome > taxLevel.incomeDiff && taxLevel.incomeDiff != -1 {
			tax += taxLevel.incomeDiff * taxLevel.incomeTax
			netIncome -= taxLevel.incomeDiff
			displayTaxLevels[idx].Tax = taxLevel.incomeDiff * taxLevel.incomeTax
			continue
		}
		tax += netIncome * taxLevel.incomeTax
		displayTaxLevels[idx].Tax = netIncome * taxLevel.incomeTax
		netIncome = 0
	}

	tax = tax - wht

	return tax, displayTaxLevels
}
