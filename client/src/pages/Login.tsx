import React, { Component, ChangeEvent, FormEvent } from 'react';
import axios from 'axios';
import { Navigate } from 'react-router-dom';

interface LoginState {
  username: string;
  password: string;
  message: string;
  isInvalid: boolean;
  endpoint: string;
  redirect: boolean;
  redirectTo: string;
}

class Login extends Component<{}, LoginState> {
  constructor(props: {}) {
    super(props);
    this.state = {
      username: '',
      password: '',
      message: '',
      isInvalid: false,
      endpoint: 'http://localhost:8080/login',
      redirect: false,
      redirectTo: '/chat?u=',
    };
  }

  // 入力変更時に状態を更新
  onChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    this.setState({
      ...this.state,
      [name]: value,
    } as Pick<LoginState, keyof LoginState>);
  };

  // フォーム送信時の処理
  onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const res = await axios.post(this.state.endpoint, {
        username: this.state.username,
        password: this.state.password,
      });

      console.log('login', res);
      if (res.data.status) {
        const redirectTo = this.state.redirectTo + this.state.username;
        this.setState({ redirect: true, redirectTo });
      } else {
        this.setState({ message: res.data.message, isInvalid: true });
      }
    } catch (error) {
      console.log(error);
      this.setState({ message: 'something went wrong', isInvalid: true });
    }
  };

  render() {
    return (
      <div>
        {this.state.redirect && (
          <Navigate to={this.state.redirectTo} replace={true} />
        )}

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
          <form onSubmit={this.onSubmit}>
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
                value={this.state.username}
                onChange={this.onChange}
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
                value={this.state.password}
                onChange={this.onChange}
                style={{
                  width: '100%',
                  padding: '10px',
                  border: '1px solid #cbd5e0',
                  borderRadius: '5px',
                }}
              />
              {!this.state.isInvalid ? (
                ''
              ) : (
                <span style={{ color: 'red' }}>
                  Invalid username or password
                </span>
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
            <i style={{ fontSize: '16px', color: 'red' }}>{this.state.message}</i>
          </div>
        </div>
      </div>
    );
  }
}

export default Login;
