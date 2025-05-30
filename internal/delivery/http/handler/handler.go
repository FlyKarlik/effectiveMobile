package http_handler

import (
	"time"

	"github.com/FlyKarlik/effectiveMobile/internal/usecase"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type HTTPHandler struct {
	startUp time.Time
	logger  logger.Logger
	usecase *usecase.Usecase
}

func New(logger logger.Logger, usecase *usecase.Usecase) *HTTPHandler {
	return &HTTPHandler{
		startUp: time.Now(),
		logger:  logger,
		usecase: usecase,
	}
}
