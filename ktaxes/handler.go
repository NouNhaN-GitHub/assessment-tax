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

// WalletHandler
//
//	@Summary		Get all allowances
//	@Description	Get all allowances
//	@Tags			allowance
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}		allowance
//	@Router			/api/v1/allowances 	[get]
//	@Failure		500	{object}		Err
func (h *Handler) AllowanceHandler(c echo.Context) error {
	allowances, err := h.store.Allowances()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, allowances)
}
