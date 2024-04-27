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
	tax := taxCalculate(taxRequest.TotalIncome)
	res := TaxResponse{tax}
	return c.JSON(http.StatusOK, res)
}

func taxCalculate(totalIncome float64) float64 {
	netIncome := totalIncome - 60000

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
		return 0
	}

	tax := 0.0
	for _, taxLevel := range taxLevels {
		if netIncome > taxLevel.incomeDiff && taxLevel.incomeDiff != -1 {
			tax += taxLevel.incomeDiff * taxLevel.incomeTax
			netIncome -= taxLevel.incomeDiff
			continue
		}
		tax += netIncome * taxLevel.incomeTax
		netIncome = 0
	}
	return tax
}
