package generators

import (
	"fmt"
	"sync"

	"github.com/schraf/assistant/pkg/models"
)

var (
	lock       sync.RWMutex
	generators map[string]Factory
)

type Config map[string]any

type Factory func(config Config) (models.ContentGenerator, error)

func Register(name string, factory Factory) error {
	lock.Lock()
	defer lock.Unlock()

	if generators == nil {
		generators = make(map[string]Factory)
	}

	if _, exists := generators[name]; exists {
		return fmt.Errorf("generator '%s' is already registered", name)
	}

	generators[name] = factory
	return nil
}

func MustRegister(name string, factory Factory) {
	if err := Register(name, factory); err != nil {
		panic(fmt.Sprintf("failed to register generator: %v", err))
	}
}

func Create(name string, config Config) (models.ContentGenerator, error) {
	lock.RLock()
	factory, exists := generators[name]
	lock.RUnlock()

	if !exists {
		return nil, fmt.Errorf("generator '%s' is not registered", name)
	}

	return factory(config)
}
