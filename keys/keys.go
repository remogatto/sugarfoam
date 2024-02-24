package keys

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/bubbles/key"
)

type Bindings struct {
	bindings  map[string]key.Binding
	shortHelp func() []key.Binding
	fullHelp  func() [][]key.Binding
}

type Option func(*Bindings)

func New(opts ...Option) *Bindings {
	kb := &Bindings{
		bindings: make(map[string]key.Binding),
	}
	for _, opt := range opts {
		opt(kb)
	}

	return kb

}

func KeyMapToMap(scope string, keyMap interface{}) (map[string]key.Binding, error) {
	result := make(map[string]key.Binding)

	val := reflect.ValueOf(keyMap)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("keyMap should be a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()
		result[fmt.Sprintf("%s.%s", scope, field.Name)] = value.(key.Binding)
	}

	return result, nil
}

func (kb *Bindings) ShortHelp() []key.Binding {
	return kb.shortHelp()
}

func (kb *Bindings) FullHelp() [][]key.Binding {
	return kb.fullHelp()
}

func (kb *Bindings) GetKey(name string) key.Binding {
	return kb.bindings[name]
}

func (kb *Bindings) SetKey(name string, key key.Binding) {
	kb.bindings[name] = key
}

func WithBinding(name string, binding key.Binding) Option {
	return func(kb *Bindings) {
		kb.bindings[name] = binding
	}
}

func WithShortHelp(keys ...string) Option {
	return func(kb *Bindings) {
		kb.shortHelp = func() []key.Binding {
			bindings := make([]key.Binding, 0)
			for _, key := range keys {
				bindings = append(bindings, kb.GetKey(key))
			}
			return bindings
		}
	}
}
