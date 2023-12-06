package render

import (
	"net/http"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/chenemiken/goland/bookings/internal/config"
)

var session *scs.SessionManager

var testApp config.AppConfig

func TestMain(m *testing.M) {

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp
}
