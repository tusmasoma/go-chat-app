import React from 'react';
import { Button, Textarea } from "@nextui-org/react";
import { Send } from "lucide-react";
// もし useChat を使わない場合は、カスタムフックを使用して状態を管理
// import { useState } from 'react';

interface ChatInputProps {
  input: string;
  handleInputChange: (e: React.ChangeEvent<any>) => void;
  handleSubmit: (e: React.FormEvent) => void;
  setInput: (value: string) => void;
}

export const ChatInput: React.FC<ChatInputProps> = ({ handleInputChange, handleSubmit, input, setInput }) => {
  return (
    <div className="z-10 bg-zinc-900 absolute bottom-0 left-0 w-full">
      <div className="mx-2 flex flex-row gap-3 md:mx-4 md:last:mb-6 lg:mx-auto lg:max-w-2xl xl:max-w-3xl">
        <div className="relative flex h-full flex-1 items-stretch md:flex-col">
          <div className="relative flex flex-col w-full flex-grow p-4">
            <form onSubmit={handleSubmit} className="relative">
              <Textarea
                minRows={4}
                autoFocus
                onChange={handleInputChange}
                value={input}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    handleSubmit(e);
                    setInput("");
                  }
                }}
                placeholder="Enter your question..."
                className="resize-none bg-zinc-800 hover:bg-zinc-900 rounded-xl text-base"
              />

              <Button
                size="sm"
                type="submit"
                className="absolute z-10 border border-border bg-zinc-900 right-2 bottom-2"
              >
                <Send className="size-4" />
              </Button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};
