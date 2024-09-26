package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.uber.org/dig"

	"github.com/tusmasoma/go-chat-app/config"
	"github.com/tusmasoma/go-chat-app/entity"
	"github.com/tusmasoma/go-chat-app/interfaces/handler"
	"github.com/tusmasoma/go-chat-app/interfaces/websocket"
	"github.com/tusmasoma/go-chat-app/repository/mysql"
	"github.com/tusmasoma/go-chat-app/usecase"
)

func BuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		log.Error("Failed to provide context")
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		config.NewCacheConfig,
		config.NewDBConfig,
		mysql.NewMySQLDB,
		mysql.NewTransactionRepository,
		mysql.NewMessageRepository,
		usecase.NewMessageUseCase,
		generateHubManager,
		handler.NewWebsocketHandler,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
		) *chi.Mux {
			r := chi.NewRouter()
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:     []string{"https://*", "http://*"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposedHeaders:     []string{"Link", "Authorization"},
				AllowCredentials:   true,
				MaxAge:             serverConfig.PreflightCacheDurationSec,
				OptionsPassthrough: true,
			}))

			r.Group(func(r chi.Router) {
				r.Get("/ws/{workspace_id}", wsHandler.WebSocket)
			})

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			log.Critical("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
			return nil, err
		}
	}

	log.Info("Container built successfully")
	return container, nil
}

func generateHubManager(ctx context.Context) *websocket.HubManager {
	//  現状、Workspaceは一つの為、containerにてHubManagerを生成して、DIする
	workspaceID := os.Getenv("WORKSPACE_ID")
	if workspaceID == "" {
		log.Critical("Failed to get workspace ID")
		return nil
	}
	hub, err := entity.NewHub(workspaceID, "DefaultWorkspace")
	if err != nil {
		log.Critical("Failed to create new hub", log.Ferror(err))
		return nil
	}
	hm := websocket.NewHubManager(hub)

	go hm.Run(ctx)

	log.Info("HubManager created successfully")

	return &hm
}
