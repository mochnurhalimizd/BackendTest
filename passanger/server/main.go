package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	pb "test-ice-house/passanger/pb"
	pbRide "test-ice-house/rideSharing/pb"
)

const (
	port            = ":50054"
	portRideSharing = ":50051"
)

type server struct {
	rideClient      pbRide.RideSharingClient
	requestUser     []*pb.RequestUser
	acceptRequest   []*pb.AcceptUser
	locationRequest []*pb.RequestLocation
}

func (s *server) SendRequest(ctx context.Context, Req *pb.RequestUser) (*pb.RequestResponse, error) {
	_, err := s.rideClient.SendRequest(ctx, &pbRide.RequestUser{
		From:   Req.From,
		Status: Req.Status,
		Lat:    Req.Lat,
		Lon:    Req.Lon,
	})

	if err != nil {
		log.Fatalf("Accept Request Error : %v", err)
	}

	log.Printf("Send Request(From: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.Lat, Req, Req.Lon)
	return &pb.RequestResponse{From: Req.From, Success: true}, nil
}

func (s *server) GetRequestStatus(req *pb.RequestFilter, stream pb.Passager_GetRequestStatusServer) error {
	for _, request := range s.acceptRequest {
		if req.UserKey != "" {
			if !strings.Contains(request.From, req.UserKey) {
				continue
			}
		}

		if err := stream.Send(request); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) GetLocation(Req *pb.RequestFilter, stream pb.Passager_GetLocationServer) error {

	for _, request := range s.locationRequest {
		if Req.UserKey != "" {
			if !strings.Contains(request.From, Req.UserKey) {
				continue
			}
		}

		if err := stream.Send(request); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ReceiveRequest(ctx context.Context, req *pb.AcceptUser) (*pb.AcceptResponse, error) {

	for _, request := range s.requestUser {
		if req.From != "" {
			if strings.Contains(request.From, req.From) {
				request.Status = "Found The Driver"
			}
		}
	}

	s.acceptRequest = append(s.acceptRequest, req)

	log.Printf("Receive Request(From: %s, Lat: %.3f, Lon: %.3f )", req.To, req.Lat, req.Lon)

	return &pb.AcceptResponse{From: req.From, Success: true}, nil
}

func (s *server) ReceiveLocation(ctx context.Context, Req *pb.RequestLocation) (*pb.ResponseLocation, error) {
	s.locationRequest = append(s.locationRequest, Req)
	log.Printf("Receive Location(From: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.Lat, Req.Lon)
	return &pb.ResponseLocation{From: Req.From, Success: true}, nil
}

func main() {

	// tcp listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc server
	srv := grpc.NewServer()
	pb.RegisterPassagerServer(srv, &server{
		rideClient: pbRide.NewRideSharingClient(mustDial(portRideSharing)),
	})
	srv.Serve(lis)

}

func mustDial(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		panic(err)
	}
	return conn
}
