import React from 'react';
import { Message as MessageComponent } from './Message';
import { MessageSquare } from 'lucide-react';

interface TMessage {
  content: string;
  role: 'user' | 'system' | 'assistant'; // 例: roleを必要に応じて定義
}

interface MessagesProps {
  messages: TMessage[];
}

export const Messages: React.FC<MessagesProps> = ({ messages }) => {
  return (
    <div className="flex max-h-[calc(100vh-3.5rem-7rem)] flex-1 flex-col overflow-y-auto">
      {messages.length ? (
        messages.map((message, i) => (
          <MessageComponent
            key={i}
            content={message.content}
            isUserMessage={message.role === 'user'}
          />
        ))
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center gap-2">
          <MessageSquare className="size-8 text-blue-500" />
          <h3 className="font-semibold text-xl text-white">You're all set!</h3>
          <p className="text-zinc-500 text-sm">Ask your first question to get started.</p>
        </div>
      )}
    </div>
  );
};
