package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/v1gn35h7/goshell/server/goshell"
	"github.com/v1gn35h7/goshell/server/pb"
)

// Shell serviec contracts
type GetScriptRequest struct {
	AgentId string `json:"agentid"`
}

type GetScriptResponse struct {
	Scripts []*goshell.ShellScript `json:"scripts"`
}

// Grpc server
type grpcServer struct {
	getScripts grpctransport.Handler
	pb.UnimplementedShellServiceServer
}

func NewGRPCServer(ep endpoint.Endpoint) *grpcServer {

	return &grpcServer{
		getScripts: grpctransport.NewServer(
			ep,
			decodeGetScriptsRequest,
			encodeGetScriptsResponse,
		),
	}

}

func (s *grpcServer) GetScripts(ctx context.Context, req *pb.ShellRequest) (*pb.ShellResponse, error) {
	_, rep, err := s.getScripts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShellResponse), nil
}

// decodeGetScriptsRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGetScriptsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShellRequest)
	return GetScriptRequest{AgentId: string(req.AgentId)}, nil
}

// encodeGRPCSumResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain sum response to a gRPC sum reply. Primarily useful in a server.
func encodeGetScriptsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetScriptResponse)
	scripts := make([]*pb.ShellScript, 0)
	for _, v := range resp.Scripts {
		scripts = append(scripts, &pb.ShellScript{
			Script: v.Script,
			Args:   v.Args,
			Type:   v.Type,
		})
	}
	return &pb.ShellResponse{Scripts: scripts}, nil
}