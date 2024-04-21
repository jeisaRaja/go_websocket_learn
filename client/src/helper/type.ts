export type EventWs = {
  type: string
  payload: {
    message: string
    from: string
  }
}

export type UserAuth = {
  username: string
}