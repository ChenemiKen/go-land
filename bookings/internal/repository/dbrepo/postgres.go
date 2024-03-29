package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/chenemiken/goland/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone,
		start_date, end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt, res.FirstName, res.LastName,
		res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID,
		time.Now(), time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(rr models.RoomRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id,
		restriction_id, reservation_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt, rr.StartDate, rr.EndDate, rr.RoomID,
		rr.RestrictionID, rr.ReservationID, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(startDate,
	endDate time.Time, roomId int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `select count(id) from room_restrictions
		where room_id = $1 and 
		$2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomId, startDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) (
	[]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select r.id, r.room_name from rooms r
		where r.id not in
		(select rr.room_id from room_restrictions rr 
			where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select r.id, r.room_name, r.created_at, r.updated_at from rooms r
		where r.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, err
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User

	query := `select first_name, last_name, email, password, access_level,
		created_at, updated_at from users where id = $1;`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.AccessLevel,
		&u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return u, err
	}
	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3,
		access_level = $4, updated_at = $5) where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now(), u.ID)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1;`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
		r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
		r.processed, rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Phone,
			&r.StartDate, &r.EndDate, &r.RoomID, &r.CreatedAt, &r.UpdatedAt,
			&r.Processed, &r.Room.ID, &r.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
		r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
		r.processed, rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where processed = 0
		order by r.start_date asc;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Phone,
			&r.StartDate, &r.EndDate, &r.RoomID, &r.CreatedAt, &r.UpdatedAt,
			&r.Processed, &r.Room.ID, &r.Room.RoomName)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone,
	r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at,
	r.processed, rm.id, rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where r.id = $1;`

	row := m.DB.QueryRowContext(ctx, query, id)

	var r models.Reservation
	err := row.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Email, &r.Phone,
		&r.StartDate, &r.EndDate, &r.RoomID, &r.CreatedAt, &r.UpdatedAt,
		&r.Processed, &r.Room.ID, &r.Room.RoomName)

	if err != nil {
		return r, err
	}

	return r, nil
}

func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name = $1, last_name = $2, email = $3,
		phone = $4, updated_at = $5 where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) UpdateReservationProcessed(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, room_name, created_at, updated_at from rooms
		order by room_name`

	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Room
		err := rows.Scan(&r.ID, &r.RoomName, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, r)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDBRepo) GetRestrictionForRoomByDate(roomId int,
	start, end time.Time) ([]models.RoomRestrictions, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var roomRestrictions []models.RoomRestrictions

	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, 
		start_date, end_date from room_restrictions where $1 < end_date and 
		$2 >= start_date and room_id = $3`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomId)
	if err != nil {
		return roomRestrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var rr models.RoomRestrictions
		err := rows.Scan(&rr.ID, &rr.ReservationID, &rr.RestrictionID,
			&rr.RoomID, &rr.StartDate, &rr.EndDate)
		if err != nil {
			return roomRestrictions, err
		}

		roomRestrictions = append(roomRestrictions, rr)
	}

	if err = rows.Err(); err != nil {
		return roomRestrictions, err
	}

	return roomRestrictions, nil
}
