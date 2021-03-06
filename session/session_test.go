package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mechwarrior1/PGL_frontend/jwtsession"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// test for setting cookie and getting cookie
func TestNewCookie(t *testing.T) {

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e := echo.New()
	c := e.NewContext(req, rec)
	NewCookie(c, 300, "id")

	req = &http.Request{Header: http.Header{"Cookie": rec.HeaderMap["Set-Cookie"]}}
	c = e.NewContext(req, rec)
	// c.Request.Header.Set("Cookie", "user=gin")
	// cookie, _ := c.Cookie("user")
	// assert.Equal(t, "gin", cookie)

	goRecycleCookie, err := c.Cookie("goRecycleCookie")
	assert.NoError(t, err)
	assert.Equal(t, "id", goRecycleCookie.Value)

}

func TestGetCookieJwt(t *testing.T) {
	jwtWrapper := &jwtsession.JwtWrapper{
		"key",
		"GoRecycle",
		10,
	}
	sessionMgr := &Session{
		MapSession: &map[string]SessionStruct{"username": SessionStruct{"uuid", 123}},
		ApiKey:     "key",
	}

	generatedToken, _, err := jwtWrapper.GenerateToken("success", "msg", "false", "lastlogin", "username", "uuid")
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:   "goRecycleCookie",
		Value:  generatedToken,
		MaxAge: 300,
		Path:   "/",
	})

	e := echo.New()
	c := e.NewContext(req, rec)
	claims, err := sessionMgr.GetCookieJwt(c, jwtWrapper)
	assert.NoError(t, err)
	assert.Equal(t, "success", claims.Context.Success)
	assert.Equal(t, "msg", claims.Context.Msg)
	assert.Equal(t, "lastlogin", claims.Context.LastLogin)
	assert.Equal(t, "username", claims.Context.Username)
	assert.Equal(t, "uuid", claims.Context.Uuid)
	assert.Equal(t, "GoRecycle", claims.StandardClaims.Issuer)
}

func TestCheckSession(t *testing.T) {
	jwtWrapper := &jwtsession.JwtWrapper{
		"key",
		"GoRecycle",
		10,
	}
	sessionMgr := &Session{
		MapSession: &map[string]SessionStruct{"username": SessionStruct{"uuid", 123}},
		ApiKey:     "key",
	}

	_, claims, _ := jwtWrapper.GenerateToken("success", "msg", "false", "lastlogin", "username", "uuid")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	e := echo.New()
	c := e.NewContext(req, rec)

	claims = sessionMgr.CheckSession(c, claims, jwtWrapper)

	assert.Equal(t, "success", claims.Context.Success)
	assert.Equal(t, "msg", claims.Context.Msg)
	assert.Equal(t, "lastlogin", claims.Context.LastLogin)
	assert.Equal(t, "username", claims.Context.Username)
	assert.Equal(t, "uuid", claims.Context.Uuid)
	assert.Equal(t, "GoRecycle", claims.StandardClaims.Issuer)
}

// test if user will be logged out if their session is not mapped
func TestCheckSessionLogout(t *testing.T) {
	jwtWrapper := &jwtsession.JwtWrapper{
		"key",
		"GoRecycle",
		10,
	}
	sessionMgr := &Session{
		MapSession: &map[string]SessionStruct{},
		ApiKey:     "key",
	}

	_, claims, _ := jwtWrapper.GenerateToken("success", "msg", "false", "lastlogin", "username", "uuid")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	e := echo.New()
	c := e.NewContext(req, rec)

	claims = sessionMgr.CheckSession(c, claims, jwtWrapper)

	assert.Equal(t, "error", claims.Context.Success)
	assert.Equal(t, "you have been logged out", claims.Context.Msg)
}

//cchecks functions that updates jwt in cookie
func TestUpdateJwtLong(t *testing.T) {
	jwtWrapper := &jwtsession.JwtWrapper{
		"key",
		"GoRecycle",
		10,
	}
	sessionMgr := &Session{
		MapSession: &map[string]SessionStruct{},
		ApiKey:     "key",
	}

	_, claims, _ := jwtWrapper.GenerateToken("success", "msg", "false", "lastlogin", "username", "uuid")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	e := echo.New()
	c := e.NewContext(req, rec)

	UpdateJwtLong("123", "456", "7", "8", "me", &claims.Context, c, jwtWrapper, sessionMgr)

	req = &http.Request{Header: http.Header{"Cookie": rec.HeaderMap["Set-Cookie"]}}
	c = e.NewContext(req, rec)

	claims, err := sessionMgr.GetCookieJwt(c, jwtWrapper)

	assert.NoError(t, err)
	assert.Equal(t, "123", claims.Context.Success)
	assert.Equal(t, "456", claims.Context.Msg)
	assert.Equal(t, "8", claims.Context.LastLogin)
	assert.Equal(t, "uuid", (*(*sessionMgr).MapSession)["me"].Uuid) //check if session is mapped
}

func TestUpdateJwt(t *testing.T) {
	jwtWrapper := &jwtsession.JwtWrapper{
		"key",
		"GoRecycle",
		10,
	}
	sessionMgr := &Session{
		MapSession: &map[string]SessionStruct{"username": SessionStruct{"uuid", 123}},
		ApiKey:     "key",
	}

	_, claims, _ := jwtWrapper.GenerateToken("success", "msg", "false", "lastlogin", "username", "uuid")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	e := echo.New()
	c := e.NewContext(req, rec)

	UpdateJwt("123", "456", &claims.Context, c, jwtWrapper)

	req = &http.Request{Header: http.Header{"Cookie": rec.HeaderMap["Set-Cookie"]}}
	c = e.NewContext(req, rec)

	claims, err := sessionMgr.GetCookieJwt(c, jwtWrapper)

	assert.NoError(t, err)
	assert.Equal(t, "123", claims.Context.Success)
	assert.Equal(t, "456", claims.Context.Msg)
	assert.Equal(t, "lastlogin", claims.Context.LastLogin)
}
