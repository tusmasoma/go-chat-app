import './output.css';
import React from 'react';
import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import Chat from './pages/Chat'; // Chat.tsx のインポート

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        {/* URLが"/chat/:url"のときにChatコンポーネントを表示 */}
        <Route path="/chat/:url" element={<ChatWrapperWithParams />} />
      </Routes>
    </Router>
  );
};

// useParamsを使ってURLパラメータを取得し、Chatコンポーネントに渡す
const ChatWrapperWithParams: React.FC = () => {
  const params = useParams<{ url: string }>(); // React RouterでURLパラメータを取得
  return <Chat params={{ url: params.url }} />;
};

export default App;
