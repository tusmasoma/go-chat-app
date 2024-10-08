CREATE DATABASE IF NOT EXISTS `go_chat_app_test_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `go_chat_app_test_db`;

DROP TABLE IF EXISTS Messages CASCADE;
DROP TABLE IF EXISTS Membership_Channels CASCADE;
DROP TABLE IF EXISTS Memberships CASCADE;
DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS Channels CASCADE;
DROP TABLE IF EXISTS Workspaces CASCADE;

CREATE TABLE Workspaces (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    name VARCHAR(50) NOT NULL
);

CREATE TABLE Channels (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    workspace_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    private BOOLEAN NOT NULL,
    UNIQUE (workspace_id, name)
);

CREATE TABLE Users (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL  -- 暗号化されたパスワードを格納
);

CREATE TABLE Memberships (
    user_id CHAR(36) NOT NULL,
    workspace_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    profile_image_url VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (user_id, workspace_id)
);

CREATE TABLE Membership_Channels (
    user_id CHAR(36) NOT NULL,
    workspace_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    PRIMARY KEY (user_id, workspace_id, channel_id),
    FOREIGN KEY (user_id, workspace_id) REFERENCES Memberships(user_id, workspace_id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES Channels(id) ON DELETE CASCADE
);

CREATE TABLE Messages (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    user_id CHAR(36) NOT NULL,
    workspace_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);