package main

import (
	"encoding/json"
	"jeisaraja/websocket_learn/models"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChangeRoom  = "change_room"
	EventAnnounce    = "announce"
)

type SendMessageEvent struct {
	models.Chat
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Room string `json:"room"`
}

type Room struct {
	Name   string   `json:"room"`
	Member []string `json:"member"`
}
