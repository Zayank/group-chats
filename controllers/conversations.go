package controllers

import (
	"github.com/Massad/gin-boilerplate/helpers"
	"github.com/dustin/go-broadcast"
)

type ConversationListenerStruct struct {
	Time      string
	UserId    string
	UserName  string
	GroupUuid string
	GroupName string
	Text      string
}

type ConversationListener struct {
	id   string
	Chan chan interface{}
}

type ConversationManager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *ConversationListener
	close        chan *ConversationListener
	delete       chan string
	messages     chan *ConversationListenerStruct
}

func NewConversationManager() *ConversationManager {
	manager := &ConversationManager{
		roomChannels: make(map[string]broadcast.Broadcaster),
		open:         make(chan *ConversationListener, 100),
		close:        make(chan *ConversationListener, 100),
		delete:       make(chan string, 100),
	}

	go manager.run()
	return manager
}

func (m *ConversationManager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case groupuuid := <-m.delete:
			m.deleteBroadcast(groupuuid)
		case message := <-m.messages:
			m.room(message.GroupUuid).Submit(message.GroupUuid + "|" + message.GroupName)
		}
	}
}

func (m *ConversationManager) register(listener *ConversationListener) {
	m.room(listener.id).Register(listener.Chan)
}

func (m *ConversationManager) deregister(listener *ConversationListener) {
	m.room(listener.id).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *ConversationManager) deleteBroadcast(groupuuid string) {
	b, ok := m.roomChannels[groupuuid]
	if ok {
		b.Close()
		delete(m.roomChannels, groupuuid)
	}
}

func (m *ConversationManager) room(groupuuid string) broadcast.Broadcaster {
	b, ok := m.roomChannels[groupuuid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.roomChannels[groupuuid] = b
	}
	return b
}

func (m *ConversationManager) OpenListener(groupuuid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &ConversationListener{
		id:   groupuuid,
		Chan: listener,
	}
	return listener
}

func (m *ConversationManager) CloseListener(groupuuid string, channel chan interface{}) {
	m.close <- &ConversationListener{
		id:   groupuuid,
		Chan: channel,
	}
}

func (m *ConversationManager) DeleteBroadcast(groupuuid string) {
	m.delete <- groupuuid
}

func (m *ConversationManager) Submit(userid, username, groupuuid, time string, text string) {

	msg := &ConversationListenerStruct{
		Time:      time,
		UserId:    userid,
		UserName:  username,
		GroupUuid: groupuuid,
	}

	m.messages <- msg
}

func (m *ConversationManager) PushConversation(userid, username, groupuuid, groupname, time string, text string) {

	m.room(helpers.GetConversationsId(userid)).Submit(groupuuid + "|" + groupname)

}
