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
	pbDriver "test-ice-house/driver/pb"
	pb "test-ice-house/rideSharing/pb"
)

const (
	address    = "localhost:50051"
	portDriver = ":50053"
)

// The person Type (more like an object)
type Driver struct {
	From string `json:"from"`
}

type Request struct {
	From string  `json:"From"`
	Lat  float32 `json:"Lat"`
	Lon  float32 `json:"Lon"`
}

type RequestAccept struct {
	From string
	To   string
	Lat  float32
	Lon  float32
}

type RequestLocation struct {
	From string
	To   string
	Lat  float32
	Lon  float32
}

type StartTripRequest struct {
	From string
	To   string
}

type EndTripRequest struct {
	From     string
	To       string
	Distance int32
}

type RequestResponseLocation struct {
	From    string
	To      string
	Success bool
}

type RequestResponseAccept struct {
	From    string
	To      string
	Success bool
}

type ResponseTrip struct {
	From   string `json:"From"`
	Succes bool   `json:"Success"`
}

var drivers []Driver

var requests []Request

// create a new item
func Connect(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var driver Driver
	_ = json.NewDecoder(r.Body).Decode(&driver)
	driver.From = params["from"]
	//driver.From = "Alice"
	drivers = append(drivers, driver)

	//Connect to gRPC
	conn, err := grpc.Dial(address, grpc.WithInsecure())

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

func getAllRequest(w http.ResponseWriter, r *http.Request) {

	var request Request
	_ = json.NewDecoder(r.Body).Decode(&request)

	conn, err := grpc.Dial(portDriver, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbDriver.NewDriverClient(conn)

	filter := &pbDriver.RequestFilter{UserKey: ""}
	stream, err := c.GetRequestUser(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
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

		request.From = Request.From
		request.Lat = Request.Lat
		request.Lon = Request.Lon

	}

	json.NewEncoder(w).Encode(request)

}

func sendLocation(w http.ResponseWriter, r *http.Request) {
	var requestResponse RequestResponseLocation
	var request RequestLocation
	//var request Request

	params := mux.Vars(r)
	request.From = params["from"]
	request.To = params["to"]

	lat, _ := strconv.ParseFloat(params["lat"], 32)

	request.Lat = float32(lat)

	lon, _ := strconv.ParseFloat(params["lon"], 32)
	request.Lon = float32(lon)

	_ = json.NewDecoder(r.Body).Decode(&requestResponse)

	conn, err := grpc.Dial(portDriver, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbDriver.NewDriverClient(conn)

	result, err := c.SendLocation(context.Background(), &pbDriver.RequestLocation{
		From: request.From,
		Lat:  request.Lat,
		Lon:  request.Lon,
		To:   request.To,
	})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Send Location(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", request.From, request.To, request.Lat, request.Lon)

	requestResponse.From = request.From
	requestResponse.Success = result.Success
	requestResponse.To = request.To
	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(requestResponse)
}

func acceptRequest(w http.ResponseWriter, r *http.Request) {
	var requestResponse RequestResponseAccept
	var requestAccept RequestAccept
	//var request Request

	params := mux.Vars(r)
	requestAccept.From = params["from"]
	requestAccept.To = params["to"]

	lat, _ := strconv.ParseFloat(params["lat"], 32)

	requestAccept.Lat = float32(lat)

	lon, _ := strconv.ParseFloat(params["lon"], 32)
	requestAccept.Lon = float32(lon)

	_ = json.NewDecoder(r.Body).Decode(&requestResponse)

	conn, err := grpc.Dial(portDriver, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbDriver.NewDriverClient(conn)

	result, err := c.AcceptRequest(context.Background(), &pbDriver.AcceptUser{
		From: requestAccept.From,
		Lat:  requestAccept.Lat,
		Lon:  requestAccept.Lon,
		To:   requestAccept.To,
	})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Accept Request(From: %s, To: %s, Lat: %.3f, Lon: %.3f )", result.From, requestAccept.To, requestAccept.Lat, requestAccept.Lon)

	requestResponse.From = result.From
	requestResponse.Success = result.Success
	requestResponse.To = requestAccept.To
	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(requestResponse)

}

func startTrip(w http.ResponseWriter, r *http.Request) {
	var requestResponse ResponseTrip
	var request StartTripRequest
	//var request Request

	params := mux.Vars(r)

	request.From = params["from"]
	request.To = params["to"]

	_ = json.NewDecoder(r.Body).Decode(&requestResponse)

	conn, err := grpc.Dial(portDriver, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbDriver.NewDriverClient(conn)

	result, err := c.StartTrip(context.Background(), &pbDriver.RequestStartTrip{
		From: request.From,
		To:   request.To,
	})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("Start Trip(From: %s, To: %s)", result.From, request.To)

	requestResponse.From = result.From
	requestResponse.Succes = result.Success
	//log.Printf("Connect to Ride Sharing From : %s", driver)
	json.NewEncoder(w).Encode(requestResponse)

}

func endTrip(w http.ResponseWriter, r *http.Request) {
	var requestResponse ResponseTrip
	var request EndTripRequest
	//var request Request

	params := mux.Vars(r)

	request.From = params["from"]
	request.To = params["to"]

	distance, _ := strconv.ParseInt(params["distance"], 0, 32)
	request.Distance = int32(distance)

	_ = json.NewDecoder(r.Body).Decode(&requestResponse)

	conn, err := grpc.Dial(portDriver, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Did not connect :%v", err)
	}

	defer conn.Close()

	c := pbDriver.NewDriverClient(conn)

	result, err := c.EndTrip(context.Background(), &pbDriver.RequestEndTrip{
		From:     request.From,
		To:       request.To,
		Distance: request.Distance,
	})
	if err != nil {
		log.Fatalf("Could not greet : %v", err)
	}

	log.Printf("End Trip(From: %s, To: %s, Distance %d)", result.From, request.To, request.Distance)

	requestResponse.From = result.From
	requestResponse.Succes = result.Success

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
	conn, err := grpc.Dial(address, grpc.WithInsecure())

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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/connect/{from}", Connect).Methods("GET")
	router.HandleFunc("/getAllRequest", getAllRequest).Methods("GET")
	router.HandleFunc("/sendPresence/{from}", sendPresence).Methods("GET")
	router.HandleFunc("/acceptRequest/{from}/{to}/{lat}/{lon}", acceptRequest).Methods("GET")
	router.HandleFunc("/sendLocation/{from}/{to}/{lat}/{lon}", sendLocation).Methods("GET")
	router.HandleFunc("/startTrip/{from}/{to}", startTrip).Methods("GET")
	router.HandleFunc("/endTrip/{from}/{to}/{distance}", endTrip).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
