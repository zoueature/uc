package sender

// SmsCodeSender 验证码发送器
type SmsCodeSender interface {
	Send(code, identify string) error
}

// CodeMessenger 验证码消息模版
type CodeMessenger interface {
	Body(code string) string
	Title() string
}
