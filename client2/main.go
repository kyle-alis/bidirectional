package main

import (
	"context"
	"fmt"
	"io"
	"time"

	users "github.com/kyle-alis/bidirectional/tree/main/proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client Stream...")
	fmt.Println()

	opts := grpc.WithInsecure()
	con, err := grpc.Dial("localhost:3000", opts)
	if err != nil {
		fmt.Errorf("Error connecting: %v \n", err)
	}

	defer con.Close()
	c := users.NewUsersClient(con)
	bulkUsers(c)
}

// container struct
type container struct {
	users []*users.User
}

// builkUsers function
func bulkUsers(c users.UsersClient) {

	// Get the stream and a possible error from the CreateUser function
	stream, err := c.CreateUser(context.Background())
	if err != nil {
		fmt.Errorf("Error when getting stream object: %v", err)
		return
	}

	// Initialize the container struct and call the initUsers function
	// to get user objects to send on the request message.
	requests := container{}.initUsers()

	// Create a new channel
	waitResponse := make(chan struct{})

	// Use a go routine to send request messages to the server
	go func() {
		// Iterate over the requests slice
		for _, req := range requests {
			// Send request message
			stream.Send(req)

			// Sleep for a little bit..
			time.Sleep(1000 * time.Millisecond)
		}
		// Close stream
		stream.CloseSend()
	}()

	// Use a go routine to receive response messages from the server
	go func() {
		for {
			// Get response and possible error message from the stream
			res, err := stream.Recv()

			// Break for loop if there are no more response messages
			if err == io.EOF {
				break
			}

			// Handle a possible error
			if err != nil {
				fmt.Errorf("Error when receiving response: %v", err)
			}

			// Log the response
			fmt.Println("Server Response: ", res)
		}

		// Close channel
		close(waitResponse)
	}()
	<-waitResponse
}

// initUsers function
func (c container) initUsers() []*users.User {
	c.users = append(c.users, c.getUser("2", "Carl", "Phill", 23))
	c.users = append(c.users, c.getUser("2", "Marisol", "Richardson", 29))
	c.users = append(c.users, c.getUser("2", "Mia", "Phill", 27))
	c.users = append(c.users, c.getUser("2", "Tomas", "Smith", 25))
	c.users = append(c.users, c.getUser("2", "Zian", "Heat", 28))

	return c.users
}

// getUser function
func (c container) getUser(id, name, lastName string, age int32) *users.User {
	return &users.User{
		Id:       id,
		Name:     name,
		LastName: lastName,
		Age:      age,
	}
}
