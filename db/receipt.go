package db

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	pb "ticketingSystem/router"
)

type UserSeatDetails struct {
	FirstName    string
	LastName     string
	Email        string
	AmountPaid   int32
	TicketNumber string
	TrainNumber  string
	SeatDetails  *Seat
}

var (
	receipt = make(map[string]*UserSeatDetails)
)

func GenerateTicket(seat *Seat, amountCharged int32, user *pb.User) *UserSeatDetails {
	// Generate a unique ticket number (you can use a more sophisticated method for production)
	ticketNumber := randomString(8)

	// Create a UserSeatDetails instance
	userDetails := &UserSeatDetails{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		AmountPaid:   amountCharged,
		TicketNumber: ticketNumber,
		TrainNumber:  seat.TrainName,
		SeatDetails:  seat,
	}

	// Store the UserSeatDetails in the receipt map
	receipt[ticketNumber] = userDetails

	return userDetails
}

func GetTicketDetails(ticketId string) *UserSeatDetails {
	if receipt[ticketId] != nil {
		return receipt[ticketId]
	}
	return nil
}

func PurchaseTicket(user *pb.User, amount int32, train *Train) (*UserSeatDetails, error) {
	if train == nil {
		return nil, errors.New("could not find a train from and to the given location")
	}
	if train.Price > amount {
		return nil, errors.New("price paid is less than the price of the train ticket")
	}
	log.Printf("Booking ticket for the user")
	seat, err := ReserveSeat(user.GetEmail(), *train)
	if err != nil || seat == nil {
		return nil, err
	}
	log.Printf("Generating receipt number")
	return GenerateTicket(seat, amount, user), nil
}

func ReserveSeat(email string, train Train) (*Seat, error) {
	// check section A and then section B for seats
	sectionA := train.Section[0]
	sectionB := train.Section[1]
	trainDetails := getTrainDetails(train.ID)
	for _, seat := range sectionA.Seats {
		if !seat.IsAssigned {
			trainDetails[email] = seat
			seat.IsAssigned = true
			log.Printf("Assigned seat to the user with email %s in section A", email)
			return seat, nil
		}
	}

	for _, seat := range sectionB.Seats {
		if !seat.IsAssigned {
			trainDetails[email] = seat
			log.Printf("Assigned seat to the user with email %s in section B", email)
			return seat, nil
		}
	}
	return nil, errors.New("could not find a seat")
}

func RemoveUser(email string) {
	for _, trainDetails := range UserSeat {
		if seat, ok := trainDetails[email]; ok {
			// mark the seat is unavailable
			seat.IsAssigned = false
			trainDetails[email] = seat
			delete(trainDetails, email)
			log.Printf("Removed user with email %s", email)
		}
	}
}

func GetAllUsersInTrain(trainId string) map[string]string {
	tempUserDetails := make(map[string]string)

	userDetailsForCurrTrain := UserSeat[trainId]

	for user, seat := range userDetailsForCurrTrain {
		tempUserDetails[user] = seat.ID
	}

	return tempUserDetails
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string
	for i := 0; i < length; i++ {
		// Generate a random index within the charset length
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result += string(charset[index.Int64()])
	}
	return result
}
