package entities

import (
	"testing"

	"github.com/douglasvolcato/task-monitor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	ID         string
	Message    string
	BinaryData []byte
	CreatedAt  string
	UpdatedAt  string
}

func MakeTestData() testData {
	return testData{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryData: faker.Binary(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

var faker = test.FakeData{}
var fakeData = MakeTestData()

func TestNewTaskShouldCreateTask(t *testing.T) {
	sut := NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	)
	assert.NotNil(t, sut)
	assert.Equal(t, fakeData.ID, sut.ID)
	assert.Equal(t, fakeData.Message, sut.Message)
	assert.Equal(t, fakeData.BinaryData, sut.BinaryData)
	assert.Equal(t, fakeData.CreatedAt, sut.CreatedAt)
	assert.Equal(t, fakeData.UpdatedAt, sut.UpdatedAt)
}

func TestValidateShouldReturnErrorIfDataIsInvalid(t *testing.T) {
	err := NewTask(
		"",
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		"",
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		([]byte)(""),
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		"",
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		"",
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		" ",
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		" ",
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		([]byte)(" "),
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		" ",
		fakeData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		" ",
	).Validate()
	assert.Error(t, err)
}

func TestValidateShouldReturnNilIfDataIsValid(t *testing.T) {
	err := NewTask(
		fakeData.ID,
		fakeData.Message,
		fakeData.BinaryData,
		fakeData.CreatedAt,
		fakeData.UpdatedAt,
	).Validate()
	assert.NoError(t, err)
}
