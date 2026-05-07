package service

import (
	"errors"
	"fmt"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type AuthService struct {
	userRepo      repo.UserRepo
	workspaceRepo repo.WorkspaceRepo
}

func NewAuthService(userRepo repo.UserRepo, workspaceRepo repo.WorkspaceRepo) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Signup(req *SignupRequest) (*domain.User, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Name,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	ws := &domain.Workspace{
		Name:    fmt.Sprintf("%s's Workspace", req.Name),
		OwnerID: user.ID,
	}
	if err := s.workspaceRepo.Create(ws); err != nil {
		return nil, err
	}

	member := &domain.WorkspaceMember{
		WorkspaceID: ws.ID,
		UserID:      user.ID,
		Role:        "owner",
	}
	if err := s.workspaceRepo.AddMember(member); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GetCurrentUser(userID int64) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetUserWorkspace(userID int64) (*domain.Workspace, error) {
	return s.workspaceRepo.FindByUserID(userID)
}
