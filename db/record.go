//go:generate protoc -I. --go_out=. ./record.proto

package db

import (
	"time"

	"github.com/golang/protobuf/ptypes"
)

// GetTime returns the timestamp as a Go *time.Time.
func (m *Record) GetTime() *time.Time {
	timestamp := m.GetTimestamp()
	if timestamp == nil {
		return nil
	}
	time, err := ptypes.Timestamp(timestamp)
	if err != nil {
		panic(err)
	}
	return &time
}
