package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/chenemiken/goland/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("invalid room Id")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(rr models.RoomRestrictions) error {
	if rr.RoomID == 1000 {
		return errors.New("invalid room Id")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(startDate,
	endDate time.Time, roomId int) (bool, error) {

	if roomId == 2 {
		return true, nil
	}
	if roomId == 3 {
		return false, errors.New("error-room3")
	}
	return false, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) (
	[]models.Room, error) {

	var rooms []models.Room
	sd, err := time.Parse("2006-01-02", "2025-01-01")
	if err != nil {
		return nil, err
	}
	ed, err := time.Parse("2006-01-02", "2025-12-12")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if start == sd && end == ed {
		room := models.Room{
			RoomName: "General's quarters",
		}
		rooms = append(rooms, room)
	}

	sd, err = time.Parse("2006-01-02", "2026-01-01")
	if err != nil {
		return nil, err
	}
	ed, err = time.Parse("2006-01-02", "2026-12-12")
	if err != nil {
		return nil, err
	}
	if start == sd && end == ed {
		return rooms, errors.New("failed to search rooms")
	}

	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("error getting room by id")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}
