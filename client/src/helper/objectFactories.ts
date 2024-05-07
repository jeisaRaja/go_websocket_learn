import { EventWs, UserAuth } from "./type";

export function newEventWs(
  type: string,
  message: string,
  room: string,
  from_name: string,
): EventWs {
  return {
    type,
    payload: { message, room, from_name },
  };
}

export function newUser(username: string): UserAuth {
  return {
    username,
  };
}
