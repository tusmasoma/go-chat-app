USE `go_chat_app_db`;

-- Workspaceデフォルトデータの挿入
INSERT INTO Workspaces (id, name) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'DefaultWorkspace');

-- Channelsデフォルトデータの挿入
INSERT INTO Channels (id, workspace_id, name, private) VALUES
('123e4567-e89b-12d3-a456-426614174000', '550e8400-e29b-41d4-a716-446655440000', 'DefaultChannel', FALSE);