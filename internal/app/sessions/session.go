package sessions

import (
	"auth/internal/app/users"
	"auth/internal/config"
	"log"
	"time"
)

type SessionId string

type Session struct {
	Id            SessionId
	UserId        users.UserId
	RefreshTokens []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewSession(userId users.UserId) *Session {
	if userId == "" {
		log.Fatalln(`"userId" should not be empty`)
	}

	return &Session{
		UserId:        userId,
		RefreshTokens: []string{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (s *Session) AddRefreshToken(t string) {
	if t == "" {
		log.Fatalln("token shouldn't be empty")
	}

	s.RefreshTokens = append(s.RefreshTokens, t)
	s.UpdatedAt = time.Now()
}

func (s *Session) Expired() bool {
	absolute := config.Values.RefreshTokenAbsoluteTtl
	idle := config.Values.RefreshTokenIdleTtl

	now := time.Now()
	absoluteExpiration := s.CreatedAt.Add(time.Second * time.Duration(absolute))
	idleExpiration := s.UpdatedAt.Add(time.Second * time.Duration(idle))

	return absoluteExpiration.Before(now) || idleExpiration.Before(now)
}
