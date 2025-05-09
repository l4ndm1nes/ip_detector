package service

import (
	"context"
	"fmt"

	"ip_detector/internal/auth"
	"ip_detector/internal/domain/model"
	"ip_detector/internal/domain/port"
	"ip_detector/internal/logger"
)

type UserService struct {
	repo   port.UserRepository
	geoIP  port.GeoIPService
	Config *Config
}

type Config struct {
	JWTSecret     string
	JWTExpiration string
}

func NewUserService(repo port.UserRepository, geoIP port.GeoIPService, cfg *Config) *UserService {
	return &UserService{
		repo:   repo,
		geoIP:  geoIP,
		Config: cfg,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	log := logger.Log.Sugar()
	log.Infow("create user called", "email", user.Email, "ip", user.IP)

	if user.IP == "" {
		log.Warn("user IP is empty")
		return fmt.Errorf("user IP is required")
	}

	country, err := s.geoIP.GetCountryByIP(user.IP)
	if err != nil {
		log.Errorw("geoIP lookup failed", "ip", user.IP, "error", err)
		return fmt.Errorf("failed to enrich user with country: %w", err)
	}
	user.Country = country

	if err := s.repo.Save(ctx, user); err != nil {
		log.Errorw("save user failed", "email", user.Email, "error", err)
		return fmt.Errorf("failed to save user: %w", err)
	}

	log.Infow("user saved", "id", user.ID, "email", user.Email, "country", user.Country)
	return nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	log := logger.Log.Sugar()
	log.Info("get all users")

	users, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Errorw("get all users failed", "error", err)
		return nil, err
	}

	log.Infow("users fetched", "count", len(users))
	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	log := logger.Log.Sugar()
	log.Infow("get user by id", "id", id)

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Errorw("get user by id failed", "id", id, "error", err)
		return nil, err
	}
	if u == nil {
		log.Infow("user not found", "id", id)
	}
	return u, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	log := logger.Log.Sugar()
	log.Infow("get user by email", "email", email)

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Errorw("get user by email failed", "email", email, "error", err)
		return nil, err
	}
	if u == nil {
		log.Infow("user not found", "email", email)
	}
	return u, nil
}

func (s *UserService) GenerateJWT(email string) (string, error) {
	log := logger.Log.Sugar()
	log.Infow("generating JWT", "email", email)

	token, err := auth.GenerateToken(email, s.Config.JWTSecret, s.Config.JWTExpiration)
	if err != nil {
		log.Errorw("generate JWT failed", "email", email, "error", err)
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return token, nil
}
