package transactions

type Transactions struct {
	Id        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	TimeTrans string `json:"timeTrans,omitempty"`
	Amount    string `json:"amount,omitempty"`
	Pass      string `json:"pass,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

type CreateUserDTO struct {
	Id       string `json:"id,omitempty"`
	Balance  string `json:"balance,omitempty"`
	Reserved string `json:"reserved,omitempty"`
}
