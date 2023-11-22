package sender

import (
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	From     string `json:"from"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type emailCli struct {
	config     *EmailConfig
	messagener CodeMessenger
}

// NewEmailSender 实例化邮件发送客户端
func NewEmailSender(cfg *EmailConfig, message CodeMessenger) SmsCodeSender {
	return &emailCli{config: cfg, messagener: message}
}

// Send 邮件发送
func (c *emailCli) Send(code, identify string) error {
	d := gomail.NewDialer(c.config.Host, c.config.Port, c.config.Username, c.config.Password)
	m := gomail.NewMessage()
	m.SetHeader("From", c.config.From)
	m.SetHeader("To", identify)
	m.SetHeader("Subject", c.messagener.Title())
	m.SetBody("text/html", c.messagener.Body(code))
	sendQueue <- func() {
		err := d.DialAndSend(m)
		if err != nil {
			// 输出日志
		}
	}
	return nil
}
