package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	// {"search-avail-post", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-03"},
	// }, http.StatusOK},
	// {"search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-03"},
	// }, http.StatusOK},
	// {"search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "Smith@john.com"},
	// 	{key: "phone", value: "444-444-4444"},
	// }, http.StatusOK},
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

func getctx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
