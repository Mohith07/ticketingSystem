package db

import (
	"log"
	"strconv"
)

type Train struct {
	ID      string `json:"id"`
	Price   int32
	Section []*Section
}

type Section struct {
	Name  string
	Seats []*Seat
}

type Seat struct {
	ID          string
	SectionName string
	TrainName   string
	IsAssigned  bool
}

type Location struct {
	From string
	To   string
}

const totalSeatsInSection = 50

var (
	UserSeat = make(map[string]map[string]*Seat)
	Trains   = make(map[Location]*Train)
)

func GetTrain(from string, to string) *Train {
	for location, train := range Trains {
		if location.To == to && location.From == from {
			return train
		}
	}
	return nil
}

func GetTrainById(id string) *Train {
	for _, train := range Trains {
		if train.ID == id {
			return train
		}
	}
	return nil
}

func getTrainDetails(id string) map[string]*Seat {
	if seat, ok := UserSeat[id]; ok {
		return seat
	}
	newMap := make(map[string]*Seat)
	UserSeat[id] = newMap
	return newMap
}

func SaveTrain(train *Train, location Location) bool {

	if _, ok := Trains[location]; ok {
		return false
	}
	var sections []*Section
	// create 2 sections

	sectionA := createSections(train, "A")
	sectionB := createSections(train, "B")

	sections = append(sections, sectionA)
	sections = append(sections, sectionB)
	train.Section = sections

	Trains[location] = train
	return true
}

func createSections(train *Train, name string) *Section {
	seats := createSeats(train, name, totalSeatsInSection)
	section := &Section{
		Name:  name,
		Seats: seats,
	}
	return section
}

func createSeats(train *Train, sectionName string, count int) []*Seat {
	var seats []*Seat
	for i := 1; i <= count; i++ {
		intValue := strconv.Itoa(i)
		seats = append(seats, &Seat{
			ID:          train.ID + "-" + sectionName + "-" + intValue,
			SectionName: sectionName,
			TrainName:   train.ID,
			IsAssigned:  false,
		})
	}
	return seats
}

func init() {
	train := &Train{
		ID:    "TRAIN-1000",
		Price: 20,
	}
	location := Location{
		From: "France",
		To:   "London",
	}
	success := SaveTrain(train, location)
	if success {
		log.Printf("Created a train from France to London with a fare of 20$")
	} else {
		log.Printf("Failed to create any trains")
	}
}
