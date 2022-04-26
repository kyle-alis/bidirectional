package main

import (
	"fmt"
	"io"
	"net"

	users "github.com/kyle-alis/bidirectional/tree/main/proto"
	"google.golang.org/grpc"
)

type server struct {
	users.UnimplementedUsersServer
}

func main() {
	fmt.Println("Starting the server...")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Errorf("Unable to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	users.RegisterUsersServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		fmt.Errorf("Failed to serve: %v", err)
	}
}

// CreateUser function
func (*server) CreateUser(stream users.Users_CreateUserServer) error {
	fmt.Println("CreateUser Function")

	for {
		// Receive the request and possible error from the stream object
		req, err := stream.Recv()

		// If there are no more requests, we return
		if err == io.EOF {
			return nil
		}

		// Handle error from the stream object
		if err != nil {
			fmt.Errorf("Error when reading client request stream: %v", err)
		}

		// Get name, last name and user id, form the request
		name, lastName, userID := req.GetName(), req.GetLastName(), req.GetId()
		fmt.Printf("Request: name: %v, last_name: %v, id: %v \n", name, lastName, userID)

		// Initialize the errors and success variables
		errors := []string{}
		success := true

		// Run some validations
		if len(name) <= 3 {
			errors = append(errors, "Name is too short")
			success = false
		}

		if lastName == "Phill" {
			errors = append(errors, "Last Name already taken")
			success = false
		}

		// Build and send response to the client
		res := stream.Send(&users.CreateUserRes{
			UserId:  userID,
			Success: success,
			Errors:  errors,
		})

		// Handle any possible error, when sending the response
		if res != nil {
			fmt.Errorf("Error when response was sent to the client: %v", res)
		}
	}
}
