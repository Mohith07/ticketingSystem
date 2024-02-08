package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "ticketingSystem/router"
	"time"
)

const PORT = "50051"

func main() {

	conn, err := grpc.Dial("localhost:"+PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	_, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5*time.Second)))
	defer cancel()

	client := pb.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := GetTicketDetails("mohith@cloudbees.com", "random", "random2")

	resp, _ := client.BuyTicket(ctx, &request)

	log.Println("response is {}", resp)

	request2 := GetTicketDetails("bro@cloudbees.com", "random-new", "random-new2")

	resp2, _ := client.BuyTicket(ctx, &request2)

	log.Println("response 2 is {}", resp2)

	ticketId := pb.TicketId{TicketId: resp2.TicketId}

	resp3, _ := client.ShowTicket(ctx, &ticketId)

	log.Println("response 3 is {}", resp3)

	req3 := pb.UserRequest{Email: "mohith@cloudbees.com"}

	resp4, _ := client.RemoveUser(ctx, &req3)

	log.Println("response 4 is {}", resp4)

	req4 := pb.TrainId{TrainId: resp3.TrainNumber}

	resp5, _ := client.ShowUsersBoarded(ctx, &req4)

	log.Println("response 4 is {}", resp5)
}

func GetTicketDetails(email string, firstName string, lastName string) pb.TicketRequest {
	user := pb.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	from := pb.Location{
		LocationCode: 12,
		Name:         "France",
	}

	to := pb.Location{
		LocationCode: 14,
		Name:         "London",
	}

	request := pb.TicketRequest{
		User:      &user,
		From:      &from,
		To:        &to,
		PricePaid: 20,
	}
	return request
}
