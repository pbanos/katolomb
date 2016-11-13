package katolomb

import (
	"fmt"
	"strconv"
)

// Translator is the interface that wraps the basic Translate method.
//
// Translate takes a string and a TranslationProperties and returns a string
// with a translation and an error.
type Translator interface {
	Translate(string, TranslationProperties) (string, error)
}

// Translatable is the interface offered by entities that can be translated into
// strings using a Translator.
//
// Translate takes a Translator and returns a string with the translation and an
// error.
type Translatable interface {
	Translate(Translator) (string, error)
}

// TranslatorFunc wraps a function with the Translator's
// Translate method signature to satisfy the Translator interface.
type TranslatorFunc func(string, TranslationProperties) (string, error)

// NewDefaultTranslator takes a default translation string and a Translator and
// returns a new Translator that wraps the Translator parameter to return the
// default translation when its Translate method returns an error.
func NewDefaultTranslator(translation string, translator Translator) Translator {
	return TranslatorFunc(func(key string, props TranslationProperties) (string, error) {
		t, err := translator.Translate(key, props)
		if err != nil {
			return translation, nil
		}
		return t, nil
	})
}

// NewKeyAsDefaultTranslator takes a Translator and returns a new Translator
// that wraps the Translator parameter to return the Translate method's key
// parameter the Translator parameter's Translate method returns an error.
func NewKeyAsDefaultTranslator(translator Translator) Translator {
	return TranslatorFunc(func(key string, props TranslationProperties) (string, error) {
		t, err := translator.Translate(key, props)
		if err != nil {
			return key, nil
		}
		return t, nil
	})
}

// NewInterpolatedTranslator takes a Translator and an Interpolator parameters
// and returns a new Translator whose Translate method obtains the translation
// provided by the Translator parameter's Translate method, interpolates it
// using the Interpolator parameter and returns the result. If the translation
// or interpolation returns an error, an error is returned right away.
func NewInterpolatedTranslator(translator Translator, interpolator Interpolator) Translator {
	return TranslatorFunc(func(key string, props TranslationProperties) (string, error) {
		t, err := translator.Translate(key, props)
		if err != nil {
			return "", err
		}
		t, err = interpolator.Interpolate(t, props)
		if err != nil {
			return "", fmt.Errorf("translating %v: %v", strconv.Quote(key), err)
		}
		return t, nil
	})
}

// Translate calls the function with the received text and properties
// parameters and returns the result.
func (tf TranslatorFunc) Translate(key string, properties TranslationProperties) (string, error) {
	return tf(key, properties)
}
