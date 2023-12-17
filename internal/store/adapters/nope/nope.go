package nope

type Nope struct {
}

func (n Nope) Validate(key string) bool {
	return true
}

func (n Nope) Remember(key string) {}
