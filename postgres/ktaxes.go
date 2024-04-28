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

func (p *Postgres) UpdateAllowance(amount float64, allowance_type string) (float64, error) {
	row := p.Db.QueryRow("INSERT INTO allowances (allowance_type, amount) VALUES ($1, $2) ON CONFLICT (allowance_type) DO UPDATE SET amount = EXCLUDED.amount RETURNING amount;", allowance_type, amount)
	err := row.Scan(&amount)
	if err != nil {
		return 0, err
	}

	return amount, nil
}
