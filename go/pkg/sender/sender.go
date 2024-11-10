package sender

type Item struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int    `json:"clock"`
	Ns    int    `json:"ns"`
}

type Sender interface {
	Send(items []Item, try int) error
}
