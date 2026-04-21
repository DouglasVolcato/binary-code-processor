package test

import (
	"github.com/go-faker/faker/v4"
)

type FakeData struct{}

func (f *FakeData) ID() string {
	return faker.UUIDDigit()
}

func (f *FakeData) Phrase() string {
	return faker.Sentence()
}

func (f *FakeData) Binary() []byte {
	return []byte(faker.Paragraph())
}

func (f *FakeData) Date() string {
	return faker.Date()
}
