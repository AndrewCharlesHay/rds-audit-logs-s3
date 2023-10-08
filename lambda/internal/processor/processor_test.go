package processor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rdsauditlogss3/internal/database"
	"rdsauditlogss3/internal/logcollector"
	parser "rdsauditlogss3/internal/parser"
	"rdsauditlogss3/internal/s3writer"
	"testing"
)

const (
	TestRdsInstanceIdentifier = "my-instance"
)

type mockLogCollector struct {
	logcollector.LogCollector
	mock.Mock
}

type mockDatabase struct {
	database.Database
	mock.Mock
}

type mockWriter struct {
	s3writer.Writer
	mock.Mock
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
