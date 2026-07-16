package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrarInterface(t *testing.T) {
	r := &mockRegistrar{}
	r.RegisterHandler("TestOp", nil)
	r.RegisterHandlerForService("svc", "TestOp", nil)
	assert.Equal(t, 2, len(r.handlers))
}

type mockRegistrar struct {
	handlers []string
}

func (m *mockRegistrar) RegisterHandler(operationName string, handler Handler) {
	m.handlers = append(m.handlers, operationName)
}

func (m *mockRegistrar) RegisterHandlerForService(serviceName, operationName string, handler Handler) {
	m.handlers = append(m.handlers, serviceName+"/"+operationName)
}
