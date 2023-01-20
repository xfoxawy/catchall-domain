package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/xfoxawy/catchall-domain/app/interface/http"
	"github.com/xfoxawy/catchall-domain/app/repository"
	db "github.com/xfoxawy/catchall-domain/providers/mongodb"
)

var serveHttp = &cobra.Command{
	Use:   "serve-http",
	Short: "start catchall service http API server",
	Run:   runHttpCmd,
}

func init() {
	rootCmd.AddCommand(serveHttp)
}

func runHttpCmd(cmd *cobra.Command, _ []string) {
	startHttpApi(cmd).Run()
}

func startHttpApi(cmd *cobra.Command) *fx.App {
	return fx.New(
		fx.Provide(
			func() *cobra.Command { return cmd },
			func() context.Context { return context.Background() },
			db.FXMongoDBConnection,
			repository.NewEventsCounterRepository,
			http.NewHttpHandlers,
			InitRouter,
		),
		fx.Invoke(registerEventApi, serveHTTP),
	)
}

func InitRouter() *gin.Engine {
	router := gin.New()
	return router
}

func registerEventApi(args http.EventsApiArgs) error {
	return http.RegisterEventsApi(args)
}

func serveHTTP(router *gin.Engine) error {
	port := viper.GetString("SERVER_PORT")
	if port == "" {
		return errors.New("failed to read server port")
	}
	fmt.Printf("The HTTP API is now available at http://localhost:%s", port)
	gin.SetMode(gin.DebugMode)
	return router.Run(port)
}
