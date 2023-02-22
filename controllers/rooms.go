package controllers

import (
	"github.com/dustin/go-broadcast"

	"github.com/Massad/gin-boilerplate/helpers"
)

type Message struct {
	Time      string
	UserId    string
	UserName  string
	GroupUuid string
	Text      string
}

type Listener struct {
	GroupUuid string
	Chan      chan interface{}
}

type Manager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
}

type User struct {
	Name     string
	Id       string
	Contacts []string
}

func NewRoomManager() *Manager {
	manager := &Manager{
		roomChannels: make(map[string]broadcast.Broadcaster),
		open:         make(chan *Listener, 100),
		close:        make(chan *Listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *Message, 100),
	}

	go manager.run()
	return manager
}

func (m *Manager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case groupuuid := <-m.delete:
			m.deleteBroadcast(groupuuid)
		case message := <-m.messages:
			m.room(message.GroupUuid).Submit(message.UserName + "|" + message.Text + "|" + message.Time + "|" + message.UserId + "|" + message.GroupUuid)
		}
	}
}

func (m *Manager) register(listener *Listener) {
	m.room(listener.GroupUuid).Register(listener.Chan)
}

func (m *Manager) deregister(listener *Listener) {
	m.room(listener.GroupUuid).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *Manager) deleteBroadcast(groupuuid string) {
	b, ok := m.roomChannels[groupuuid]
	if ok {
		b.Close()
		delete(m.roomChannels, groupuuid)
	}
}

func (m *Manager) room(groupuuid string) broadcast.Broadcaster {
	b, ok := m.roomChannels[groupuuid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.roomChannels[groupuuid] = b
	}
	return b
}

func (m *Manager) OpenListener(groupuuid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		GroupUuid: groupuuid,
		Chan:      listener,
	}
	return listener
}

func (m *Manager) CloseListener(groupuuid string, channel chan interface{}) {
	m.close <- &Listener{
		GroupUuid: groupuuid,
		Chan:      channel,
	}
}

func (m *Manager) DeleteBroadcast(groupuuid string) {
	m.delete <- groupuuid
}

func (m *Manager) Submit(userid, username, groupuuid, time string, text string) {

	msg := &Message{
		Time:      time,
		UserId:    userid,
		UserName:  username,
		GroupUuid: groupuuid,
		Text:      text,
	}

	m.messages <- msg
}

func (m *Manager) PushConversation(userid, username, groupuuid, time string, text string) {

	m.room(helpers.GetConversationsId(userid)).Submit(username + "|" + text + "|" + time + "|" + userid + "|" + groupuuid)

}
