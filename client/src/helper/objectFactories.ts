import { EventWs } from "./type";

export function newEventWs(type: string, payload: string): EventWs {
  return {
    type, payload
  }
}