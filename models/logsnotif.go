package models

type Notif struct {
	ID           int
	Receiver     int
	NotifType    int
	NotifMessage string
	MessageLink  string
	Read         bool
}

type Log struct {
	ID          int
	LogType     string
	LogMessage  string
	Date        string
	MessageLink string
	Requester   int
	Handled     bool
}

type Notifications struct {
	Notifs  []Notif
	NotRead int
}
