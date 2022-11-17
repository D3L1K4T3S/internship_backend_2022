package orders

type Orders struct {
	Id            string `json:"orderId,omitempty"`
	UserId        string `json:"userId,omitempty"`
	ServiceId     string `json:"serviceId,omitempty"`
	TransactionIs string `json:"transactionIs,omitempty"`
	Status        string `json:"status,omitempty"`
	Cost          string `json:"cost,omitempty"`
	Date          string `json:"date,omitempty"`
}
