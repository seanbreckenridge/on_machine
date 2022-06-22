package on_machine

// keeps track of successes and failures separately, so we can memoize errors
type Memoizer struct {
	Storage map[string]interface{}
	Errors  map[string]error
}

func NewMemoizer() *Memoizer {
	return &Memoizer{
		Storage: make(map[string]interface{}),
		Errors:  make(map[string]error),
	}
}

func (m *Memoizer) Memoize(key string, fn func() (interface{}, error)) (interface{}, error, bool) {
	// if cached value succeeded, return value
	if cachedVal, ok := m.Storage[key]; ok {
		return cachedVal, nil, true
	}
	if errorVal, ok := m.Errors[key]; ok {
		return nil, errorVal, true
	}
	// call the function
	data, innerErr := fn()
	// succeeded, store and return
	if innerErr == nil {
		m.Storage[key] = data
		return data, innerErr, false
	} else {
		// failed, store error
		m.Errors[key] = innerErr
		return data, innerErr, false
	}
}

var Cache *Memoizer = NewMemoizer()
