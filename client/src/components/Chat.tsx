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
    <div className="chat-element w-full flex flex-col my-3">
      <div className="message bg-slate-200 p-3 rounded-md">{msg}</div>
      <div className="flex gap-3 px-3">
        <div className="username">{uname}</div>
        <div className="time">{`${day}/${month}/${year}`}</div>
      </div>
    </div>
  );
}
