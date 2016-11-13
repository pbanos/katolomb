package katolomb_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/pbanos/katolomb"
)

func TestNewTranslationProperties(t *testing.T) {
	testCases := []struct {
		m           map[string]string
		p           string
		v           string
		eNotNil     bool
		description string
	}{
		{map[string]string{}, "prop", "", true, "empty map"},
		{map[string]string{"prop": "my value"}, "prop", "my value", false, "single key-value pair map with the property"},
		{map[string]string{"prop": "my value"}, "other prop", "", true, "single key-value pair map without the property"},
		{map[string]string{"prop": "my value", "my other prop": "other value"}, "my other prop", "other value", false, "multiple key-value pair map with the property last"},
		{map[string]string{"my other prop": "other value", "prop": "my value"}, "my other prop", "other value", false, "multiple key-value pair map with the property first"},
		{map[string]string{"my other prop": "other value", "prop": "my value", "yet another prop": "yet another value"}, "prop", "my value", false, "multiple key-value pair map with the property in the middle"},
		{map[string]string{"my other prop": "other value", "prop": "my value"}, "yet another prop", "", true, "multiple key-value pair map without the property"},
	}
	for _, tc := range testCases {
		tp := katolomb.NewTranslationProperties(tc.m)
		v, err := tp.Property(tc.p)
		eNotNil := err != nil
		if eNotNil != tc.eNotNil {
			if tc.eNotNil {
				t.Errorf("expected Property to return error when called with %v", tc.description)
			} else {
				t.Errorf("expected Property not to return error when called with %v", tc.description)
			}
		}
		if v != tc.v {
			t.Errorf("expected Property to return %v when called with %v", strconv.Quote(tc.v), tc.description)
		}
	}
}

func TestNewTranslationPropertiesWithDefault(t *testing.T) {
	defVal := "default value"
	valuePropsVal := "value"
	errorProps := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
		return "", fmt.Errorf("some error")
	})
	valueProps := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
		return valuePropsVal, nil
	})
	emptyValueProps := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
		return "", nil
	})
	testCases := []struct {
		tp          katolomb.TranslationProperties
		v           string
		description string
	}{
		{valueProps, valuePropsVal, "wrapped TranslationProperties provides a non-empty value"},
		{emptyValueProps, "", "wrapped TranslationProperties provides an empty value"},
		{errorProps, defVal, "wrapped TranslationProperties returns an error"},
	}
	for _, tc := range testCases {
		tp := katolomb.NewTranslationPropertiesWithDefault(tc.tp, defVal)
		v, err := tp.Property("property")
		if err != nil {
			t.Errorf("expected Property not to return error when %v", tc.description)
		}
		if v != tc.v {
			t.Errorf("expected Property to return %v when called with %v", strconv.Quote(tc.v), tc.description)
		}
	}
}
