import { ChangeRoom, EventWs, SendMessagePayload, UserAuth } from "./type";

export function newEventWs(
  type: string,
  message: string,
  room: string,
  from_name: string,
): EventWs {
  return {
    type,
    payload: { message, room, from_name } as SendMessagePayload,
  };
}

export function newChangeRoomWs(room: string): ChangeRoom {
  return {
    type: "change_room",
    payload: { room },
  };
}

export function newUser(username: string): UserAuth {
  return {
    username,
  };
}
