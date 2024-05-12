export type EventWs = {
  type: string;
  payload: NewMessagePayload | SendMessagePayload | AnnouncementPayload;
};

export type NewMessagePayload = {
  type: "new_message";
  id: string;
  message: string;
  room: string;
  from_id: string;
  from_name: string;
  sent: Date;
};

export type SendMessagePayload = {
  type: "send_message";
  room: string;
  from_name: string;
  message: string;
};

export type AnnouncementPayload = {
  type: "announce";
  member: [string];
};

export type ChangeRoom = {
  type: string;
  payload: {
    room: string;
  };
};

export type Chat = {
  message: string;
  from: string;
  sent: Date;
};

export type UserAuth = {
  username: string;
};

export type Announce = {
  room: string;
  member: [string];
};
