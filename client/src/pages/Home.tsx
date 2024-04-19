import { FormEvent, useState } from "react";

const Home = () => {
  const [message, setMessage] = useState("");
  const [chatroom, setChatroom] = useState("general");
  const [inputChatroom, setInputChatroom] = useState("");
  const onMessageSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log(message);
    if(message===''){
      return
    }
    conn.send(message)
  };
  const onChangeChatroom = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!inputChatroom || chatroom == inputChatroom) {
      return false;
    }
    setChatroom(inputChatroom);
    setInputChatroom("");
  };

  const isWebSocketSupported = 'WebSocket' in window;
  const conn = new WebSocket("ws://localhost:3000/ws")
  
  setTimeout(() => {
    if (conn.readyState !== WebSocket.OPEN){
      conn.close()
      console.log("websocket failed to connect")
    }
  }, 2000);
  return (
    !isWebSocketSupported ? <p>WebSockets are not supported in this browser.</p> :
    <div className="flex flex-col bg-white p-5">
      <h1 className="">Chatgo</h1>
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
