package sender

import "fmt"

type emailMessage struct {
	title        string
	bodyTemplate string
}

func (m emailMessage) Title() string {
	return m.title
}

func (m emailMessage) Body(code string) string {
	return fmt.Sprintf(m.bodyTemplate, code)
}

func NewEmailMessage(title, template string) CodeMessenger {
	return emailMessage{
		title:        title,
		bodyTemplate: template,
	}
}

func DefaultMailMessage() CodeMessenger {
	return NewEmailMessage("Verify Code", "You code is %s")
}
