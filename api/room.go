package main

import (
	"encoding/json"
	"fmt"
)

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

	err = addRoomMember(c)
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

func addRoomMember(c *Client) error {
	var ann Room
	c.manager.RoomMap[c.chatroom] = append(c.manager.RoomMap[c.chatroom], c.username)
	ann.Name = c.chatroom
	ann.Member = c.manager.RoomMap[c.chatroom]
	data, err := json.Marshal(ann)
	if err != nil {
		return err
	}
	var event Event
	event.Type = EventAnnounce
	event.Payload = data

	fmt.Println(c.manager.RoomMap)

	for client := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- event
		}
	}
	fmt.Println("sending announce join room")
	return nil
}
