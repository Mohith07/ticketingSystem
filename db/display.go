package db

import "log"

func PrintUserSeats(train *Train) {
	log.Println(UserSeat[train.ID])
}

func PrintAllSeats(train *Train) {
	sectionA := train.Section[0]
	for _, seat := range sectionA.Seats {
		log.Println(seat)
	}
}
