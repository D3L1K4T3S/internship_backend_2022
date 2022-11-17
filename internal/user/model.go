package user

type User struct {
	Id       string `json:"id,omitempty"`
	Balance  string `json:"balance,omitempty"`
	Reserved string `json:"reserved,omitempty"`
}

type CreateUserDTO struct {
	Id       string `json:"id,omitempty"`
	Balance  string `json:"balance,omitempty"`
	Reserved string `json:"reserved,omitempty"`
}
