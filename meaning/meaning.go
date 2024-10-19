package meaning

type MeaningGet interface {
	Get(key string) (string, error)
}

type Meaning struct {
	dbHandler MeaningGet
}

type Settings struct {
	DBHandler MeaningGet
}

func New(s Settings) Meaning {
	return Meaning{
		dbHandler: s.DBHandler,
	}
}

func (m Meaning) GetMeaning(word string) string {
	meaning, err := m.dbHandler.Get(word)
	if err != nil {
		return ""
	}
	return meaning
}
