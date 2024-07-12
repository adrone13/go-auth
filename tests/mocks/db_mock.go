package mocks

import "context"

type DatabaseMock struct{}

func (db *DatabaseMock) Ping(ctx context.Context) error {
	return nil
}
