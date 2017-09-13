package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"strconv"
	pbPassanger "test-ice-house/passanger/pb"
	pb "test-ice-house/rideSharing/pb"
)

const (
	rideSharingPort = "localhost:50051"
	passangerPort   = "localhost:50054"
)

// The person Type (more like an object)
type Driver struct {
	From string `json:"from"`
}

type RequestData struct {
	From   string  `json:"From"`
	Lat    float32 `json:"lat"`
	Lon    float32 `json:"lon"`
	Status string  `json:"Status"`
}

type AcceptRequest struct {
	From string  `json:"From"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
	To   string  `json:"To"`
}

type RequestLocation struct {
	From string  `json:"From"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
	To   string  `json:"To"`
}

type RequestResponse struct {
	From    string `json:"From"`
	Success bool   `json:"Success"`
}

var drivers []Driver

var requests []RequestData

// create a new item
func Connect(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var driver Driver
	_ = json.NewDecoder(r.Body).Decode(&driver)
	//driver.From = params["from"]
	params := mux.Vars(r)
	driver.From = params["from"]
	drivers = append(drivers, driver)

	//Connect to gRPC
	conn, err := grpc.Dial(rideSharingPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pb.NewRideSharingClient(conn)

	result, err := c.Connect(context.Background(), &pb.UserRequest{UserName: driver.From})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Connet to Ride Sharing Username :%s", result.UserName)

	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(driver)
}

func sendRequest(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//var driver Driver

	var request RequestData
	var requestResponse RequestResponse
	_ = json.NewDecoder(r.Body).Decode(&requestResponse)
	//driver.From = params["from"]
	params := mux.Vars(r)

	lat, _ := strconv.ParseFloat(params["lat"], 32)
	request.Lat = float32(lat)

	request.From = params["from"]

	lon, _ := strconv.ParseFloat(params["lon"], 32)
	request.Lon = float32(lon)
	request.Status = "Waiting"

	//driver.From = "Nurhalim"
	requests = append(requests, request)

	//Connect to gRPC
	conn, err := grpc.Dial(passangerPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbPassanger.NewPassagerClient(conn)

	result, err := c.SendRequest(context.Background(), &pbPassanger.RequestUser{
		From:   request.From,
		Lat:    request.Lat,
		Lon:    request.Lon,
		Status: request.Status,
	})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Send Request(From: %s, Lat: %.3f, Lon: %.3f )", result.From, request.Lat, request.Lon)

	requestResponse.From = result.From
	requestResponse.Success = result.Success
	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(requestResponse)
}

func sendPresence(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var driver Driver
	_ = json.NewDecoder(r.Body).Decode(&driver)
	//driver.From = params["from"]
	params := mux.Vars(r)
	driver.From = params["from"]
	drivers = append(drivers, driver)

	//Connect to gRPC
	conn, err := grpc.Dial(rideSharingPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pb.NewRideSharingClient(conn)

	result, err := c.SendPresence(context.Background(), &pb.UserRequest{UserName: driver.From})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Send Presence to Ride Sharing Username :%s", result.UserName)

	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(driver)
}

func getLocationDriver(w http.ResponseWriter, r *http.Request) {
	var request RequestLocation
	_ = json.NewDecoder(r.Body).Decode(&request)

	conn, err := grpc.Dial(passangerPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbPassanger.NewPassagerClient(conn)

	filter := &pbPassanger.RequestFilter{UserKey: ""}
	stream, err := c.GetLocation(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get Data Request: %v", err)
	}
	for {
		// Receiving the stream of data
		Request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf(" Error : %v", err)
		}
		log.Printf("Got Driver Location: %v", Request)

		request.From = Request.From
		request.Lat = Request.Lat
		request.Lon = Request.Lon
		request.To = Request.To
	}

	json.NewEncoder(w).Encode(request)

}

func getAllRequest(w http.ResponseWriter, r *http.Request) {

	var request AcceptRequest
	_ = json.NewDecoder(r.Body).Decode(&request)

	conn, err := grpc.Dial(passangerPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbPassanger.NewPassagerClient(conn)

	filter := &pbPassanger.RequestFilter{UserKey: ""}
	stream, err := c.GetRequestStatus(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get Data Request: %v", err)
	}
	for {
		// Receiving the stream of data
		Request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf(" Error : %v", err)
		}
		log.Printf("Request: %v", Request)

		request.From = Request.To
		request.Lat = Request.Lat
		request.Lon = Request.Lon
		request.To = Request.From
	}

	json.NewEncoder(w).Encode(request)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/connect/{from}", Connect).Methods("GET")
	router.HandleFunc("/sendPresence/{from}", sendPresence).Methods("GET")
	router.HandleFunc("/sendRequest/{from}/{lat}/{lon}", sendRequest).Methods("GET")
	router.HandleFunc("/getAllRequest", getAllRequest).Methods("GET")
	router.HandleFunc("/getLocationDriver", getLocationDriver).Methods("GET")
	log.Fatal(http.ListenAndServe(":8001", router))
}
