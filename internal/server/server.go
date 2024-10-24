package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ctxTimeout = 5
	certFile   = "./ssl/server.crt"
	keyFile    = "./ssl/server.key"
)

type Server struct {
	gin        *gin.Engine
	logger     *logrus.Logger
	cfg        *config.Config
	postgresDb *gorm.DB
}

// NewServer server constructor
func NewServer(
	log *logrus.Logger,
	config *config.Config,
	postgresDb *gorm.DB,
) *Server {
	return &Server{
		gin:        gin.New(),
		logger:     log,
		cfg:        config,
		postgresDb: postgresDb,
	}
}

// running server with SLL/TLS or not
func (s *Server) Run() error {
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port),
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	if s.cfg.Server.SSL {
		serverError := make(chan error)
		quit := make(chan os.Signal)

		// listen signal interrupt/terminate from os
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		go func() {
			s.logger.Infof("TLS server listening on %s", server.Addr)
			serverError <- server.ListenAndServeTLS(certFile, keyFile)
		}()

		select {
		case err := <-serverError:
			{
				s.logger.Fatalf("Error starting TLS server: %v", err)
			}
		case <-quit:
			{
				ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
				defer shutdown()
				if err := server.Shutdown(ctx); err != nil {
					s.logger.Fatalf("Error gracefully shutting down server: %v", err)
				}
			}
		}
		s.logger.Info("Server Exited Properly")
		return nil
	}

	serverError := make(chan error, 1)
	quit := make(chan os.Signal, 1)

	// listen signal interrupt/terminate from OS
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.logger.Infof("Server listening on %s", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		{
			s.logger.Fatalf("Srror listening and serving: %v", err)
		}
	case <-quit:
		{
			ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
			defer shutdown()
			if err := server.Shutdown(ctx); err != nil {
				s.logger.Fatalf("Error gracefully shutting down server: %v", err)
			}
		}
	}
	s.logger.Info("Server Exited Properly")
	return nil
}
