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
	UpdateAllowance(amount float64, allowance_type string) (float64, error)
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

func (h *Handler) PersonalDeductionHandler(c echo.Context) error {
	body := struct {
		Amount float64 `json:"amount"`
	}{}

	err := c.Bind(&body)
	if body.Amount <= 10000 || body.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "personal amount must between 10,001 - 100,000"})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	updatedAllowance, err := h.store.UpdateAllowance(body.Amount, "personal")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, struct {
		PersonalDeduction float64 `json:"personalDeduction"`
	}{updatedAllowance})
}

func (h *Handler) TaxCalculationsHandler(c echo.Context) error {
	taxRequest := TaxRequest{}
	if err := c.Bind(&taxRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	allowances, err := h.store.Allowances()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	personalDeduction := 0.0
	for _, allowance := range allowances {
		if allowance.AllowanceType == "personal" {
			personalDeduction = allowance.Amount
		}
	}

	tax, taxLevels := taxCalculate(taxRequest.TotalIncome, taxRequest.Wht, taxRequest.Allowances, personalDeduction)
	res := TaxResponse{}
	res.TaxLevels = taxLevels
	res.Tax = tax
	if tax < 0 {
		res.Tax = 0
		res.TaxRefund = 0 - tax
	}
	return c.JSON(http.StatusOK, res)
}

func taxCalculate(totalIncome float64, wht float64, allowances []Allowance, personalDeduction float64) (float64, []TaxLevel) {
	netIncome := totalIncome - personalDeduction

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
