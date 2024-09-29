import WebSocket from "ws";
import axios from "axios";

jest.setTimeout(1000000); // タイムアウトを延長

describe("WebSocket E2E Tests with Go Server", () => {
  let ws: WebSocket; // WebSocketクライアントのインスタンスを保持する変数
  let authToken: string;
  let clientID: string;
  let workspaceID: string = "550e8400-e29b-41d4-a716-446655440000";
  let channelID: string = "123e4567-e89b-12d3-a456-426614174000";
  let msgID: string;
  let membershipIDofMsg: string;

  // 全てのテストの前に実行されるセットアップ処理
  beforeAll(async () => {
    try {
      // 認証リクエストを送信してトークンを取得
      const uniqueEmail = `e2e_test_${Date.now()}@test.com`;
      const response = await axios.post(
        "http://localhost:8080/api/user/signup",
        {
          email: uniqueEmail,
          password: "e2e_test_password",
        }
      );
      authToken = response.headers["authorization"].split(" ")[1];

      // 取得したトークンを使ってWebSocket接続を確立
      const headers = {
        Authorization: `Bearer ${authToken}`,
      };

      console.log("authToken: ", authToken);

      ws = new WebSocket(`ws://localhost:8080/ws`, { headers });

      await new Promise<void>((resolve, reject) => {
        ws.on("open", () => {
          console.log("WebSocket connection established");
          resolve(); // WebSocket接続が確立されたらセットアップ完了
        });

        ws.on("error", (error) => {
          console.error("WebSocket connection error:", error);
          reject(error); // 接続エラーが発生した場合はテスト失敗
        });
      });
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Setup failed: ${error.message}`); // 認証リクエストに失敗した場合はテスト失敗
      } else {
        throw new Error("Setup failed: An unknown error occurred"); // 不明なエラーの場合
      }
    }
  });

  // 全てのテストの後に実行されるクリーンアップ処理
  afterAll(() => {
    ws.close(); // WebSocket接続を閉じる
  });

  // WebSocket接続が正しく確立されるかをテスト
  test("TEST: WebSocket connection", (done) => {
    // WebSocket接続が確立された時に発生するイベント
    ws.on("open", () => {
      expect(ws.readyState).toBe(WebSocket.OPEN); // WebSocketがOPEN状態であることを確認
      console.log("SUCCESS: WebSocket connection");
      done(); // テスト完了
    });

    // すでにオープンしている場合の処理
    if (ws.readyState === WebSocket.OPEN) {
      console.log("WebSocket was already open");
      expect(ws.readyState).toBe(WebSocket.OPEN);
      done(); // テスト完了
    }

    ws.on("error", (error) => {
      console.error("FAIL: WebSocket connection", error);
      done(error); // エラー発生時はテスト失敗
    });
  });

  // メッセージの送信と受信をテスト
  test("TEST: Create Message", (done) => {
    const testMessage = {
      id: "",
      user_id: "",
      workspace_id: workspaceID,
      text: "送信したいメッセージの内容",
      created_at: "2024-06-11T15:48:00Z",
      action: "CREATE_MESSAGE",
      target_id: channelID,
      // sender_id: clientID,
      // membership_id: "",
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action === "CREATE_MESSAGE") {
        expect(receivedMessage.action).toBe(testMessage.action);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.id).not.toBe("");
        expect(receivedMessage.membership_id).not.toBe("");
        expect(receivedMessage.text).toBe(testMessage.text);
        expect(receivedMessage.created_at).toBe(
          testMessage.created_at
        );
        console.log("SUCCESS: CREATE_MESSAGE");
        msgID = receivedMessage.id;
        membershipIDofMsg = receivedMessage.membership_id;
        done();
      }
    });

    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(testMessage));
    } else {
      ws.once("open", () => {
        ws.send(JSON.stringify(testMessage));
      });
    }

    ws.once("error", (error) => {
      console.error("FAIL: CREATE_MESSAGE", error);
      done(error);
    });
  });

  // メッセージの更新をテスト
  test("TEST: Update Message", (done) => {
    const testMessage = {
      id: msgID,
      user_id: "",
      workspace_id: workspaceID,
      text: "更新したいメッセージの内容",
      created_at: "2024-06-11T15:48:00Z",
      action: "UPDATE_MESSAGE",
      target_id: channelID,
      // sender_id: clientID,
      membership_id: membershipIDofMsg,
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action === "UPDATE_MESSAGE") {
        expect(receivedMessage.action).toBe(testMessage.action);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.id).toBe(testMessage.id);
        expect(receivedMessage.membership_id).toBe(
          testMessage.membership_id
        );
        expect(receivedMessage.text).toBe(testMessage.text);
        expect(receivedMessage.created_at).toBe(
          testMessage.created_at
        );
        console.log("SUCCESS: UPDATE_MESSAGE");
        done();
      }
    });

    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(testMessage));
    } else {
      ws.once("open", () => {
        ws.send(JSON.stringify(testMessage));
      });
    }

    ws.once("error", (error) => {
      console.error("FAIL: UPDATE_MESSAGE", error);
      done(error);
    });
  });

  // メッセージの削除をテスト
  test("TEST: Delete Message", (done) => {
    const testMessage = {
      id: msgID,
      user_id: "",
      workspace_id: workspaceID,
      text: "",
      created_at: "2024-06-11T15:48:00Z",
      action: "DELETE_MESSAGE",
      target_id: channelID,
      // sender_id: clientID,
      membership_id: membershipIDofMsg,
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action === "DELETE_MESSAGE") {
        expect(receivedMessage.action).toBe(testMessage.action);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.membership_id).toBe(
          testMessage.membership_id
        );
        expect(receivedMessage.text).toBe(testMessage.text);
        expect(receivedMessage.created_at).toBe(
          testMessage.created_at
        );
        console.log("SUCCESS: DELETE_MESSAGE");
        done();
      }
    });

    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(testMessage));
    } else {
      ws.once("open", () => {
        ws.send(JSON.stringify(testMessage));
      });
    }

    ws.once("error", (error) => {
      console.error("FAIL: DELETE_MESSAGE", error);
      done(error);
    });
  });
});