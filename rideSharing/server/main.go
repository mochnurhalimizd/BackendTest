package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbDriver "test-ice-house/driver/pb"
	pbPassanger "test-ice-house/passanger/pb"
	pb "test-ice-house/rideSharing/pb"
)

const (
	port          = ":50051"
	portDriver    = ":50053"
	portPassanger = ":50054"
)

type server struct {
	driverClient    pbDriver.DriverClient
	passangerClient pbPassanger.PassagerClient
}

func (s *server) Connect(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	log.Printf("Connect {from : %s}", req.UserName)
	return &pb.UserResponse{UserName: req.UserName}, nil
}

func (s *server) SendPresence(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Send Presence {from : %s}", req.UserName)
	return &pb.UserResponse{UserName: req.UserName}, nil
}

func (s *server) SendRequest(ctx context.Context, Req *pb.RequestUser) (*pb.RequestResponse, error) {
	_, err := s.driverClient.ReceiveRequest(ctx, &pbDriver.RequestUser{
		From:   Req.From,
		Status: Req.Status,
		Lat:    Req.Lat,
		Lon:    Req.Lon,
	})

	log.Printf("Send Request(From: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.Lat, Req.Lon)
	if err != nil {
		log.Fatalf("Send Request Error : %v", err)
	}

	return &pb.RequestResponse{From: Req.From, Success: true}, nil
}

func (s *server) SendLocation(ctx context.Context, Req *pb.RequestLocation) (*pb.ResponseLocation, error) {
	_, err := s.passangerClient.ReceiveLocation(ctx, &pbPassanger.RequestLocation{
		From: Req.From,
		To:   Req.To,
		Lat:  Req.Lat,
		Lon:  Req.Lon,
	})

	if err != nil {
		log.Fatalf("Send location Error %s :", err)
	}

	log.Printf("Send Location(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", Req.From, Req.To, Req.Lat, Req.Lon)

	return &pb.ResponseLocation{From: Req.From, Success: true}, nil

}

func (s *server) AcceptRequest(ctx context.Context, req *pb.AcceptUser) (*pb.AcceptResponse, error) {
	_, err := s.passangerClient.ReceiveRequest(ctx, &pbPassanger.AcceptUser{
		From: req.From,
		To:   req.To,
		Lat:  req.Lat,
		Lon:  req.Lon,
	})

	if err != nil {
		log.Fatalf("Accept Request Error %s :", err)
	}

	log.Printf("Accept Request(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", req.From, req.To, req.Lat, req.Lon)

	return &pb.AcceptResponse{From: req.From, Success: true}, nil
}

func (s *server) StartTrip(ctx context.Context, Req *pb.RequestStartTrip) (*pb.ResponseTrip, error) {
	log.Printf("Start Trip(From: %s, To: %s)", Req.From, Req.To)
	return &pb.ResponseTrip{From: Req.From, Success: true}, nil

}

func (s *server) EndTrip(ctx context.Context, Req *pb.RequestEndTrip) (*pb.ResponseTrip, error) {
	log.Printf("End Trip(From: %s, To: %s, Distance: %d)", Req.From, Req.To, Req.Distance)
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
	pb.RegisterRideSharingServer(srv, &server{
		driverClient:    pbDriver.NewDriverClient(mustDial(portDriver)),
		passangerClient: pbPassanger.NewPassagerClient(mustDial(portPassanger)),
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
