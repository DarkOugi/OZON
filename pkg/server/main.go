package server

import (
	"github.com/DarkOugi/OZON/pkg/grpc/pb"
	"github.com/DarkOugi/OZON/pkg/service"
)

type XmlServer struct {
	sv *service.Service
}

func NewXmlServer(sv *service.Service) *XmlServer {
	return &XmlServer{
		sv: sv,
	}
}

type ProtoServer struct {
	pb.UnimplementedMockDailyValueServiceServer
	sv *service.Service
}

func NewProtoServer(sv *service.Service) *ProtoServer {
	return &ProtoServer{
		sv: sv,
	}
}
