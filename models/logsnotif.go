package models

type Notif struct {
	ID           int
	Receiver     int
	NotifType    int
	NotifMessage string
	Read         bool
}

type Log struct {
	ID          int
	LogType     string
	LogMessage  string
	Date        string
	Requester   int
	Handled     bool
	MessageLink string
}

type Notifications struct {
	Notifs  []Notif
	NotRead int
}
