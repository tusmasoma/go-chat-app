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
	"github.com/tusmasoma/go-chat-app/repository"
	"github.com/tusmasoma/go-chat-app/repository/auth"
	"github.com/tusmasoma/go-chat-app/repository/mysql"
	"github.com/tusmasoma/go-chat-app/repository/redis"
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
		mysql.NewUserRepository,
		mysql.NewMembershipRepository,
		auth.NewAuthRepository,
		redis.NewRedisClient,
		redis.NewPubSubRepository,
		usecase.NewMessageUseCase,
		usecase.NewUserUseCase,
		generateHubManager,
		handler.NewWebsocketHandler,
		handler.NewUserHandler,
		func(
			serverConfig *config.ServerConfig,
			wsHandler *handler.WebsocketHandler,
			userHandler handler.UserHandler,
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

			r.Route("/api", func(r chi.Router) {
				r.Route("/user", func(r chi.Router) {
					r.Post("/create", userHandler.CreateUser)
					// r.Post("/login", userHandler.Login)
					// r.Group(func(r chi.Router) {
					// 	r.Use(authMiddleware.Authenticate)
					// 	r.Get("/logout", userHandler.Logout)
					// })
				})
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

func generateHubManager(ctx context.Context, psr repository.PubSubRepository) *websocket.HubManager {
	//  現状、Workspaceは一つの為、containerにてHubManagerを生成して、DIする
	//  同様に、ChannelManagerも生成してDIする
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

	go hm.Run()

	channelID := os.Getenv("CHANNEL_ID")
	if channelID == "" {
		log.Critical("Failed to get channel ID")
		return nil
	}
	channel, err := entity.NewChannel(channelID, "DefaultChannel", false)
	if err != nil {
		log.Critical("Failed to create new channel", log.Ferror(err))
		return nil
	}

	cm := websocket.NewChannelManager(channel, psr)
	go cm.Run(ctx)
	hm.RegisterChannelManager(cm)

	log.Info("HubManager created successfully")

	return &hm
}
