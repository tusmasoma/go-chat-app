import React, { useState, useEffect } from 'react';
import { Message } from '../types/Message'; // メッセージの型は適宜定義
import { Messages } from './Messages';
import { ChatInput } from './ChatInput';

interface ChatWrapperProps {
  sessionId: string;
  initialMessages: Message[];
}

export const ChatWrapper: React.FC<ChatWrapperProps> = ({
  sessionId,
  initialMessages,
}) => {
  const [messages, setMessages] = useState<Message[]>(initialMessages);
  const [input, setInput] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<any>) => {
    setInput(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // メッセージの送信ロジックをここに実装
    const newMessage: Message = {
      id: Date.now().toString(),
      content: input,
      role: 'user', // 送信者を指定
    };

    // メッセージの状態を更新
    setMessages((prevMessages) => [...prevMessages, newMessage]);

    // 送信後に入力をクリア
    setInput('');
  };

  useEffect(() => {
    // sessionIdを使ってサーバーからのメッセージ取得ロジックを実装
  }, [sessionId]);

  return (
    <div className="relative min-h-full bg-zinc-900 flex divide-y divide-zinc-700 flex-col justify-between gap-2" style={{marginTop: "100px"}}>
      <div className="flex-1 text-black bg-zinc-800 justify-between flex flex-col">
        <Messages messages={messages} />
      </div>

      <ChatInput
        input={input}
        handleInputChange={handleInputChange}
        handleSubmit={handleSubmit}
        setInput={setInput}
      />
    </div>
  );
};
