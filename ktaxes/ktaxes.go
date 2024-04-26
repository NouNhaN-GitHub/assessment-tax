package ktaxes

import "time"

type Allowance struct {
	AllowanceType string    `json:"allowance_type" example:"k-receipt"`
	Amount        float64   `json:"amount" example:"100.00"`
	CreatedAt     time.Time `json:"created_at" example:"2024-04-24T14:19:00.729237Z"`
}
