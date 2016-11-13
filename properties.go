package katolomb

import "fmt"

// TranslationProperties is the interface that wraps the Property method used by
// interpolators to obtain the values the text requires.
//
// Property takes a string (the name of the property) and returns a string with
// the value and an error.
type TranslationProperties interface {
	Property(string) (string, error)
}

// TranslationPropertiesFunc wraps a function with the TranslationProperties's
// Property method signature to satisfy the TranslationProperties interface.
type TranslationPropertiesFunc func(string) (string, error)

// NewTranslationProperties wraps a map of strings to strings and returns a
// TranslationProperties providing access to the elements of the map with its
// Property method
func NewTranslationProperties(ps map[string]string) TranslationProperties {
	return TranslationPropertiesFunc(func(p string) (string, error) {
		v, ok := ps[p]
		if !ok {
			return "", fmt.Errorf("property not available")
		}
		return v, nil
	})
}

// NewTranslationPropertiesWithDefault takes a TranslationProperties parameter
// and a defaultValue string and returns a TranslationProperties with an
// Interpolate method that returns the defaultValue when the wrapped
// TranslationProperties parameter has no value available.
func NewTranslationPropertiesWithDefault(ps TranslationProperties, defaultValue string) TranslationProperties {
	return TranslationPropertiesFunc(func(p string) (string, error) {
		v, err := ps.Property(p)
		if err != nil {
			v = defaultValue
		}
		return v, nil
	})
}

// Property calls the function with the received text and properties
// parameters and returns the result.
func (tpf TranslationPropertiesFunc) Property(p string) (string, error) {
	return tpf(p)
}
