package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type EventsApiArgs struct {
	fx.In

	Router   *gin.Engine
	Handlers *Handlers
}

func RegisterEventsApi(args EventsApiArgs) error {
	events := args.Router.Group("/events/:domain")
	{
		events.PUT("/delivered", args.Handlers.EventProcessorHandler)
		events.PUT("/bounced", args.Handlers.EventProcessorHandler)
	}
	args.Router.GET("/domains/:domain", args.Handlers.GetDomainStatusHandler)
	return nil
}
