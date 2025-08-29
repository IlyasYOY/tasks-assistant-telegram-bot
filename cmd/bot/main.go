package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/config"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/handler/msgsender"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/store"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	bot, updates := newTelegramClient(cfg)
	aiClient := newOpenAIClient(cfg)

	db, sqlStore := newSQLStore(cfg)
	defer func() {
		if cerr := sqlStore.Close(); cerr != nil {
			log.Printf("error closing store: %v", cerr)
		}
		if cerr := db.Close(); cerr != nil {
			log.Printf("error closing DB: %v", cerr)
		}
	}()

	sender := msgsender.NewMessageSender(bot)

	h := handler.New(
		cfg,
		sender,
		handler.NewStartHandler(sender),
		handler.NewHelpHandler(sender),
		handler.NewNewTaskHandler(sender, sqlStore, cfg, &aiClient),
		handler.NewUnknownHandler(sender),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go runUpdatesHandling(updates, h)

	<-quit
	log.Println("Received shutdown signal, exiting...")
	bot.StopReceivingUpdates()
}

func runUpdatesHandling(updates tgbotapi.UpdatesChannel, h *handler.Handler) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if err := h.HandleUpdate(&update); err != nil {
			log.Printf(
				"error handling update from user %d: %v",
				update.Message.From.ID,
				err,
			)
		}
	}
}

func newSQLStore(cfg *config.Config) (*sql.DB, *store.SQLStore) {
	db, err := sql.Open("sqlite", cfg.SQLDSN)
	if err != nil {
		log.Fatalf("failed to open SQLite DB: %v", err)
	}

	// Run database migrations before creating the store.
	if migrateErr := store.Migrate(db, "./internal/store/migrations"); migrateErr != nil {
		log.Fatalf("database migration failed: %v", migrateErr)
	}
	log.Println("Database migrations applied successfully")

	sqlStore, err := store.NewSQLStore(db)
	if err != nil {
		log.Fatalf("failed to initialise SQL store: %v", err)
	}
	return db, sqlStore
}

func newOpenAIClient(cfg *config.Config) openai.Client {
	clientOpts := []option.RequestOption{
		option.WithAPIKey(cfg.OpenAPIKey),
	}
	if cfg.OpenAPIBasePath != "" {
		clientOpts = append(clientOpts, option.WithBaseURL(cfg.OpenAPIBasePath))
	}
	aiClient := openai.NewClient(clientOpts...)
	return aiClient
}

func newTelegramClient(
	cfg *config.Config,
) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}
	bot.Debug = false

	log.Printf("Bot authorized as %s", bot.Self.UserName)

	// Set up update handling (long polling)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := bot.GetUpdatesChan(u)
	return bot, updates
}
