import React, { useState, useEffect, useRef } from 'react';
import InputForm from '../components/chat/inputForm';
import Messages,{ CustomMessage } from '../components/chat/messages';

let workspaceID: string = "550e8400-e29b-41d4-a716-446655440000";
let channelID: string = "123e4567-e89b-12d3-a456-426614174000";

const decodeToken = (token: string) => {
  // トークンをドットで分割
  const parts = token.split('.');

  if (parts.length !== 3) {
    throw new Error("Invalid token format");
  }

  // ヘッダー部分とペイロード部分をデコード
  const header = JSON.parse(atob(parts[0].replace(/-/g, '+').replace(/_/g, '/')));
  const payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));

  console.log("Decoded Header:", header);  // ヘッダーの内容を表示
  console.log("Decoded Payload:", payload);  // ペイロードの内容を表示

  return payload; // 必要に応じてペイロードの内容を返す
};

const Chat = () => {
  const [userID, setUserID] = useState<string>(''); // ユーザーID
  const [messages, setMessages] = useState<CustomMessage[]>([]); // メッセージのリスト
  const [input, setInput] = useState(''); // 入力内容
  const [isLoading, setIsLoading] = useState(false); // ローディング状態
  const [bottomPadding, setBottomPadding] = useState(0); // メッセージの高さに基づいた余白
  const lastMessageRef = useRef<HTMLDivElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const ws = useRef<WebSocket | null>(null); // WebSocketインスタンス

  // WebSocketの接続
  useEffect(() => {
    const token = localStorage.getItem('jwtToken');
    if (!token) {
      console.error('JWT token not found');
      return;
    }

    // JWTデコードしてuserIDをセット
    const decodedPayload = decodeToken(token);
    setUserID(decodedPayload.userId);
    console.log('JWT userID:', decodedPayload.userId);

  }, []);

  useEffect(() => {
    if (!userID) return; // userIDが設定されるまで待つ

    const token = localStorage.getItem('jwtToken');
    if (!token) {
      console.error('JWT token not found');
      return;
    }

    // WebSocket接続を開始
    ws.current = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

    ws.current.onopen = () => {
      console.log("WebSocket connection established");
    };

    ws.current.onmessage = (event) => {
      const message: CustomMessage = JSON.parse(event.data);
      if (message.user_id !== userID) {
        console.log("Received message useID:", message.user_id);
        console.log("userID:", userID);
        setMessages((prevMessages) => [...prevMessages, message]); // 受信したメッセージを追加
      }
    };

    return () => {
      ws.current?.close(); // クリーンアップ時にWebSocketを閉じる
    };
  }, [userID]);

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
      id: String(messages.length + 1),
      user_id: "",
      workspace_id: workspaceID,
      target_id: channelID,
      text: input,
      role: "user",
      timestamp: new Date().toISOString(),
      action: "CREATE_MESSAGE",
    };

    // メッセージを追加し、入力内容をクリア
    setMessages((prevMessages) => [...prevMessages, newMessage]);
    setInput('');
    setIsLoading(true);

    // WebSocketでメッセージを送信
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(newMessage));
    }

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
