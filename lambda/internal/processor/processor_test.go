package processor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rdsauditlogss3/internal/database"
	"rdsauditlogss3/internal/entity"
	"rdsauditlogss3/internal/logcollector"
	parser "rdsauditlogss3/internal/parser"
	"rdsauditlogss3/internal/s3writer"
	"io"
	"testing"
)

const (
	TestRdsInstanceIdentifier = "my-instance"
)

type mockDatabase struct {
	database.Database
	mock.Mock
}

func (m *mockDatabase) StoreCheckpoint(record *entity.CheckpointRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *mockDatabase) GetCheckpoint(id string) (*entity.CheckpointRecord, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.CheckpointRecord), args.Error(1)
}

type mockLogCollector struct {
	logcollector.LogCollector
	mock.Mock
}

func (m *mockLogCollector) GetLogs(timestamp int64) (io.Reader, bool, int64, error) {
	args := m.Called(timestamp)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool), args.Get(2).(int64), args.Error(3)
	}
	return args.Get(0).(io.Reader), args.Get(1).(bool), args.Get(2).(int64), args.Error(3)
}

func (m *mockLogCollector) ValidateAndPrepareRDSInstance() error {
	args := m.Called()
	return args.Error(0)
}

type mockWriter struct {
	s3writer.Writer
	mock.Mock
}

func (m *mockWriter) WriteLogEntry(data entity.LogEntry) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestProcessMultiLogCallback(t *testing.T) {
	p := parser.NewAuditLogParser()
	lc := new(mockLogCollector)
	w := new(mockWriter)

	processor := NewProcessor(lc, w, p, TestRdsInstanceIdentifier)
	err := processor.Process()
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, 2, 1)
}
