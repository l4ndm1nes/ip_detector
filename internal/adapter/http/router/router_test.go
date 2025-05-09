package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"ip_detector/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"ip_detector/internal/adapter/http/router"
	"ip_detector/internal/app/service"
	"ip_detector/internal/domain/model"
	"ip_detector/internal/domain/port"
)

type mockRepo struct {
	users map[string]*model.User
}

func newMockRepo() *mockRepo { return &mockRepo{users: map[string]*model.User{}} }

func (m *mockRepo) Save(_ context.Context, u *model.User) error {
	m.users[u.Email] = u
	return nil
}
func (m *mockRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	return m.users[email], nil
}
func (m *mockRepo) GetByID(_ context.Context, id string) (*model.User, error) { // not used
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}
func (m *mockRepo) GetAll(_ context.Context) ([]*model.User, error) {
	var out []*model.User
	for _, u := range m.users {
		out = append(out, u)
	}
	return out, nil
}

var _ port.UserRepository = (*mockRepo)(nil)

type geoIPMock struct{}

func (g geoIPMock) GetCountryByIP(_ string) (string, error) { return "UA", nil }

func setupTestRouter() http.Handler {
	logger.Init()

	repo := newMockRepo()
	geo := geoIPMock{}
	cfg := &service.Config{JWTSecret: "supersecretkey", JWTExpiration: "24h"}
	us := service.NewUserService(repo, geo, cfg)
	return router.SetupRouter(us, "supersecretkey")
}

func TestRegisterOK(t *testing.T) {
	r := setupTestRouter()
	rec := httptest.NewRecorder()

	body := `{"name":"Alice","email":"alice@example.com","ip":"8.8.8.8","password":"secret123"}`
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRegisterValidationFail(t *testing.T) {
	r := setupTestRouter()
	rec := httptest.NewRecorder()

	body := `{"name":"","email":"bad","ip":"not_ip","password":"1"}`
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

func TestLoginAndProtected(t *testing.T) {
	r := setupTestRouter()

	reg := httptest.NewRecorder()
	regBody := `{"name":"Bob","email":"bob@example.com","ip":"1.1.1.1","password":"hunter2"}`
	reqReg, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(regBody))
	reqReg.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(reg, reqReg)
	if reg.Code != http.StatusCreated {
		t.Fatalf("register failed: %d", reg.Code)
	}

	login := httptest.NewRecorder()
	loginBody := `{"email":"bob@example.com","password":"hunter2"}`
	reqLogin, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(login, reqLogin)
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d", login.Code)
	}

	var resp struct{ Token string }
	if err := json.Unmarshal(login.Body.Bytes(), &resp); err != nil || resp.Token == "" {
		t.Fatalf("cannot parse token: %v, body: %s", err, login.Body.String())
	}

	noTok := httptest.NewRecorder()
	reqNoTok, _ := http.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(noTok, reqNoTok)
	if noTok.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 without token, got %d", noTok.Code)
	}

	withTok := httptest.NewRecorder()
	reqTok, _ := http.NewRequest(http.MethodGet, "/users", nil)
	reqTok.Header.Set("Authorization", "Bearer "+resp.Token)
	r.ServeHTTP(withTok, reqTok)
	if withTok.Code != http.StatusOK {
		t.Fatalf("want 200 with token, got %d", withTok.Code)
	}
}
