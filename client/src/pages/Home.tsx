import React, { useState } from 'react';
import InputForm from '../components/inputForm';
import Messages,{ CustomMessage } from '../components/messages';

const Home = () => {
    const [messages, setMessages] = useState<CustomMessage[]>([]); // メッセージのリスト
    const [input, setInput] = useState(''); // 入力内容
    const [isLoading, setIsLoading] = useState(false); // ローディング状態

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

  return (
    <main className="flex min-h-screen flex-col items-center p-12 text-lg">
      <InputForm
        input={input}
        handleInputChange={handleInputChange}
        handleSubmit={handleSubmit}
        isLoading={isLoading}
        stop={stop}
      />
      <Messages messages={messages} isLoading={isLoading} />
    </main>
  );
};

export default Home;
