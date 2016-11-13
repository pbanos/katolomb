package katolomb_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/pbanos/katolomb"
)

func TestNewDefaultTranslator(t *testing.T) {
	errTranslator := katolomb.TranslatorFunc(func(k string, p katolomb.TranslationProperties) (string, error) {
		p.Property("prop")
		return "", fmt.Errorf("some error")
	})
	successfulTranslator := katolomb.TranslatorFunc(func(key string, p katolomb.TranslationProperties) (string, error) {
		p.Property("prop")
		return fmt.Sprintf("translated %v", key), nil
	})
	key := "my translation key"
	testCases := []struct {
		translator         katolomb.Translator
		defaultTranslation string
		result             string
		description        string
	}{
		{errTranslator, "lost in translation", "lost in translation", "wrapped translator cannot translate the key"},
		{errTranslator, "my default translation", "my default translation", "wrapped translator cannot translate the key"},
		{successfulTranslator, "my default translation", "translated my translation key", "wrapped translator can translate the key"},
		{successfulTranslator, "lost in translation", "translated my translation key", "wrapped translator can translate the key"},
	}
	for _, tc := range testCases {
		translator := katolomb.NewDefaultTranslator(tc.defaultTranslation, tc.translator)
		rightPropertiesPassed := false
		properties := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
			rightPropertiesPassed = true
			return "", nil
		})
		result, err := translator.Translate(key, properties)
		if err != nil {
			t.Errorf("expected Translate not to return error when %v", tc.description)
		}
		if !rightPropertiesPassed {
			t.Errorf("expected Translate to pass the given TranslationProperties to the base interpolator when %v", tc.description)
		}
		if result != tc.result {
			t.Errorf("expected Translate to return %v when %v", strconv.Quote(tc.result), tc.description)
		}
	}
}

func TestNewKeyAsDefaultTranslator(t *testing.T) {
	errTranslator := katolomb.TranslatorFunc(func(k string, p katolomb.TranslationProperties) (string, error) {
		p.Property("prop")
		return "", fmt.Errorf("some error")
	})
	successfulTranslator := katolomb.TranslatorFunc(func(key string, p katolomb.TranslationProperties) (string, error) {
		p.Property("prop")
		return fmt.Sprintf("translated %v", key), nil
	})
	testCases := []struct {
		translator  katolomb.Translator
		key         string
		result      string
		description string
	}{
		{errTranslator, "my.key", "my.key", "wrapped translator cannot translate the key"},
		{errTranslator, "some other translation key", "some other translation key", "wrapped translator cannot translate the key"},
		{successfulTranslator, "my.key", "translated my.key", "wrapped translator can translate the key"},
		{successfulTranslator, "some other translation key", "translated some other translation key", "wrapped translator can translate the key"},
	}
	for _, tc := range testCases {
		translator := katolomb.NewKeyAsDefaultTranslator(tc.translator)
		rightPropertiesPassed := false
		properties := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
			rightPropertiesPassed = true
			return "", nil
		})
		result, err := translator.Translate(tc.key, properties)
		if err != nil {
			t.Errorf("expected Translate not to return error when %v", tc.description)
		}
		if !rightPropertiesPassed {
			t.Errorf("expected Translate to pass the given TranslationProperties to the base interpolator when %v", tc.description)
		}
		if result != tc.result {
			t.Errorf("expected Translate to return %v when %v", strconv.Quote(tc.result), tc.description)
		}
	}
}

func TestNewInterpolatedTranslator(t *testing.T) {
	errTranslator := katolomb.TranslatorFunc(func(k string, p katolomb.TranslationProperties) (string, error) {
		p.Property("translator")
		return "", fmt.Errorf("some error")
	})
	successfulTranslator := katolomb.TranslatorFunc(func(key string, p katolomb.TranslationProperties) (string, error) {
		p.Property("translator")
		return fmt.Sprintf("translated %v", key), nil
	})
	errInterpolator := katolomb.InterpolatorFunc(func(k string, p katolomb.TranslationProperties) (string, error) {
		p.Property("interpolator")
		return "", fmt.Errorf("some error")
	})
	successfulInterpolator := katolomb.InterpolatorFunc(func(k string, p katolomb.TranslationProperties) (string, error) {
		p.Property("interpolator")
		return fmt.Sprintf("interpolated %v", k), nil
	})
	testCases := []struct {
		translator             katolomb.Translator
		interpolator           katolomb.Interpolator
		key                    string
		result                 string
		errNotNil              bool
		propertiesInterpolated bool
		description            string
	}{
		{errTranslator, errInterpolator, "my.key", "", true, false, "translator and interpolator return errors"},
		{errTranslator, successfulInterpolator, "my.key", "", true, false, "translator return error but interpolator does not"},
		{successfulTranslator, errInterpolator, "my.key", "", true, true, " translator is successful but interpolator returns error"},
		{successfulTranslator, successfulInterpolator, "my.key", "interpolated translated my.key", false, true, " translator and interpolator are successful"},
	}
	for _, tc := range testCases {
		translator := katolomb.NewInterpolatedTranslator(tc.translator, tc.interpolator)
		propertiesForTranslator := false
		propertiesForInterpolator := false
		properties := katolomb.TranslationPropertiesFunc(func(property string) (string, error) {
			switch property {
			case "translator":
				propertiesForTranslator = true
			case "interpolator":
				propertiesForInterpolator = true
			}
			return "", nil
		})
		result, err := translator.Translate(tc.key, properties)
		errNotNil := err != nil
		if errNotNil != tc.errNotNil {
			if errNotNil {
				t.Errorf("expected Translate not to return error when %v", tc.description)
			} else {
				t.Errorf("expected Translate to return error when %v", tc.description)
			}
		}
		if !propertiesForTranslator {
			t.Errorf("expected Translate to pass the given TranslationProperties to the wrapped Translator when %v", tc.description)
		}
		if propertiesForInterpolator != tc.propertiesInterpolated {
			t.Errorf("expected Translate to pass the given TranslationProperties to the Interpolator when %v", tc.description)
		}
		if result != tc.result {
			t.Errorf("expected Translate to return %v when %v", strconv.Quote(tc.result), tc.description)
		}
	}
}
