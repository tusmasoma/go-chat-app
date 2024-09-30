import React, { useState, ChangeEvent, FormEvent } from 'react';
import axios from 'axios';
import { Navigate } from 'react-router-dom';

const Login: React.FC = () => {
  // useState フックで状態を定義
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [message, setMessage] = useState<string>('');
  const [isInvalid, setIsInvalid] = useState<boolean>(false);
  const [redirect, setRedirect] = useState<boolean>(false);
  const [redirectTo, setRedirectTo] = useState<string>('/chat?u=');
  const endpoint = 'http://localhost:8080/api/user/login';

  // 入力変更時に状態を更新
  const onChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    if (name === 'username') {
      setUsername(value);
    } else if (name === 'password') {
      setPassword(value);
    }
  };

  // フォーム送信時の処理
  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const res = await axios.post(endpoint, {
        username,
        password,
      });

      console.log('login', res);
      if (res.data.status) {
        setRedirectTo(redirectTo + username);
        setRedirect(true);
      } else {
        setMessage(res.data.message);
        setIsInvalid(true);
      }
    } catch (error) {
      console.log(error);
      setMessage('something went wrong');
      setIsInvalid(true);
    }
  };

  // リダイレクトの処理
  if (redirect) {
    return <Navigate to={redirectTo} replace={true} />;
  }

  return (
    <div
      style={{
        marginTop: '40px',
        maxWidth: '600px',
        marginLeft: 'auto',
        marginRight: 'auto',
        textAlign: 'left',
        padding: '20px',
        border: '2px solid #e2e8f0',
        borderRadius: '8px',
      }}
    >
      <form onSubmit={onSubmit}>
        <div style={{ marginBottom: '20px' }}>
          <label
            style={{
              display: 'block',
              marginBottom: '5px',
              fontWeight: 'bold',
            }}
          >
            Username
          </label>
          <input
            type="text"
            placeholder="Username"
            name="username"
            value={username}
            onChange={onChange}
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #cbd5e0',
              borderRadius: '5px',
            }}
          />
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label
            style={{
              display: 'block',
              marginBottom: '5px',
              fontWeight: 'bold',
            }}
          >
            Password
          </label>
          <input
            type="password"
            placeholder="Password"
            name="password"
            value={password}
            onChange={onChange}
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #cbd5e0',
              borderRadius: '5px',
            }}
          />
          {!isInvalid ? (
            ''
          ) : (
            <span style={{ color: 'red' }}>Invalid username or password</span>
          )}
        </div>

        <button
          type="submit"
          style={{
            width: '100%',
            padding: '15px',
            backgroundColor: 'green',
            color: 'white',
            fontSize: '16px',
            border: 'none',
            borderRadius: '5px',
            cursor: 'pointer',
          }}
        >
          Login
        </button>
      </form>

      <div style={{ paddingTop: '10px' }}>
        <i style={{ fontSize: '16px', color: 'red' }}>{message}</i>
      </div>
    </div>
  );
};

export default Login;
