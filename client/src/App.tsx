import React from 'react';
import './globals.css';
import { Providers } from './components/Providers';
import Home from './pages/Home';

const App: React.FC = () => {
  return (
    <div className="min-h-screen antialiased">
      <Providers>
        <main className="h-screen dark text-foreground bg-background">
          <Home />
        </main>
      </Providers>
    </div>
  );
};

export default App;
