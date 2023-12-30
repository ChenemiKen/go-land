package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/chenemiken/goland/bookings/internal/drivers"
	"github.com/chenemiken/goland/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"rs", "/reservation-summary", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestNewRepo(t *testing.T) {
	var db drivers.DB
	testRepo := NewRepo(&db, &app)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, "+
			"wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name,
				e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepositoryReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
	}
	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getctx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("make reservation did not return appropriate response code"+
			"expected %d but got %d", http.StatusOK, rr.Code)
	}

	//Testing a case where the reservation is not in session
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusTemporaryRedirect {
		t.Errorf("make reservation did not return appropriate response code"+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	//Testing the case where could not get room by Id
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusTemporaryRedirect {
		t.Errorf("make reservation did not return appropriate response code"+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func TestRepositoryPostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, err := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusSeeOther, rr.Code)
	}

	// testing the parse form
	req, err = http.NewRequest("POST", "/make-reservation", nil)
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing invalid start-date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing invalid end-date
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing invalid roomId
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing form valid
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=n")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusSeeOther, rr.Code)
	}

	// testing insert reservation
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=nJoh")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing insert reservation
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=nJoh")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing insert reservation
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-04")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=nJoh")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Sule")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sule@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=20544 343 334")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}
	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation did not return appropriate response code, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepositoryPostAvailability(t *testing.T) {
	// test for when parsing form fails
	req, err := http.NewRequest(
		"POST", "/search-availability", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for invalid start date
	reqBody := "start=2025-15-21"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-10-10")
	req, err = http.NewRequest(
		"POST", "/search-availability", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for invalid end date
	reqBody = "start=2025-03-21"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-10-37")
	req, err = http.NewRequest(
		"POST", "/search-availability", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for when searching room fails
	reqBody = "start=2026-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2026-12-12")
	req, err = http.NewRequest(
		"POST", "/search-availability", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// test for when no room is available
	reqBody = "start=2024-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-10-10")
	req, err = http.NewRequest(
		"POST", "/search-availability", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusSeeOther, rr.Code)
	}

	// test for when a room is available
	reqBody = "start=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2025-12-12")
	req, err = http.NewRequest(
		"POST", "/search-availability", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post search availability returned wrong response, "+
			"expected %d but got %d", http.StatusOK, rr.Code)
	}

}

func TestRepositoryPostAvailabilityJson(t *testing.T) {
	// test parsing form
	req, err := http.NewRequest(
		"POST", "/search-availability-json", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailabilityJson)

	handler.ServeHTTP(rr, req)

	var resp jsonResponse
	err = json.Unmarshal([]byte(rr.Body.String()), &resp)
	if err != nil {
		t.Error(err)
	}

	if resp.OK {
		t.Errorf("invalid result for testing parse form")
	}

	// test failed db search
	reqBody := "start=2024-10-22"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-11-22")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	req, err = http.NewRequest(
		"POST", "/search-availability-json", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailabilityJson)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)

	if resp.OK {
		t.Errorf("invalid result for testing search, passed when should fail")
	}

	// test when room is available
	reqBody = "start=2024-10-22"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-11-22")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, err = http.NewRequest(
		"POST", "/search-availability-json", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailabilityJson)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)

	if !resp.OK {
		t.Errorf("invalid result for testing search, passed when should fail")
	}

	// test when room is not available
	reqBody = "start=2024-10-22"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2024-11-22")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, err = http.NewRequest(
		"POST", "/search-availability-json", strings.NewReader(reqBody))
	if err != nil {
		t.Error(err)
	}

	ctx = getctx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailabilityJson)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Error(err)
	}
	log.Println(resp)

	if resp.OK {
		t.Errorf("invalid result for testing search, found room avail when shouldn't")
	}
}

func TestRepositoryReservationSummary(t *testing.T) {
	//Testing instance where there is no reservation in session
	req, err := http.NewRequest("GET", "/reservation-summary", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getctx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation summary returned wrong response for no reservation, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// Testing case where there is a reservation
	session.Put(req.Context(), "reservation", models.Reservation{})
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation summary returned wrong response for no reservation, "+
			"expected %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestRepositoryChooseRoom(t *testing.T) {
	// Test no room id param
	req, err := http.NewRequest("GET", "/choose-room/1", nil)
	if err != nil {
		t.Error(err)
	}

	ctx := getctx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/"

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
	//	Testing no reservation in session
	req.RequestURI = "/choose-room/2"

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom returned wrong response, "+
			"expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// Test when reservation and id param are available

	session.Put(req.Context(), "reservation", models.Reservation{})
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom returned wrong response, "+
			"expected %d but got %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepositoryBookRoom(t *testing.T) {
	req, err := http.NewRequest("GET", "/book-room?id=3&s=2024-04-02&e=2024-05-02", nil)
	if err != nil {
		t.Error(err)
	}
	ctx := getctx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom returned wrong response, "+
			"expected %d but got %d", http.StatusSeeOther, rr.Code)
	}
}

func getctx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
