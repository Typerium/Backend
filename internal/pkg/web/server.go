package web

import (
	"sync"

	"github.com/gramework/gramework"
	"go.uber.org/zap"

	"typerium/internal/pkg/logging"
)

func NewServer(log *zap.Logger) *Server {
	gramework.DisableFlags()

	server := &Server{
		App: gramework.New(),
		log: log.Named("http_server"),
	}

	server.App.Logger = logging.NewGrameworkLogger(server.log)

	err := server.App.Use(server.App.CORSMiddleware())
	if err != nil {
		server.log.Fatal("can't create cors middleware", zap.Error(err))
	}
	server.App.NotFound(notFoundHandler)
	server.App.MethodNotAllowed(methodNotAllowed)

	return server
}

type Server struct {
	*gramework.App
	wg  sync.WaitGroup
	log *zap.Logger
}

func (s *Server) Start(addr string) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		err := s.App.ListenAndServe(addr)
		if err != nil {
			s.log.Fatal("can't start server", zap.Error(err))
		}
	}()
}

func (s *Server) Stop() {
	err := s.App.Shutdown()
	if err != nil {
		s.log.Error("failed stopping server", zap.Error(err))
	}
	s.wg.Wait()
}
