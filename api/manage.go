package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jeisaraja/websocket_learn/database"
	"jeisaraja/websocket_learn/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

type Manager struct {
	clients  ClientList
	handlers map[string]EventHandler
	sync.RWMutex
	otps     OTPMap
	DB       *database.Queries
	Validate *validator.Validate
	RoomMap  models.RoomMap
}

func NewManager(ctx context.Context, db *database.Queries) *Manager {
	validate := validator.New()
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
		otps:     NewOTPMap(ctx, 20*time.Second),
		DB:       db,
		Validate: validate,
		RoomMap:  models.RoomMap{},
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) serveWs(w http.ResponseWriter, r *http.Request) {
	otp := r.URL.Query().Get("otp")
	username := r.URL.Query().Get("u")
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
	client := NewClient(conn, m, username)
	m.addClient(client)

	chats, err := m.DB.LoadChats("general")
	if err != nil {
		return
	}

	go client.readMessages()
	go client.writeMessages()

	for _, chat := range chats {
		jsonChat, err := json.Marshal(chat)
		if err != nil {
			return
		}
		var sendEvent Event
		sendEvent.Payload = jsonChat
		sendEvent.Type = EventNewMessage

		client.egress <- sendEvent
	}
	err = announceJoinRoom(client)
	if err != nil {
		return
	}
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

func announceJoinRoom(c *Client) error {
	var ann JoinRoom
  c.manager.RoomMap[c.chatroom] = append(c.manager.RoomMap[c.chatroom], c.username)
	ann.Room = c.chatroom
  ann.Member = c.manager.RoomMap[c.chatroom]
	data, err := json.Marshal(ann)
	if err != nil {
		return err
	}
	var event Event
	event.Type = EventAnnounce
	event.Payload = data

	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- event
		}
	}
	fmt.Println("sending announce join room")
	return nil
}

func sendMessage(event Event, c *Client) error {
	var msgEvent models.Chat
	if err := json.Unmarshal(event.Payload, &msgEvent); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}

	msgEvent.Room = c.chatroom
	msgEvent.Sent = time.Now()
	msgEvent.ID = uuid.New()
	err := c.manager.DB.InsertChat(&msgEvent)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	data, err := json.Marshal(msgEvent)
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

	chats, err := c.manager.DB.LoadChats(roomEvent.Room)
	if err != nil {
		return err
	}

	err = announceJoinRoom(c)
	if err != nil {
		return err
	}

	for _, chat := range chats {
		jsonChat, err := json.Marshal(chat)
		if err != nil {
			return err
		}
		var sendEvent Event
		sendEvent.Payload = jsonChat
		sendEvent.Type = EventNewMessage

		c.egress <- sendEvent
	}
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

func (m *Manager) signup(w http.ResponseWriter, r *http.Request) {
	var reqBody = models.AuthSignup{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(reqBody.Username) < 4 {
		http.Error(w, "Username must be at least 4 characters long", 400)
		return
	}
	if len(reqBody.InputPassword) < 8 {
		http.Error(w, "Password must be at least 8 characters long", 400)
		return
	}
	if len(reqBody.Email) < 8 {
		http.Error(w, "Email must be at least 8 characters long", 400)
		return
	}

	var user = models.User{}
	user.ID = uuid.New()
	user.Username = reqBody.Username

	err = user.Password.Set(reqBody.InputPassword)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}
	user.Email = reqBody.Email

	err = m.DB.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	jwtKey := []byte("JWT_KEY")
	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := models.Response{
		Error: "",
		Msg:   "Account Created.",
		Data:  map[string]string{"token": tokenStr},
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)
}

func (m *Manager) login(w http.ResponseWriter, r *http.Request) {
	var reqBody = models.AuthSignin{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = m.Validate.Struct(reqBody)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	password := m.DB.GetUserPassword(reqBody.Username)

	if password == nil {
		http.Error(w, "User does not exists", 400)
		return
	}
	valid, err := ComparePasswords(*password, reqBody.InputPassword)
	if err != nil {
		http.Error(w, "User does not exists", 400)
		return
	}

	if valid {
		type response struct {
			OTP string `json:"otp"`
		}
		otp := m.otps.NewOTP()
		res := models.Response{
			Error: "",
			Msg:   "Login Successful.",
			Data:  map[string]string{"otp": otp.Key},
		}
		jsonRes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonRes)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized!"))
		return
	}
}
