package app

import (
	"github.com/hashicorp/go-hclog"
	"os"
	"os/signal"
	"syscall"
)

func (s *server) gracefulShutdown(logger hclog.Logger) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		logger.Warn("Received shutdown signal")

		err := s.app.Shutdown()
		if err != nil {
			logger.Error("Error shutting down app", "error", err)
		}

		//err = s.db.Close()
		//if err != nil {
		//	logger.Error("Failed to close database", "error", err)
		//}

		logger.Info("Shutdown complete.")
	}()
}
