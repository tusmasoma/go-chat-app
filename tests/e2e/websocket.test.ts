import WebSocket from "ws";
import axios from "axios";

jest.setTimeout(1000000); // タイムアウトを延長

describe("WebSocket E2E Tests with Go Server", () => {
  let ws: WebSocket; // WebSocketクライアントのインスタンスを保持する変数
  let authToken: string;
  let clientID: string;
  let channelID: string;
  let msgID: string;
  let membershipIDofMsg: string;

  // 全てのテストの前に実行されるセットアップ処理
  beforeAll(async () => {
    try {
      // 認証リクエストを送信してトークンを取得
      const uniqueEmail = `e2e_test_${Date.now()}@test.com`;
      const response = await axios.post(
        "http://localhost:8083/api/user/signup",
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

      // テスト用のworkspaceを作成
      const response2 = await axios.post(
        "http://localhost:8083/api/workspace/create",
        {
          name: "TestWorkspace"
        },
        { headers }
      );
      const workspaceID = response2.data.workspace_id;

      // TODO: Membershipsの作成をする必要がある
      const response3 = await axios.post(
        `http://localhost:8083/api/membership/create/${workspaceID}`,
        {
          name: "Test User",
          profile_image_url: "http://example.com/image.jpg",
          is_admin: false
        },
        { headers }
      );

      ws = new WebSocket(`ws://localhost:8083/ws/${workspaceID}`, { headers });

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

  // 公開チャンネルの作成をテスト
  test("TEST: Create Public Channel", (done) => {
    const testChannel = `test_${Date.now()}_Channel`;
    const createChannelMessage = {
      action_tag: "CREATE_PUBLIC_CHANNEL",
      target_id: "",
      sender_id: clientID,
      content: {
        id: "",
        membership_id: "",
        text: testChannel,
        created: "2024-06-11T15:48:00Z",
        updated: null,
      },
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action_tag === "CREATE_PUBLIC_CHANNEL") {
        expect(receivedMessage.content.text).toBe(testChannel);
        channelID = receivedMessage.target_id;
        console.log("SUCCESS: CREATE_PUBLIC_CHANNEL");
        done();
      }
    });

    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(createChannelMessage));
    } else {
      ws.once("open", () => {
        ws.send(JSON.stringify(createChannelMessage));
      });
    }

    ws.once("error", (error) => {
      console.error("FAIL: CREATE_PUBLIC_CHANNEL", error);
      done(error);
    });
  });

  // メッセージの送信と受信をテスト
  test("TEST: Create Message", (done) => {
    const testMessage = {
      action_tag: "CREATE_MESSAGE",
      target_id: channelID,
      sender_id: clientID,
      content: {
        id: "",
        membership_id: "",
        text: "送信したいメッセージの内容",
        created_at: "2024-06-11T15:48:00Z",
        updated_at: null,
      },
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action_tag === "CREATE_MESSAGE") {
        expect(receivedMessage.action_tag).toBe(testMessage.action_tag);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.content.id).not.toBe("");
        expect(receivedMessage.content.membership_id).not.toBe("");
        expect(receivedMessage.content.text).toBe(testMessage.content.text);
        expect(receivedMessage.content.created_at).toBe(
          testMessage.content.created_at
        );
        console.log("SUCCESS: CREATE_MESSAGE");
        msgID = receivedMessage.content.id;
        membershipIDofMsg = receivedMessage.content.membership_id;
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

  // 初回メッセージの取得をテスト
  test("TEST: List Messages", (done) => {
    const listMessagesRequest = {
      action_tag: "LIST_MESSAGES",
      target_id: channelID,
      sender_id: clientID,
      content: {
        id: "",
        membership_id: "",
        text: "",
        created_at: "2024-06-11T15:48:00Z",
        updated_at: null,
      },
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action_tag === "LIST_MESSAGES") {
        expect(receivedMessage.action_tag).toBe(listMessagesRequest.action_tag);
        expect(receivedMessage.target_id).toBe(channelID);
        //expect(receivedMessage.sender_id).toBe(clientID);
        expect(Array.isArray(receivedMessage.contents)).toBe(true);
        expect(receivedMessage.contents.length).toBeGreaterThan(0);
        expect(receivedMessage.contents[0]).toHaveProperty("id");
        expect(receivedMessage.contents[0]).toHaveProperty("membership_id");
        expect(receivedMessage.contents[0]).toHaveProperty("text");
        expect(receivedMessage.contents[0].text).toBe(
          "送信したいメッセージの内容"
        );
        expect(receivedMessage.contents[0]).toHaveProperty("created_at");
        console.log("SUCCESS: LIST_MESSAGES");
        done();
      }
    });

    if (ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(listMessagesRequest));
    } else {
      ws.once("open", () => {
        ws.send(JSON.stringify(listMessagesRequest));
      });
    }

    ws.once("error", (error) => {
      console.error("FAIL: LIST_MESSAGES", error);
      done(error);
    });
  });

  // メッセージの更新をテスト
  test("TEST: Update Message", (done) => {
    const testMessage = {
      action_tag: "UPDATE_MESSAGE",
      target_id: channelID,
      sender_id: clientID,
      content: {
        id: msgID,
        membership_id: membershipIDofMsg,
        text: "更新したいメッセージの内容",
        created_at: "2024-06-11T15:48:00Z",
        updated_at: "2024-07-12T15:48:00Z",
      },
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action_tag === "UPDATE_MESSAGE") {
        expect(receivedMessage.action_tag).toBe(testMessage.action_tag);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.content.id).toBe(testMessage.content.id);
        expect(receivedMessage.content.membership_id).toBe(
          testMessage.content.membership_id
        );
        expect(receivedMessage.content.text).toBe(testMessage.content.text);
        expect(receivedMessage.content.created_at).toBe(
          testMessage.content.created_at
        );
        expect(receivedMessage.content.updated_at).toBe(
          testMessage.content.updated_at
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
      action_tag: "DELETE_MESSAGE",
      target_id: channelID,
      sender_id: clientID,
      content: {
        id: msgID,
        membership_id: membershipIDofMsg,
        text: "",
        created_at: "2024-06-11T15:48:00Z",
        updated_at: null,
      },
    };

    ws.once("message", (data) => {
      const receivedMessage = JSON.parse(data.toString());
      if (receivedMessage.action_tag === "DELETE_MESSAGE") {
        expect(receivedMessage.action_tag).toBe(testMessage.action_tag);
        expect(receivedMessage.target_id).toBe(testMessage.target_id);
        //expect(receivedMessage.sender_id).toBe(testMessage.sender_id);
        expect(receivedMessage.content.membership_id).toBe(
          testMessage.content.membership_id
        );
        expect(receivedMessage.content.text).toBe(testMessage.content.text);
        expect(receivedMessage.content.created_at).toBe(
          testMessage.content.created_at
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