package main

import (
	"context"
	"errors"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	"ticketingSystem/db"
	pb "ticketingSystem/router"
)

var (
	port = flag.Int("port", 5001, "The server port")
)

type ticketingSystemServer struct {
	pb.UnimplementedRouteGuideServer
	apiCount int
}

func (s *ticketingSystemServer) BuyTicket(ctx context.Context, req *pb.TicketRequest) (*pb.TicketResponse, error) {
	s.apiCount++
	toLocation := req.To
	fromLocation := req.From
	train := db.GetTrain(fromLocation.GetName(), toLocation.GetName())

	seat, err := db.PurchaseTicket(req.GetUser(), req.GetPricePaid(), train)

	if err != nil {
		return nil, err
	}
	log.Printf("Ticket successfully booked for the user")
	response := pb.TicketResponse{
		TicketId:      seat.TicketNumber,
		SeatNumber:    seat.SeatDetails.ID,
		SectionNumber: seat.SeatDetails.SectionName,
	}
	return &response, nil
}

func (s *ticketingSystemServer) ShowTicket(ctx context.Context, ticket *pb.TicketId) (*pb.Ticket, error) {
	s.apiCount++

	userSeatDetails := db.GetTicketDetails(ticket.TicketId)

	if userSeatDetails == nil {
		log.Fatalf("no tickets exist with the given id")
	}

	tickResponse := &pb.Ticket{
		TrainNumber: userSeatDetails.TrainNumber,
		Seat: &pb.Seat{
			Section:    userSeatDetails.SeatDetails.SectionName,
			SeatNumber: userSeatDetails.SeatDetails.ID,
		},
		User: &pb.User{
			FirstName: userSeatDetails.FirstName,
			LastName:  userSeatDetails.LastName,
			Email:     userSeatDetails.Email,
		},
	}

	return tickResponse, nil
}

func (s *ticketingSystemServer) ShowUsersBoarded(ctx context.Context, req *pb.TrainId) (*pb.SectionResponse, error) {
	train := db.GetTrainById(req.TrainId)
	if train == nil {
		return nil, errors.New("could not find a train from and to the given location")
	}
	details := db.GetAllUsersInTrain(train.ID)
	resp := pb.SectionResponse{
		Name:          train.ID,
		UserToSeatMap: details,
	}
	return &resp, nil
}

func (s *ticketingSystemServer) RemoveUser(ctx context.Context, req *pb.UserRequest) (*pb.BooleanObj, error) {
	s.apiCount++
	if req == nil {
		return &pb.BooleanObj{Response: false}, errors.New("user request is nil")
	}
	email := req.Email
	if email == "" {
		return &pb.BooleanObj{Response: false}, errors.New("user email is empty")
	}

	db.RemoveUser(email)

	db.PrintAllSeats(db.GetTrainById("TRAIN-1000"))

	db.PrintUserSeats(db.GetTrainById("TRAIN-1000"))

	return &pb.BooleanObj{Response: true}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	ser := &ticketingSystemServer{apiCount: 0}
	log.Println("server starting..")
	pb.RegisterRouteGuideServer(grpcServer, ser)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error while loading the service")
		return
	}
}
