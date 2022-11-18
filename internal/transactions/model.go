package transactions

import "github.com/jackc/pgx/v5/pgtype"

type Transactions struct {
	Name      string           `json:"name,omitempty"`
	Type      string           `json:"type,omitempty"`
	TimeTrans pgtype.Timestamp `json:"timeTrans,omitempty"`
	Amount    string           `json:"amount,omitempty"`
	Pass      string           `json:"pass,omitempty"`
	Comment   *string          `json:"comment,omitempty"`
}
