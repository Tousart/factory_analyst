package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const shutdownTimeout = 10 * time.Second

// Service регистрирует свои ручки на gRPC-сервере и grpc-gateway.
type Service interface {
	Register(s grpc.ServiceRegistrar)
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type Server struct {
	Addr        string
	GatewayAddr string

	grpcServer *grpc.Server
	services   []Service
}

func NewServer(addr, gatewayAddr string, services ...Service) *Server {
	grpcServer := grpc.NewServer()

	for _, service := range services {
		service.Register(grpcServer)
	}

	return &Server{
		Addr:        addr,
		GatewayAddr: gatewayAddr,
		grpcServer:  grpcServer,
		services:    services,
	}
}

// Run поднимает gRPC-сервер и grpc-gateway и блокируется до SIGINT/SIGTERM,
// после чего даёт активным запросам shutdownTimeout на завершение.
func (s *Server) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	gatewayServer, conn, err := s.newGatewayServer(ctx)
	if err != nil {
		lis.Close()
		return err
	}
	defer conn.Close()

	errCh := make(chan error, 2)

	go func() {
		log.Println("gRPC server is listening on", s.Addr)
		if err := s.grpcServer.Serve(lis); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()

	go func() {
		log.Println("gRPC-gateway is listening on", s.GatewayAddr)
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("grpc-gateway: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		gatewayServer.Close()
		s.grpcServer.Stop()
		return err
	case <-ctx.Done():
		log.Println("Shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Сначала гасим gateway: он проксирует запросы в gRPC-сервер.
	if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
		log.Println("gRPC-gateway shutdown:", err)
	}

	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Graceful shutdown")
	case <-shutdownCtx.Done():
		s.grpcServer.Stop()
	}

	log.Println("Server stopped")

	return nil
}

// newGatewayServer создаёт HTTP-сервер grpc-gateway, который ходит в gRPC-сервер по его адресу.
func (s *Server) newGatewayServer(ctx context.Context) (*http.Server, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(s.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("grpc-gateway client: %w", err)
	}

	mux := runtime.NewServeMux()

	for _, service := range s.services {
		if err := service.RegisterGateway(ctx, mux, conn); err != nil {
			conn.Close()
			return nil, nil, fmt.Errorf("register grpc-gateway: %w", err)
		}
	}

	return &http.Server{
		Addr:    s.GatewayAddr,
		Handler: mux,
	}, conn, nil
}
