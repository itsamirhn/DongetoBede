package cmd

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/itsamirhn/dongetobede/internal/bot"
	"github.com/itsamirhn/dongetobede/internal/config"
	"github.com/itsamirhn/dongetobede/internal/database"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Telegram Bot serving",
	Run:   serve,
}

func serve(cmd *cobra.Command, _ []string) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Panic(err)
		}
	}()
	loadConfigOrPanic(cmd)

	mc := createMongoClientOrPanic()

	defer func() {
		if err := mc.Disconnect(context.Background()); err != nil {
			logrus.WithError(err).Error("failed to disconnect from mongo")
		}
	}()

	db := database.NewMongoClient(mc, config.GlobalConfig.DB.Name)

	b := createBotOrPanic(db)

	logrus.Info("Starting the telegram bot server...")
	b.Start()
}

func createBotOrPanic(db database.Client) *bot.Bot {
	b, err := bot.NewBot(db, config.GlobalConfig.Token, config.GlobalConfig.Endpoint, config.GlobalConfig.ListenPort)
	if err != nil {
		logrus.WithError(err).Panic("failed to create bot")
	}
	return b
}

func createMongoClientOrPanic() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.GlobalConfig.DB.URI))
	if err != nil {
		logrus.WithError(err).Panic("failed to create mongo client")
	}

	mCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Ping(mCtx, readpref.Primary())
	if err != nil {
		logrus.WithError(err).Panic("failed to ping mongo")
	}
	logrus.Info("Successfully connected to mongodb")
	return client
}
