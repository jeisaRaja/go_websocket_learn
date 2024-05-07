type ChatProps = {
  msg: string;
  uname: string;
  sent: Date;
};
export default function Chat({ msg, uname, sent }: ChatProps) {
  const sentDate = new Date(sent);
  const day = sentDate.getDate();
  const month = sentDate.getMonth();
  const year = sentDate.getFullYear();

  return (
    <div className="chat-element">
      <div className="message">{msg}</div>
      <div className="username">{uname}</div>
      <div className="time">{`${day}/${month}/${year}`}</div>
    </div>
  );
}
