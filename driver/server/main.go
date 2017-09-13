package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	pb "test-ice-house/driver/pb"
	pbRide "test-ice-house/rideSharing/pb"
)

const (
	port            = ":50053"
	portRideSharing = ":50051"
)

type server struct {
	requestUser []*pb.RequestUser
	rideClient  pbRide.RideSharingClient
}

func (s *server) GetRequestUser(req *pb.RequestFilter, stream pb.Driver_GetRequestUserServer) error {
	for _, request := range s.requestUser {
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

func (s *server) ReceiveRequest(ctx context.Context, req *pb.RequestUser) (*pb.RequestResponse, error) {
	s.requestUser = append(s.requestUser, req)
	log.Printf("Receive Request(From: %s, Lat: %.3f, Lon: %.3f )", req.From, req.Lat, req.Lon)
	return &pb.RequestResponse{From: req.From, Success: true}, nil
}

func (s *server) AcceptRequest(ctx context.Context, Req *pb.AcceptUser) (*pb.AcceptResponse, error) {
	_, err := s.rideClient.AcceptRequest(ctx, &pbRide.AcceptUser{
		From: Req.From,
		To:   Req.To,
		Lat:  Req.Lat,
		Lon:  Req.Lon,
	})
	log.Printf("Accept Request(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.To, Req.Lat, Req.Lon)

	if err != nil {
		log.Fatalf("Accept Request Error : %v", err)
	}

	return &pb.AcceptResponse{From: Req.From, Success: true}, nil
}

func (s *server) SendLocation(ctx context.Context, Req *pb.RequestLocation) (*pb.ResponseLocation, error) {
	_, err := s.rideClient.SendLocation(ctx, &pbRide.RequestLocation{
		From: Req.From,
		To:   Req.To,
		Lat:  Req.Lat,
		Lon:  Req.Lon,
	})

	if err != nil {
		log.Fatalf("Send Location Error Error : %v", err)
	}

	log.Printf("Send Location(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.To, Req.Lat, Req.Lon)
	return &pb.ResponseLocation{From: Req.From, Success: true}, nil

}

func (s *server) EndTrip(ctx context.Context, Req *pb.RequestEndTrip) (*pb.ResponseTrip, error) {
	_, err := s.rideClient.EndTrip(ctx, &pbRide.RequestEndTrip{
		From:     Req.From,
		To:       Req.To,
		Distance: Req.Distance,
	})

	if err != nil {
		log.Fatalf("End Trip Error : %v", err)
	}

	log.Printf("End Trip(From: %s, To: %s, Distance: %d)", Req.From, Req.To, Req.Distance)
	return &pb.ResponseTrip{From: Req.From, Success: true}, nil
}

func (s *server) StartTrip(ctx context.Context, Req *pb.RequestStartTrip) (*pb.ResponseTrip, error) {
	_, err := s.rideClient.StartTrip(ctx, &pbRide.RequestStartTrip{
		From: Req.From,
		To:   Req.To,
	})

	if err != nil {
		log.Fatalf("Start Trip Error : %v", err)
	}

	log.Printf("Start Trip(From: %s, To: %s)", Req.From, Req.To)
	return &pb.ResponseTrip{From: Req.From, Success: true}, nil
}

func main() {

	// tcp listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpc server
	srv := grpc.NewServer()
	pb.RegisterDriverServer(srv, &server{
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
