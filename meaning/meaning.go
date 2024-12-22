package meaning

import (
	"errors"
	"fmt"
)

type MeaningGetter interface {
	Get(key string) (string, error)
}

type Settings struct {
	GetMeaning MeaningGetter
}

type Meaning struct {
	getMeaning MeaningGetter
}

func check(s Settings) error {
	var aErr error

	if s.GetMeaning == nil {
		aErr = errors.Join(aErr, errors.New("no meaginig getter"))
	}

	return aErr
}

func New(settings Settings) (*Meaning, error) {
	err := check(settings)
	if err != nil {
		return nil, fmt.Errorf("checking settings: %w", err)
	}

	return &Meaning{
		getMeaning: settings.GetMeaning,
	}, nil
}

func (m Meaning) GetMeaning(word string) string {
	meaning, err := m.getMeaning.Get(word)
	if err != nil {
		return ""
	}

	return meaning
}
