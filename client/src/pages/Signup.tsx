import React, { useState, ChangeEvent, FormEvent } from 'react';
import axios from 'axios';
import { Navigate } from 'react-router-dom';

const SignUp: React.FC = () => {
  // useState フックで状態を定義
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [message, setMessage] = useState<string>('');
  const [isInvalid, setIsInvalid] = useState<boolean>(false);
  const [redirect, setRedirect] = useState<boolean>(false);
  const [redirectTo, setRedirectTo] = useState<string>('/chat');
  const endpoint = 'http://localhost:8080/api/user/signup';

  // 入力変更時に状態を更新
  const onChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    if (name === 'email') {
      setEmail(value);
    } else if (name === 'password') {
      setPassword(value);
    }
  };

  // フォーム送信時の処理
  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const res = await axios.post(endpoint, {
        email,
        password,
      });

      if (res.status === 200) {
        console.log('Success: signup', res);
        localStorage.setItem('jwtToken', res.headers.authorization);

        setRedirectTo(redirectTo);
        setRedirect(true);
      } else {
        console.log('Failed: signup', res);
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
            Email
          </label>
          <input
            type="text"
            placeholder="Email"
            name="email"
            value={email}
            onChange={onChange}
            style={{
              width: '100%',
              padding: '10px',
              border: '1px solid #cbd5e0',
              borderRadius: '5px',
            }}
          />
          {isInvalid && <span style={{ color: 'red' }}>{message}</span>}
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
          <small style={{ color: '#718096' }}>Use a dummy password</small>
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
          SignUp
        </button>
      </form>
    </div>
  );
};

export default SignUp;
