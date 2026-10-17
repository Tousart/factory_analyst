package analyst

import (
	"context"
	"fmt"

	pb "github.com/tousart/factory_analyst/pkg/api/analyst"
)

func (t *Transport) GetUsersByIDs(_ context.Context, req *pb.GetUsersByIDsRequest) (*pb.GetUsersByIDsResponse, error) {
	return &pb.GetUsersByIDsResponse{
		Message: fmt.Sprintf("Пользователь: %s", req.GetId()),
	}, nil
}
