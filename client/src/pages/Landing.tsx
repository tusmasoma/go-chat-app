import React from 'react';
import { Link } from 'react-router-dom';

const Landing: React.FC = () => {
  return (
    <div
      style={{
        maxWidth: '600px',
        marginTop: '3rem',
        marginLeft: 'auto',  // 左右中央揃え
        marginRight: 'auto', // 左右中央揃え
        textAlign: 'center',
        display: 'flex',       // フレックスボックスで垂直中央揃え
        flexDirection: 'column',
        justifyContent: 'center',
        minHeight: 'calc(100vh - 18rem)',
      }}
    >
      <div style={{ padding: '20px', marginBottom: '20px' }}>
        <h1 style={{ fontSize: '2rem', marginBottom: '20px' }}>
          This is a simple chat application written in ReactJS.
        </h1>
        <p>It is using websocket for communication.</p>
      </div>
      <div>
        <div style={{ display: 'flex', justifyContent: 'center', gap: '20px' }}>
          <Link to="signup">
            <button
              style={{
                padding: '10px 20px',
                fontSize: '1.2rem',
                backgroundColor: 'green',
                color: 'white',
                border: 'none',
                borderRadius: '5px',
                cursor: 'pointer',
              }}
            >
              SignUp
            </button>
          </Link>
          <Link to="login">
            <button
              style={{
                padding: '10px 20px',
                fontSize: '1.2rem',
                backgroundColor: 'transparent',
                color: 'green',
                border: '2px solid green',
                borderRadius: '5px',
                cursor: 'pointer',
              }}
            >
              Login
            </button>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Landing;
