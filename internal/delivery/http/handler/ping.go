package http_handler

import (
	"net/http"
	"time"

	http_response "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/response"
	"github.com/gin-gonic/gin"
)

type pingResponse struct {
	Uptime   string `json:"uptime"`
	DateTime string `json:"datetime"`
}

func (h *HTTPHandler) Ping(c *gin.Context) {
	resp := &pingResponse{
		Uptime:   time.Since(h.startUp).String(),
		DateTime: time.Now().Format(time.RFC1123),
	}

	http_response.New(c, http.StatusOK, true, resp, nil)
}
