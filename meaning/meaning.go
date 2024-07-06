package meaning

type Meaning struct{}

func (m Meaning) GetMeaning(word string) string {
	return "test"
}

func New() Meaning {
	return Meaning{}
}
