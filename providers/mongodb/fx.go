package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

func FXMongoDBConnection(cmd *cobra.Command, ctx context.Context, lc fx.Lifecycle) (*mongo.Database, error) {
	var err error
	host := viper.GetString(hostFlag)
	if host == "" {
		host, err = cmd.Flags().GetString(hostFlag)
		if err != nil {
			return nil, ErrCredentialsNotFound
		}
	}
	port := viper.GetString(portFlag)
	if port == "" {
		port, err = cmd.Flags().GetString(portFlag)
		if port == "" || err != nil {
			port = "27017"
		}
	}
	password := viper.GetString(passwordFlag)
	if password == "" {
		password, err = cmd.Flags().GetString(passwordFlag)
		if err != nil {
			return nil, ErrCredentialsNotFound
		}
	}
	username := viper.GetString(userFlag)
	if username == "" {
		username, err = cmd.Flags().GetString(userFlag)
		if err != nil {
			return nil, ErrCredentialsNotFound
		}
	}
	db := viper.GetString(dbFlag)
	if db == "" {
		db, err = cmd.Flags().GetString(dbFlag)
		if err != nil {
			return nil, ErrCredentialsNotFound
		}
	}
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
	clientOptions = clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx, nil)
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return client.Database(db), err
}
