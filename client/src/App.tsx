import React from 'react';
import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import Home from './pages/Home';

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/chat/" element={<Home />} />
      </Routes>
    </Router>
  );
};

export default App;