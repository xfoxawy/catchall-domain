package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mailgun/catchall"

	"github.com/xfoxawy/catchall-domain/app/repository"
)

type Handlers struct {
	repo *repository.EventsCounterRepository
}

func NewHttpHandlers(repo *repository.EventsCounterRepository) *Handlers {
	return &Handlers{repo: repo}
}

func (h *Handlers) GetDomainStatusHandler(ctx *gin.Context) {
	var err error
	domain := ctx.Param("domain")
	decounter, err := h.repo.FindOneByDomain(ctx, domain)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	if decounter.Bounced > 0 {
		ctx.JSON(200, gin.H{
			"status": "not catch-all",
		})
		return
	}
	if decounter.Delivered > 1000 {
		ctx.JSON(200, gin.H{
			"status": "catch-all",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": "unknown",
	})
}

func (h *Handlers) EventProcessorHandler(ctx *gin.Context) {
	var err error
	var i catchall.Event
	err = ctx.ShouldBindJSON(&i)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	_, err = h.repo.IncrementCounter(ctx, i.Domain, i.Type)
	if err != nil {
		ctx.AbortWithError(500, err)
	}
	ctx.JSON(202, i)
}
