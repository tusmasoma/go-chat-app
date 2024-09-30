import React from "react";
import Markdown from "./markdown";
import { Bot, User2 } from "lucide-react";

export type CustomMessage = {
  id: string;
  user_id: string;
  workspace_id: string;
  target_id: string;
  membership_id?: string;
  text: string;
  role: "user" | "bot" | "assistant" | "system" | "function" | "tool";
  timestamp?: string;
  action?: string;
};


type Props = {
  messages: CustomMessage[];
  isLoading: boolean;
};

const Messages = ({ messages, isLoading }: Props) => {
  return (
    <div
      id="chatbox"
      className="flex flex-col w-full text-left mt-4 gap-4 whitespace-pre-wrap"
    >
      {messages.map((m, index) => {
        return (
          <div
            key={m.id}  // メッセージの一意なIDを使用
            className={`p-4 shadow-md rounded-3xl relative w-fit max-w-[80%] break-words ${
              m.role === "user"
                ? "bg-purple-600 text-white ml-auto mr-5"
                : "bg-gray-300 text-black ml-5"
            }`}
          >
            <Markdown text={m.text} />
            {m.role === "user" ? (
              <User2 className="absolute -right-10 top-2 border rounded-full p-1 shadow-lg bg-purple-600 text-white" />
            ) : (
              <Bot
                className={`absolute top-2 -left-10 border rounded-full p-1 shadow-lg stroke-[#0842A0] bg-gray-300 text-black ${
                  isLoading && index === messages.length - 1 ? "animate-bounce" : ""
                }`}
              />
            )}
            {/* 任意でタイムスタンプなども表示可能 */}
            {m.timestamp && <div className="text-xs text-gray-500">{m.timestamp}</div>}
          </div>
        );
      })}
    </div>
  );
};

export default Messages;
