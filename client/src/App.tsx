import React from 'react';
import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import Home from './pages/Home';
import Landing from './pages/Landing';
import SignUp from './pages/Signup';
import Login from './pages/Login';

const App: React.FC = () => {
  return (
      <Router>
        <Routes>
          <Route path="/" element={<Landing />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/login" element={<Login />} />
          <Route path="/chat/" element={<Home />} />
        </Routes>
      </Router>
  );
};

export default App;