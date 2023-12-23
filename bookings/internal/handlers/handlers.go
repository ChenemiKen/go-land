package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/chenemiken/goland/bookings/helpers"
	"github.com/chenemiken/goland/bookings/internal/config"
	"github.com/chenemiken/goland/bookings/internal/drivers"
	"github.com/chenemiken/goland/bookings/internal/forms"
	"github.com/chenemiken/goland/bookings/internal/models"
	"github.com/chenemiken/goland/bookings/internal/render"
	"github.com/chenemiken/goland/bookings/internal/repository"
	"github.com/chenemiken/goland/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(db *drivers.DB, a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	resvn := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	data := make(map[string]interface{})

	room, err := m.DB.GetRoomById(resvn.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	resvn.Room = room
	m.App.Session.Put(r.Context(), "reservation", resvn)
	data["reservation"] = resvn

	sd := resvn.StartDate.Format("2006-01-02")
	ed := resvn.EndDate.Format("2006-01-02")

	stringData := make(map[string]string)
	stringData["start_date"] = sd
	stringData["end_date"] = ed

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringData,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	resvn := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resvn.FirstName = r.Form.Get("first_name")
	resvn.LastName = r.Form.Get("last_name")
	resvn.Email = r.Form.Get("email")
	resvn.Phone = r.Form.Get("phone")

	sd := resvn.StartDate.Format("2006-01-02")
	ed := resvn.EndDate.Format("2006-01-02")

	stringData := make(map[string]string)
	stringData["start_date"] = sd
	stringData["end_date"] = ed

	form := forms.New(r.PostForm)
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = resvn
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringData,
		})
		return
	}

	newResID, err := m.DB.InsertReservation(resvn)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restriction := models.RoomRestrictions{
		StartDate:     resvn.StartDate,
		EndDate:       resvn.EndDate,
		RoomID:        resvn.RoomID,
		RestrictionID: 1,
		ReservationID: newResID,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", resvn)

	http.Redirect(w, r, "reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) < 1 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	data := make(map[string]interface{})
	data["rooms"] = rooms

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

	// w.Write([]byte(fmt.Sprintf("The selected start is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK        bool      `json:"ok"`
	StartDate time.Time `json:"start"`
	EndDate   time.Time `json:"end"`
	RoomID    int       `json:"room_id"`
}

func (m *Repository) PostAvailabilityJson(w http.ResponseWriter, r *http.Request) {
	sd := r.PostForm.Get("start")
	ed := r.PostForm.Get("end")
	startDate, _ := time.Parse("2006-01-02", sd)
	endDate, _ := time.Parse("2006-01-02", ed)
	roomId, _ := strconv.Atoi(r.PostForm.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)

	resp := jsonResponse{
		OK:        available,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	resvn, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("reservation not found")
		m.App.Session.Put(r.Context(), "error", "no reservation found")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = resvn

	sd := resvn.StartDate.Format("2006-01-02")
	ed := resvn.EndDate.Format("2006-01-02")

	stringData := make(map[string]string)
	stringData["start_date"] = sd
	stringData["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringData,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("unable to get reservation details from session"))
	}

	res.RoomID = roomId

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	startDate, _ := time.Parse("2006-01-02", sd)
	endDate, _ := time.Parse("2006-01-02", ed)

	reservation := models.Reservation{
		RoomID:    roomID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
