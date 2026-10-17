package analyst

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/tousart/factory_analyst/pkg/api/analyst"
)

type Transport struct {
	pb.UnimplementedAnalystServer

	analystSrv AnalystService
}

func NewTransport(
	analystSrv AnalystService,
) *Transport {
	return &Transport{
		analystSrv: analystSrv,
	}
}

// Register регистрирует ручки аналитики на переданном gRPC-сервере.
func (t *Transport) Register(s grpc.ServiceRegistrar) {
	pb.RegisterAnalystServer(s, t)
}

// RegisterGateway регистрирует HTTP-ручки аналитики на grpc-gateway.
func (t *Transport) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return pb.RegisterAnalystHandler(ctx, mux, conn)
}
