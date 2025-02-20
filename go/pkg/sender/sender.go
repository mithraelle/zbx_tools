package sender

type Item struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int    `json:"clock"`
	Ns    int    `json:"ns"`
}

type ItemSendError struct {
	Items []Item
	Err   error
}

type Sender interface {
	Send([]Item, chan<- ItemSendError)
}
