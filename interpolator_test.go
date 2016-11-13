package katolomb_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/pbanos/katolomb"
)

func TestNewNoErrorInterpolator(t *testing.T) {
	errorInterpolator := katolomb.InterpolatorFunc(func(text string, ps katolomb.TranslationProperties) (string, error) {
		ps.Property("prop")
		return "", fmt.Errorf("property in %s not available", strconv.Quote(text))
	})
	workingInterpolator := katolomb.InterpolatorFunc(func(text string, ps katolomb.TranslationProperties) (string, error) {
		ps.Property("prop")
		return fmt.Sprintf("%s - some text", text), nil
	})
	testCases := []struct {
		baseInterpolator katolomb.Interpolator
		text             string
		result           string
		description      string
	}{
		{workingInterpolator, "this text", "this text - some text", "wrapping a succesful interpolation"},
		{errorInterpolator, "this other text", "this other text", "wrapping a failing interpolation"},
	}
	for _, tc := range testCases {
		i := katolomb.NewNoErrorInterpolator(tc.baseInterpolator)
		rightPropertiesPassed := false
		properties := katolomb.TranslationPropertiesFunc(func(string) (string, error) {
			rightPropertiesPassed = true
			return "", nil
		})
		result, err := i.Interpolate(tc.text, properties)
		if err != nil {
			t.Errorf("expected Interpolate not to return error when %v", tc.description)
		}
		if !rightPropertiesPassed {
			t.Errorf("expected Interpolate to pass the given TranslationProperties to the base interpolator")
		}
		if result != tc.result {
			t.Errorf("expected Interpolate to return %v when called with %v", strconv.Quote(tc.result), tc.description)
		}
	}
}

func TestNewInterpolator(t *testing.T) {
	profile := katolomb.NewTranslationProperties(map[string]string{
		"name":           "Alan Ginsberg",
		"hobby":          "writing poems",
		"favorite music": "jazz",
	})
	testCases := []struct {
		text        string
		result      string
		errNotNil   bool
		description string
	}{
		{"", "", false, "no content"},
		{"this text %{}", "this text %{}", false, "no interpolation needed"},
		{"my name is %{name}", "my name is Alan Ginsberg", false, "one default-value-less interpolation needed with property available"},
		{"my name is %{name|Frida Kahlo}", "my name is Alan Ginsberg", false, "one default-valued interpolation needed with property available"},
		{"my name is %{firstname}", "", true, "one default-value-less interpolation needed without property available"},
		{"my name is %{firstname|Frida}", "my name is Frida", false, "one default-valued interpolation needed without property available"},
		{"my name is %{name} and I like %{hobby}", "my name is Alan Ginsberg and I like writing poems", false, "two default-value-less interpolations needed with properties available"},
		{"my name is %{name|Frida Kahlo} and I like %{hobby|feminist activism}", "my name is Alan Ginsberg and I like writing poems", false, "two default-valued interpolations needed with property available"},
		{"my name is %{lastname|Ginsberg}, %{name}", "my name is Ginsberg, Alan Ginsberg", false, "one default-valued interpolation needed with default value and one default-value-less interpolation needed with property available"},
		{"my name is %{lastname}, %{name}", "", true, "two default-value-less interpolations needed and only one with property available"},
		{"my name is %{lastname|Ginsberg}, %{name}. My hobby is %{hobby|watching TV} and I like %{favorite music} music", "my name is Ginsberg, Alan Ginsberg. My hobby is writing poems and I like jazz music", false, "multiple interpolations with different scenarios"},
	}
	for _, tc := range testCases {
		i := katolomb.NewInterpolator()
		result, err := i.Interpolate(tc.text, profile)
		errNotNil := err != nil
		if errNotNil != tc.errNotNil {
			t.Errorf("expected Interpolate not to return error for a text with %v", tc.description)
		}
		if result != tc.result {
			t.Errorf("expected Interpolate to return %v when called for a text with %v", strconv.Quote(tc.result), tc.description)
		}
	}
}
