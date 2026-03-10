package events

type HelloEvent struct {
	ID   string
	Name string
	Type string
}

func (e HelloEvent) EventName() string {
	return "notifications.hello"
}

func (e HelloEvent) EventKey() string {
	return e.ID
}
