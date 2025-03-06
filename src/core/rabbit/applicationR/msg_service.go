package applicationR

type MessagePublisher interface {
	PublishMessage(queueName string, message []byte) error
}

type MessageService struct {
	messagePublisher MessagePublisher
}

func NewMessageService(publisher MessagePublisher) *MessageService {
	return &MessageService{messagePublisher: publisher}
}

func (service *MessageService) PublishMessage(queueName string, message string) error {
	err := service.messagePublisher.PublishMessage(queueName, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
