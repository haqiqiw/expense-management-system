package repository_test

import (
	"context"
	"errors"
	"expense-management-system/internal/entity"
	"expense-management-system/internal/repository"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/suite"
)

type UserRepositorySuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
	repo *repository.UserRepository
	ctx  context.Context
	now  time.Time
}

func (s *UserRepositorySuite) SetupTest() {
	s.mock, _ = pgxmock.NewPool()
	s.repo = repository.NewUserRepository(s.mock)
	s.ctx = context.Background()
	s.now = time.Now()
}

func (s *UserRepositorySuite) TearDownTest() {
	s.mock.Close()
}

func (s *UserRepositorySuite) TestUserRepository_Create() {
	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		param    *entity.User
		wantErr  error
	}{
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO users (email, name, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				)).
					WithArgs("john@mail.com", "John Doe", "password", "manager", pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint64(1)))
			},
			param: &entity.User{
				Email:        "john@mail.com",
				Name:         "John Doe",
				PasswordHash: "password",
				Role:         "manager",
			},
			wantErr: nil,
		},
		{
			name: "unexpected error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO users (email, name, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
				)).
					WithArgs("john@mail.com", "John Doe", "password", "manager", pgxmock.AnyArg()).
					WillReturnError(errors.New("something error"))
			},
			param: &entity.User{
				Email:        "john@mail.com",
				Name:         "John Doe",
				PasswordHash: "password",
				Role:         "manager",
			},
			wantErr: errors.New("something error"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			err := s.repo.Create(s.ctx, tt.param)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *UserRepositorySuite) TestUserRepository_FindByID() {
	tests := []struct {
		name     string
		mockFunc func(pgxmock.PgxPoolIface)
		paramID  uint64
		wantUser *entity.User
		wantErr  error
	}{
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "email", "name", "password_hash", "role", "created_at"}).
					AddRow(uint64(1), "john@mail.com", "John Doe", "password", "manager", s.now)
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnRows(rows)
			},
			paramID: uint64(1),
			wantUser: &entity.User{
				ID:           uint64(1),
				Email:        "john@mail.com",
				Name:         "John Doe",
				PasswordHash: "password",
				Role:         "manager",
				CreatedAt:    s.now,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(pgx.ErrNoRows)
			},
			paramID:  uint64(1),
			wantUser: nil,
			wantErr:  nil,
		},
		{
			name: "unexpected error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE id = $1 LIMIT 1`,
				)).
					WithArgs(uint64(1)).
					WillReturnError(errors.New("something error"))
			},
			paramID:  uint64(1),
			wantUser: nil,
			wantErr:  errors.New("something error"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.FindByID(s.ctx, tt.paramID)
			s.Equal(tt.wantUser, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *UserRepositorySuite) TestUserRepository_FindByEmail() {
	tests := []struct {
		name       string
		mockFunc   func(pgxmock.PgxPoolIface)
		paramEmail string
		wantUser   *entity.User
		wantErr    error
	}{
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "email", "name", "password_hash", "role", "created_at"}).
					AddRow(uint64(1), "john@mail.com", "John Doe", "password", "manager", s.now)
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE email = $1 LIMIT 1`,
				)).
					WithArgs("john@mail.com").
					WillReturnRows(rows)
			},
			paramEmail: "john@mail.com",
			wantUser: &entity.User{
				ID:           uint64(1),
				Email:        "john@mail.com",
				Name:         "John Doe",
				PasswordHash: "password",
				Role:         "manager",
				CreatedAt:    s.now,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE email = $1 LIMIT 1`,
				)).
					WithArgs("john@mail.com").
					WillReturnError(pgx.ErrNoRows)
			},
			paramEmail: "john@mail.com",
			wantUser:   nil,
			wantErr:    nil,
		},
		{
			name: "unexpected error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(
					`SELECT id, email, name, password_hash, role, created_at FROM users WHERE email = $1 LIMIT 1`,
				)).
					WithArgs("john@mail.com").
					WillReturnError(errors.New("something error"))
			},
			paramEmail: "john@mail.com",
			wantUser:   nil,
			wantErr:    errors.New("something error"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.FindByEmail(s.ctx, tt.paramEmail)
			s.Equal(tt.wantUser, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func (s *UserRepositorySuite) TestUserRepository_CountByEmail() {
	tests := []struct {
		name       string
		mockFunc   func(pgxmock.PgxPoolIface)
		paramEmail string
		wantTotal  int
		wantErr    error
	}{
		{
			name: "success",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM users WHERE email = $1`)).
					WithArgs("john@mail.com").
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
			},
			paramEmail: "john@mail.com",
			wantTotal:  1,
			wantErr:    nil,
		},
		{
			name: "not found",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM users WHERE email = $1`)).
					WithArgs("john@mail.com").
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))
			},
			paramEmail: "john@mail.com",
			wantTotal:  0,
			wantErr:    nil,
		},
		{
			name: "unexpected error",
			mockFunc: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM users WHERE email = $1`)).
					WithArgs("john@mail.com").
					WillReturnError(errors.New("something error"))
			},
			paramEmail: "john@mail.com",
			wantTotal:  0,
			wantErr:    errors.New("something error"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mockFunc(s.mock)

			res, err := s.repo.CountByEmail(s.ctx, tt.paramEmail)
			s.Equal(tt.wantTotal, res)
			s.Equal(tt.wantErr, err)
		})
	}
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
