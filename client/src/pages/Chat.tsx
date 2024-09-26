import React, { useEffect, useState } from 'react';
import { ChatWrapper } from '../components/ChatWrapper'; // 適切なパスに変更
import Cookies from 'js-cookie'; // クライアントサイドでのクッキー操作
import { fetchInitialMessages, fetchIsIndexed, indexUrl } from '../apis/api'; // API呼び出し関数を実装

interface ChatProps {
  params: {
    url: string | string[] | undefined;
  };
}

function reconstructUrl(url: string[]): string {
  const decodedComponents = url.map((component) => decodeURIComponent(component));
  return decodedComponents.join("//");
}

const Chat: React.FC<ChatProps> = ({ params }) => {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [initialMessages, setInitialMessages] = useState<any[]>([]);

//   useEffect(() => {
//     const fetchData = async () => {
//       const sessionCookie = Cookies.get("sessionId"); // クッキー取得
//       const reconstructedUrl = reconstructUrl(params.url as string[]);

//       const sessionIdValue = (reconstructedUrl + "--" + sessionCookie).replace(/\//g, "");
//       setSessionId(sessionIdValue);

//       const isAlreadyIndexed = await fetchIsIndexed(reconstructedUrl); // Redisチェック

//       const messages = await fetchInitialMessages({ sessionId: sessionIdValue });
//       setInitialMessages(messages);

//       if (!isAlreadyIndexed) {
//         await indexUrl(reconstructedUrl); // URLのインデックスをサーバーに指示
//       }
//     };

//     fetchData();
//   }, [params.url]);

//  if (!sessionId) return <div>Loading...</div>;

  return <ChatWrapper sessionId={""} initialMessages={initialMessages} />;
};

export default Chat;
