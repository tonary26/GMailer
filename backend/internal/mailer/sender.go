package mailer

type Message struct {
	To 		string
	Subject	string
	Body 	string
}

type Sender interface {
	Send(msg *Message) error
}