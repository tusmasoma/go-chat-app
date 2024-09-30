import React, { useState, useEffect, useRef } from 'react';
import InputForm from '../components/chat/inputForm';
import Messages,{ CustomMessage } from '../components/chat/messages';

const Chat = () => {
    const [messages, setMessages] = useState<CustomMessage[]>([]); // メッセージのリスト
    const [input, setInput] = useState(''); // 入力内容
    const [isLoading, setIsLoading] = useState(false); // ローディング状態
    const [bottomPadding, setBottomPadding] = useState(0); // メッセージの高さに基づいた余白
    const lastMessageRef = useRef<HTMLDivElement>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    // 入力内容が変わったときのハンドラー
    const handleInputChange = (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
      ) => {
        setInput(event.target.value);
    };

    // フォーム送信時のハンドラー
    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      if (!input.trim()) return;

      const newMessage: CustomMessage = {
        id: String(messages.length + 1), // IDを生成
        content: input, // inputをcontentに代入
        role: "user", // ロールを設定
        timestamp: new Date().toISOString(), // タイムスタンプを生成
      };

      // ローディング状態にセット
      setIsLoading(true);

      // 仮のメッセージ送信ロジック
      setTimeout(() => {
        setMessages((prevMessages) => [...prevMessages, newMessage]); // ユーザーメッセージを追加
        setInput(''); // 入力をクリア
        setIsLoading(false); // ローディング終了

        // 固定メッセージを1秒後に返却
        const fixedMessage: CustomMessage = {
          id: String(messages.length + 2), // 固定メッセージ用のIDを生成
          content: "これは固定メッセージです。", // 固定の返答メッセージ
          role: "bot", // ボットからのメッセージとして設定
          timestamp: new Date().toISOString(),
        };

        setTimeout(() => {
          setMessages((prevMessages) => [...prevMessages, fixedMessage]);
        }, 1000); // 1秒後に固定メッセージを返す
      }, 1000); // 1秒後にメッセージを送信するように見せる
    };

    // メッセージの送信を止める（仮のロジック）
    const stop = () => {
      setIsLoading(false);
    };

    useEffect(() => {
      if (lastMessageRef.current) {
        // 最新メッセージの高さに基づいて余白を設定
        messagesEndRef.current!.style.paddingBottom = `100px`;
      }
    }, [messages]);

    // 最新メッセージにスクロール
    const scrollToBottom = () => {
      if (messagesEndRef.current) {
        messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
      }
    };

    useEffect(() => {
      scrollToBottom();
    }, [messages]);

  return (
    <main className="flex flex-col items-center min-h-screen p-12 text-lg">
      <div className="flex flex-col w-full max-w-2xl h-full flex-grow overflow-y-auto" style={{ paddingBottom: `${bottomPadding}px` }}>
        <div>
          <Messages messages={messages} isLoading={isLoading} /> {/* メッセージ全体を1回だけ表示 */}
        </div>
        <div ref={lastMessageRef}>
          {/* 最後のメッセージ位置に ref を設定 */}
          <div ref={lastMessageRef} /> {/* 最新メッセージの位置を監視 */}
        </div>
        <div ref={messagesEndRef} /> {/* 最新メッセージ位置の参照 */}
      </div>
      <InputForm
        input={input}
        handleInputChange={handleInputChange}
        handleSubmit={handleSubmit}
        isLoading={isLoading}
        stop={() => setIsLoading(false)}
      />
    </main>
  );
};

export default Chat;
