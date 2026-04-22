package entities

import (
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	"github.com/stretchr/testify/assert"
)

type testDataTaskEntity struct {
	ID         string
	Message    string
	BinaryCode string
	CreatedAt  string
	UpdatedAt  string
}

func makeFakeDataTaskEntity() *testDataTaskEntity {
	faker := test.FakeData{}
	return &testDataTaskEntity{
		ID:         faker.ID(),
		Message:    faker.Phrase(),
		BinaryCode: faker.Phrase(),
		CreatedAt:  faker.Date(),
		UpdatedAt:  faker.Date(),
	}
}

func TestNewTaskShouldCreateTask(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	sut := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	)
	assert.NotNil(t, sut)
	assert.Equal(t, testData.ID, sut.ID)
	assert.Equal(t, testData.Message, sut.Message)
	assert.Equal(t, testData.BinaryCode, sut.BinaryCode)
	assert.Equal(t, testData.CreatedAt, sut.CreatedAt)
	assert.Equal(t, testData.UpdatedAt, sut.UpdatedAt)
}

func TestValidateShouldReturnErrorIfTaskDataIsInvalid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	err := NewTask(
		"",
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		"",
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		"",
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		"",
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		"",
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		" ",
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		" ",
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		"",
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		" ",
		testData.UpdatedAt,
	).Validate()
	assert.Error(t, err)
	err = NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		" ",
	).Validate()
	assert.Error(t, err)
}

func TestValidateShouldReturnNilIfTaskDataIsValid(t *testing.T) {
	testData := makeFakeDataTaskEntity()
	err := NewTask(
		testData.ID,
		testData.Message,
		testData.BinaryCode,
		testData.CreatedAt,
		testData.UpdatedAt,
	).Validate()
	assert.NoError(t, err)
}
