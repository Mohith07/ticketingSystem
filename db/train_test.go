package db

import (
	"testing"
)

func TestGetTrain(t *testing.T) {
	// Test positive case: Train exists
	train := &Train{
		ID:    "TRAIN-1000",
		Price: 20,
	}
	location := Location{
		From: "France",
		To:   "London",
	}
	Trains[location] = train

	gotTrain := GetTrain("France", "London")
	if gotTrain == nil || gotTrain.ID != "TRAIN-1000" {
		t.Errorf("GetTrain(France, London) returned unexpected train: %v", gotTrain)
	}

	// Test negative case: Train doesn't exist
	gotTrain = GetTrain("Germany", "Italy")
	if gotTrain != nil {
		t.Errorf("GetTrain(Germany, Italy) expected to return nil, got: %v", gotTrain)
	}
}

func TestGetTrainById(t *testing.T) {
	// Test positive case: Train exists
	train := &Train{
		ID:    "TRAIN-1000",
		Price: 20,
	}
	Trains[Location{From: "France", To: "London"}] = train

	gotTrain := GetTrainById("TRAIN-1000")
	if gotTrain == nil || gotTrain.ID != "TRAIN-1000" {
		t.Errorf("GetTrainById(TRAIN-1000) returned unexpected train: %v", gotTrain)
	}

	// Test negative case: Train doesn't exist
	gotTrain = GetTrainById("INVALID-ID")
	if gotTrain != nil {
		t.Errorf("GetTrainById(INVALID-ID) expected to return nil, got: %v", gotTrain)
	}
}

func TestSaveTrain(t *testing.T) {
	// Test successful save
	location := Location{
		From: "Germany",
		To:   "Spain",
	}
	train := &Train{
		ID:    "TRAIN-2000",
		Price: 30,
	}

	success := SaveTrain(train, location)
	if !success || len(Trains) != 2 {
		t.Errorf("SaveTrain(train, location) failed or didn't save to Trains")
	}

	// Test duplicate ID
	success = SaveTrain(train, location)
	if success || len(Trains) != 2 {
		t.Errorf("SaveTrain(train, location) should have failed due to duplicate ID")
	}
}

func TestCreateSections(t *testing.T) {
	// Test section creation
	train := &Train{
		ID: "TRAIN-3000",
	}

	section := createSections(train, "C")
	if section.Name != "C" || len(section.Seats) != totalSeatsInSection {
		t.Errorf("createSections(train, name) created unexpected section: %v", section)
	}
}

func TestCreateSeats(t *testing.T) {
	// Test seat creation
	train := &Train{
		ID: "TRAIN-4000",
	}

	seats := createSeats(train, "D", 20)
	if len(seats) != 20 || seats[0].ID != "TRAIN-4000-D-1" {
		t.Errorf("createSeats(train, sectionName, count) created unexpected seats: %v", seats)
	}
}

func TestMain(m *testing.M) {
	Trains = make(map[Location]*Train)
	m.Run()
}
