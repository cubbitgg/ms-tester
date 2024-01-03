package endpoints

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

const (
	HealthPath = "/health"
	LoggerName = "efesto"
)

type Server struct {
	engine *gin.Engine
	host   string
	port   int
}

func NewServer() *Server {
	engine := gin.New()

	engine.Use(gintrace.Middleware(LoggerName))

	engine.SetTrustedProxies(nil)

	server := &Server{
		engine: engine,
		host:   "",
		port:   8080,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {

	unauthenticatedRoute := s.engine.Group("/")
	unauthenticatedRoute.GET(HealthPath, s.createHealthRoute())
	unauthenticatedRoute.GET("/500", s.create500Route())
	unauthenticatedRoute.GET("/502", s.create502Route())
	unauthenticatedRoute.GET("/503", s.create503Route())
	unauthenticatedRoute.GET("/504", s.create504Route())
	unauthenticatedRoute.GET("/timeout", s.createTimeoutRoute())
}

func (s *Server) createHealthRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	}
}

func (s *Server) create500Route() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "500",
		})
	}
}

func (s *Server) create502Route() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status": "502",
		})
	}
}

func (s *Server) create503Route() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "503",
		})
	}
}

func (s *Server) create504Route() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"status": "504",
		})
	}
}

func (s *Server) createTimeoutRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmeoutString := ctx.Query("timeout")
		if len(tmeoutString) <= 0 {
			tmeoutString = "60"
		}

		timeout, err := strconv.Atoi(tmeoutString)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "500",
				"msg":    err,
			})
			return
		}
		time.Sleep(time.Duration(timeout) * time.Second)

		ctx.JSON(http.StatusOK, gin.H{
			"status": "200",
			"msg":    fmt.Sprintf("wait for %d sec", timeout),
		})
	}
}

func (s *Server) Listen() {
	address := fmt.Sprintf("%s:%d", s.host, s.port)
	log.Info().Msgf("Listening on %s", address)
	s.engine.Run(address)
}
