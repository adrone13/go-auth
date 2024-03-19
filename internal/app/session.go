package app

import (
	"auth/internal/config"
	"log"
	"time"
)

type SessionId string

type Session struct {
	Id            SessionId
	UserId        UserId
	RefreshTokens []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	// https://stackoverflow.com/questions/24564619/nullable-time-time
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

	absoluteExpire := s.CreatedAt.Add(time.Second * time.Duration(absolute))
	if absoluteExpire.Before(now) {
		return false
	}

	idleExpire := s.UpdatedAt.Add(time.Second * time.Duration(idle))

	return idleExpire.Before(now)
}
