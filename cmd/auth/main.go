package main

import (
	"auth/internal/app"
	_ "auth/internal/config"
	"auth/internal/db"
	"auth/internal/logger"
	"context"
	"log"
	"time"
)

func main() {

	sr := db.SessionRepository{}

	now := time.Now()
	session := &app.Session{
		UserId:        "aac787b1-2ce5-41c0-9e62-54129068454a",
		RefreshTokens: []string{"refresh_token:1"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	err := sr.Insert(context.Background(), session)
	if err != nil {
		log.Fatalf("Failed to insert. Err: %s", err)
	}

	s, err := sr.FindByUserId(context.Background(), session.UserId)
	if err != nil {
		log.Fatalf("Failed to find by user id. Err: %s", err)
	}

	logger.Debug(s)

	time.Sleep(1 * time.Second)

	session.AddRefreshToken("refresh_token:2")

	err = sr.Update(context.Background(), session)
	if err != nil {
		log.Fatalf("Failed to update session. Err: %s\n", err)
	}

	//err = sr.Delete(context.Background(), session.Id)
	//if err != nil {
	//	log.Fatalf("Failed to delete. Err: %s", err)
	//}

	//s := server.New(config.Values.Port, database)
	//
	//logger.Info(
	//	fmt.Sprintf("Server listening on 0.0.0.0%s", s.Addr),
	//)
	//
	//err := s.ListenAndServe()
	//if err != nil {
	//	fmt.Println("Failed to start")
	//
	//	panic(err)
	//}
}
