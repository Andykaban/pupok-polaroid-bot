package bot

import (
	"sync"
	"time"
)

const maxExpirationDelta = 60000000000

type MessageItem struct {
	message string
	expiration int64
}

type MessagesStorage struct {
	lock *sync.RWMutex
	messages map[int64]MessageItem
}

func NewMessagesStorage() *MessagesStorage {
	messages := make(map[int64]MessageItem)
	return &MessagesStorage{lock:&sync.RWMutex{}, messages:messages}
}

func (m *MessagesStorage) GetMessage(chatId int64) string {
	msg := ""
	m.lock.RLock()
	defer m.lock.RUnlock()
	if item, ok := m.messages[chatId]; ok {
		msg = item.message
	}
	return msg
}

func (m *MessagesStorage) SetMessage(chatId int64, msg string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	item := MessageItem{message:msg, expiration:time.Now().UnixNano()}
	m.messages[chatId] = item
}

func (m *MessagesStorage) RemoveExpiredMessages() {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentTime := time.Now().UnixNano()
	for key,value := range m.messages {
		startTime := value.expiration
		if currentTime - startTime > maxExpirationDelta {
			delete(m.messages, key)
		}
	}
}
