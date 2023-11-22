package sender

type identifyCode struct {
	identify string
	code     string
}

type sendFunc func()

var sendQueue = make(chan sendFunc, 100)

func init() {
	for i := 0; i < 3; i++ {
		go handleSendQueue(sendQueue)
	}
}

func handleSendQueue(queue <-chan sendFunc) {
	for {
		doSend := <-queue
		doSend()
	}
}
