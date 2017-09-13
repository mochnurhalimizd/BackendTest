Instalation 

1. go get golang.org/x/net/context
2. go get google.golang.org/grpc

Running Service Ride Sharing

1. cd/{HOME}/test-ice-house/rideSharing/server
2. go run main.go

Run Service Passenger

1. cd/{HOME}/test-ice-house/passanger/server
2. go run main.go

Run API Passenger

1. cd/{HOME}/test-ice-house/passanger/API
2. go run main.go

Run Service Driver

1. cd/{HOME}/test-ice-house/driver/server
2. go run main.go

Run API Driver

1. cd/{HOME}/test-ice-house/driver/API
2. go run main.go




Running Rest API Driver

localhost:8000/connect/{from} GET
localhost:8000/getAllRequest GET
localhost:8000/sendPresence/{from} GET
localhost:8000/acceptRequest/{from}/{to}/{lat}/{lon} GET
localhost:8000/sendLocation/{from}/{to}/{lat}/{lon} GET
localhost:8000/startTrip/{from}/{to} GET 
localhost:8000/endTrip/{from}/{to}/{distance} GET


Running Rest API Passenger

localhost:8001//connect/{from} GET
localhost:8001//sendPresence/{from} GET
localhost:8001//sendRequest/{from}/{lat}/{lon} GET
localhost:8001//getAllRequest GET
localhost:8001//getLocationDriver GET

