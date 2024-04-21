import { EventWs, UserAuth } from "./type";

export function newEventWs(type: string, message: string, from: string): EventWs {
  return {
    type, payload: { message, from }
  }
}

export function newUser(username: string): UserAuth {
  return {
    username
  }
}