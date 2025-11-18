package faulttolerance

type DeadLetterQueue struct {
	messages []string
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		messages: make([]string, 0),
	}
}

func (s *DeadLetterQueue) GetMessages() []string {
	return s.messages
}

func ProcessWithDLQ(messages []string, handler func(msg string) error, dlq *DeadLetterQueue) {
	var err error

	for _, msg := range messages {
		err = handler(msg)
		if err != nil {
			dlq.messages = append(dlq.messages, msg)
		}
	}
}
