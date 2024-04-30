import { FormEvent, useEffect, useRef, useState } from "react";
import { newEventWs, newUser } from "../helper/objectFactories";
import { EventWs, UserAuth } from "../helper/type";

const Home = () => {
  const textArea = useRef<HTMLTextAreaElement | null>(null);
  const [message, setMessage] = useState("");
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
    const res = await fetch("https://localhost:3000/login", {
      method: "POST",
      body: JSON.stringify(authBody),
    });
    if (res.status === 200) {
      const data = await res.json();
      const user = newUser(username);
      setOTP(data.otp);
      setAuth(user);
    }
  };

  const routeEvent = (event: EventWs) => {
    console.log(event);
    if (event.type === undefined) {
      return;
    }
    switch (event.type) {
      case "new_message":
        console.log("new_message", event.payload);
        if (textArea.current) {
          const currVal = textArea.current.value;
          textArea.current.value =
            currVal +
            "\n" +
            event.payload.message +
            " \n -" +
            event.payload.from;
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
    const eventData = newEventWs(type, message, username);
    const jsonData = JSON.stringify(eventData);
    console.log(jsonData);
    conn.send(jsonData);
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
      return false;
    }
    if (!conn) {
      return;
    }
    sendMessage("change_room", inputChatroom);
    setChatroom(inputChatroom);
    setInputChatroom("");
  };

  const isWebSocketSupported = "WebSocket" in window;

  useEffect(() => {
    console.log(OTP);
    if (
      isWebSocketSupported &&
      conn === null &&
      auth !== undefined &&
      OTP !== undefined
    ) {
      console.log("new websocket");
      setConn(() => new WebSocket(`wss://localhost:3000/ws?otp=${OTP}`));
    }
    return () => {
      conn?.close();
    };
  }, [conn, isWebSocketSupported, auth, OTP]);

  useEffect(() => {
    if (conn) {
      const handleMessage = (ev: MessageEvent) => {
        console.log(ev);
        if (ev.type === "ping") {
          console.log("ping received");
        }
        const eventData = JSON.parse(ev.data) as EventWs;
        routeEvent(eventData);
      };
      conn.onmessage = handleMessage;
      console.log("message");
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

      <textarea
        name="messagearea"
        id="messagearea"
        cols={30}
        rows={10}
        placeholder="welcome to chatgo"
        readOnly
        className="p-3"
        ref={textArea}
      ></textarea>

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
