import { FormEvent, useEffect, useRef, useState } from "react";
import {
  newChangeRoomWs,
  newEventWs,
  newUser,
} from "../helper/objectFactories";
import { Chat as ChatType, EventWs, UserAuth } from "../helper/type";
import { scrollToBottom } from "../helper/autoScroll";
import Chat from "../components/Chat";

const Home = () => {
  const messageArea = useRef<HTMLDivElement | null>(null);
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<Array<ChatType>>([]);
  const [chatroom, setChatroom] = useState("general");
  const [inputChatroom, setInputChatroom] = useState("");
  const [conn, setConn] = useState<null | WebSocket>(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [auth, setAuth] = useState<undefined | UserAuth>(undefined);
  const [OTP, setOTP] = useState("");

  const handleAuth = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const authBody = {
      username: username,
      password: password,
    };
    const res = await fetch("http://localhost:5000/login", {
      method: "POST",
      body: JSON.stringify(authBody),
    });
    if (res.status === 200) {
      const resData = await res.json();
      const user = newUser(username);
      setOTP(resData.data.otp);
      setAuth(user);
    }
  };

  const routeEvent = (event: EventWs) => {
    if (event.type === undefined) {
      return;
    }
    switch (event.type) {
      case "new_message":
        if (messageArea.current) {
          if (!event.payload.sent) {
            return;
          }
          const newMessage: ChatType = {
            message: event.payload.message,
            from: event.payload.from_name,
            sent: event.payload.sent,
          };
          setMessages((prev) => [...prev, newMessage]);
        }
        break;
      default:
        console.log("unsupported message type");
        break;
    }
  };

  const sendMessage = (type: string, message: string) => {
    if (!conn) {
      return;
    }
    switch (type) {
      case "change_room": {
        const eventData = newChangeRoomWs(inputChatroom);
        const jsonData = JSON.stringify(eventData);
        conn.send(jsonData);
        break;
      }
      default: {
        const eventData = newEventWs(type, message, chatroom, username);
        const jsonData = JSON.stringify(eventData);
        console.log(jsonData);
        conn.send(jsonData);
        break;
      }
    }
  };

  const onMessageSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (message === "") {
      return;
    }
    if (conn === null) {
      return;
    }
    sendMessage("send_message", message);
    setMessage("");
  };

  const onChangeChatroom = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!inputChatroom || chatroom == inputChatroom) {
      return;
    }
    if (!conn) {
      return;
    }
    sendMessage("change_room", inputChatroom);
    setMessages([])
    setChatroom(inputChatroom);
    setInputChatroom("");
  };

  const isWebSocketSupported = "WebSocket" in window;

  useEffect(() => {
    scrollToBottom("messagearea");
  }, [messages]);

  useEffect(() => {
    if (
      isWebSocketSupported &&
      conn === null &&
      auth !== undefined &&
      OTP !== undefined
    ) {
      console.log("new websocket");
      setConn(() => new WebSocket(`ws://localhost:5000/ws?otp=${OTP}`));
    }
    return () => {
      conn?.close();
    };
  }, [conn, isWebSocketSupported, auth, OTP]);

  useEffect(() => {
    if (conn) {
      const handleMessage = (ev: MessageEvent) => {
        //if (ev.type === "ping") {
        // console.log("ping received");
        //}
        const eventData = JSON.parse(ev.data) as EventWs;
        routeEvent(eventData);
      };
      conn.onmessage = handleMessage;
      return () => {
        conn.onmessage = null;
      };
    }
  }, [conn]);

  setTimeout(() => {
    if (conn && conn.readyState !== WebSocket.OPEN) {
      conn.close();
      console.log("websocket failed to connect");
    }
  }, 10000);

  return !isWebSocketSupported ? (
    <p>WebSockets are not supported in this browser.</p>
  ) : (
    <div className="flex flex-col bg-white p-5">
      <h1 className="">Chatgo</h1>
      <form action="" className="flex gap-3" onSubmit={(e) => handleAuth(e)}>
        <label htmlFor="username">Username</label>
        <input
          type="text"
          name="username"
          id="username"
          onChange={(e) => setUsername(e.target.value)}
          value={username}
          className="border border-grey-200 p-2"
        />
        <label htmlFor="password">Password</label>
        <input
          type="password"
          name="password"
          id="password"
          onChange={(e) => setPassword(e.target.value)}
          value={password}
          className="border border-grey-200 p-2"
        />
        <button type="submit">login</button>
      </form>
      <h3>Currently in chat: {chatroom}</h3>
      <form
        action=""
        className="flex items-center gap-3 my-3"
        onSubmit={(e) => onChangeChatroom(e)}
      >
        <label htmlFor="chatroom">Chatroom: </label>
        <input
          type="text"
          name="chatroom"
          id="chatroom"
          className="py-2 px-3 outline-none border-2 border-gray-200"
          onChange={(e) => setInputChatroom(e.target.value)}
          value={inputChatroom}
        />
        <input
          type="submit"
          value="Change chatroom"
          className="cursor-pointer py-2 px-3 bg-gray-200 rounded-md hover:bg-gray-300"
        />
      </form>

      <div
        id="messagearea"
        className="w-full p-3 flex flex-col h-[400px] overflow-y-auto"
        ref={messageArea}
      >
        {messages.map((item) => (
          <Chat msg={item.message} uname={item.from} sent={item.sent} />
        ))}
      </div>

      <form
        onSubmit={(e) => onMessageSubmit(e)}
        className="flex items-center gap-3 my-3"
      >
        <label htmlFor="message">Message: </label>
        <input
          type="text"
          name="message"
          id="message"
          onChange={(e) => setMessage(e.target.value)}
          className="py-2 px-3 outline-none border-2 border-gray-200"
          value={message}
        />
        <input
          type="submit"
          value="send"
          className="cursor-pointer py-2 px-3 bg-gray-200 rounded-md hover:bg-gray-300"
        />
      </form>
    </div>
  );
};

export default Home;
