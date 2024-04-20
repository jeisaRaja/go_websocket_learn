import { EventWs, UserAuth } from "./type";

export function newEventWs(type: string, payload: string): EventWs {
  return {
    type, payload
  }
}

export function newUser(username: string): UserAuth {
  return {
    username
  }
}