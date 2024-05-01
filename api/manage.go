package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Println(origin)
			return origin == "http://localhost:5173"
		},
	}
)

const (
	username = "testing"
	password = "password"
)

type Manager struct {
	clients  ClientList
	handlers map[string]EventHandler
	sync.RWMutex
	otps OTPMap
	DB   *sql.DB
}

func NewManager(ctx context.Context, db *sql.DB) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewOTPMap(ctx, 20*time.Second),
    DB: db,
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) serveWs(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("unauthorized!")
		return
	}
	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println((err))
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = sendMessage
	m.handlers[EventChangeRoom] = changeRoom
}

func sendMessage(event Event, c *Client) error {
	var msgEvent SendMessageEvent
	log.Println(event)
	if err := json.Unmarshal(event.Payload, &msgEvent); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}
	var msg NewMessageEvent
	msg.Sent = time.Now()
	msg.From = msgEvent.From
	msg.Message = msgEvent.Message

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}
	var sendEvent Event
	sendEvent.Payload = data
	sendEvent.Type = EventNewMessage

	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- sendEvent
		}
	}
	return nil
}

func changeRoom(event Event, c *Client) error {
	var roomEvent ChangeRoomEvent
	err := json.Unmarshal(event.Payload, &roomEvent)
	if err != nil {
		return fmt.Errorf("bad payload")
	}
	c.chatroom = roomEvent.Room
	return nil
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var reqBody = UserAuth{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	if reqBody.Username == username && reqBody.Password == password {
		type response struct {
			OTP string `json:"otp"`
		}
		otp := m.otps.NewOTP()
		resp := response{
			OTP: otp.Key,
		}
		data, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
