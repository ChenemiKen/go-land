package repository

import (
	"time"

	"github.com/chenemiken/goland/bookings/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(rr models.RoomRestrictions) error

	SearchAvailabilityByDatesByRoomId(startDate, endDate time.Time,
		roomId int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)
}
