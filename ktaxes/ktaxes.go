package ktaxes

type Allowance struct {
	AllowanceType string  `json:"allowanceType" example:"k-receipt"`
	Amount        float64 `json:"amount" example:"100.00"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome" example:"500000.00"`
	Wht         float64     `json:"wht" example:"0.0"`
	Allowances  []Allowance `json:"allowances" example:"{donation,0.0}"`
}

type TaxResponse struct {
	Tax       float64    `json:"tax" example:"29000"`
	TaxLevels []TaxLevel `json:"taxLevel" example:"{0-150,000, 0.0}"`
	TaxRefund float64    `json:"taxRefund,omitempty" example:"0.0"`
}

type TaxLevel struct {
	Level string  `json:"level" example:"0-150,000"`
	Tax   float64 `json:"tax" example:"0.0"`
}
