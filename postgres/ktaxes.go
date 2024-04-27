package postgres

import (
	"github.com/NouNhaN-GitHub/assessment-tax/ktaxes"
)

type Allowance struct {
	AllowanceType string  `postgres:"allowance_type"`
	Amount        float64 `postgres:"amount"`
}

func (p *Postgres) Allowances() ([]ktaxes.Allowance, error) {
	rows, err := p.Db.Query("SELECT * FROM allowances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allowances []ktaxes.Allowance
	for rows.Next() {
		var a Allowance
		err := rows.Scan(&a.AllowanceType,
			&a.Amount,
		)
		if err != nil {
			return nil, err
		}
		allowances = append(allowances, ktaxes.Allowance{
			AllowanceType: a.AllowanceType,
			Amount:        a.Amount,
		})
	}
	return allowances, nil
}
