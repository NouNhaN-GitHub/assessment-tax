package ktaxes

type Allowance struct {
	AllowanceType string  `json:"allowance_type" example:"k-receipt"`
	Amount        float64 `json:"amount" example:"100.00"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome" example:"500000.00"`
	Wht         float64     `json:"wht" example:"0.0"`
	Allowances  []Allowance `json:"allowances" example:"{donation,0.0}"`
}

type TaxResponse struct {
	Tax       float64 `json:"tax" example:"29000"`
	TaxRefund float64 `json:"taxRefund,omitempty" example:"0.0"`
}
