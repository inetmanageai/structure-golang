package logs

import (
	"github.com/stretchr/testify/mock"
)

type appLogsMock struct {
	mock.Mock
}

func NewAppLogsMock() *appLogsMock {
	return &appLogsMock{}
}

func (m *appLogsMock) Info(msg string) {
}

func (m *appLogsMock) Debug(msg string) {

}

func (m *appLogsMock) Warning(msg string) {

}

func (m *appLogsMock) Error(msg interface{}) {

}
