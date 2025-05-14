package main

import (
	"context"
	"github.com/caiomarcatti12/nanogo/pkg/log"
	"github.com/stackmesh/api-log-trail-proto/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type EventConsumerHandler struct {
	pb.UnimplementedEventConsumerServer
	marshaler protojson.MarshalOptions
	logger    log.ILog
}

func NewEventConsumerHandler(logger log.ILog) *EventConsumerHandler {
	return &EventConsumerHandler{
		marshaler: protojson.MarshalOptions{
			Multiline:       false,
			EmitUnpopulated: true,
			UseProtoNames:   true,
		},
		logger: logger,
	}
}

func (h *EventConsumerHandler) Register(s *grpc.Server) {
	pb.RegisterEventConsumerServer(s, h)
}

// Implementa o serviço gRPC
func (h *EventConsumerHandler) Store(ctx context.Context, e *pb.Event) (*pb.Empty, error) {
	if e.GetId() == "" {
		return &pb.Empty{}, nil
	}

	h.logger.Info(e.GetId())
	return &pb.Empty{}, nil
}
