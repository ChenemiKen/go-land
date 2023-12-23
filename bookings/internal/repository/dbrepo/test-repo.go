package dbrepo

import (
	"errors"
	"time"

	"github.com/chenemiken/goland/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(rr models.RoomRestrictions) error {
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(startDate,
	endDate time.Time, roomId int) (bool, error) {
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) (
	[]models.Room, error) {

	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("error getting room by id")
	}
	return room, nil
}
