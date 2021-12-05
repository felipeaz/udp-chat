package messages

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
	"udp-chat/internal/app/chat/entity"
	error_messages "udp-chat/internal/app/chat/server/constants"
	database "udp-chat/internal/cache"
	"udp-chat/internal/logger"
)

type Message struct {
	Cache  database.CacheInterface
	Logger logger.LogInterface
	Size   int
}

func NewMessage(cache database.CacheInterface, log logger.LogInterface, maxSize int) Message {
	return Message{
		Cache:  cache,
		Logger: log,
		Size:   maxSize,
	}
}

func (m Message) Store(msg string) (*entity.Message, error) {
	var msgObj entity.Message
	messages, err := m.Get()
	if err != nil {
		err = errors.Wrap(err, error_messages.FailedToGetMessagesFromChat)
		m.Logger.Error(err)
	}

	bmsg := []byte(msg)
	err = json.Unmarshal(bmsg, &msgObj)
	if err != nil {
		return nil, err
	}
	msgObj.Date = time.Now()

	messages = m.addMessageToQueue(messages, msgObj)

	b, err := json.Marshal(messages)
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	err = m.Cache.Set("CHAT", b)
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	return &msgObj, nil
}

func (m Message) Get() ([]entity.Message, error) {
	var messages []entity.Message
	b, err := m.Cache.Get("CHAT")
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	err = json.Unmarshal(b, &messages)
	if err != nil {
		m.Logger.Error(err)
		return nil, err
	}

	return messages, nil
}

func (m Message) addMessageToQueue(queue []entity.Message, msg entity.Message) []entity.Message {
	var newQueue []entity.Message
	if len(queue) < m.Size {
		newQueue = append(queue, msg)
		return newQueue
	}

	newQueue = append(queue[1:], msg)
	return newQueue
}

func (m Message) Delete(id string) error {
	return nil
}
