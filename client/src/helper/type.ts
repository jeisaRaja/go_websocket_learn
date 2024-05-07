export type EventWs = {
  type: string;
  payload: {
    id?: string;
    message: string;
    room: string;
    from_id?: string;
    from_name: string;
    sent?: Date;
  };
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
