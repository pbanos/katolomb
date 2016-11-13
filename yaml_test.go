package katolomb_test

import (
	"strconv"
	"testing"

	"github.com/pbanos/katolomb"
)

func TestNewYAMLTranslator(t *testing.T) {
	empty := ""
	emptyYAML := "---"
	emptyJSON := "{}"
	invalidYAML := "asaa"
	invalidJSON := `{"hello": "Hola"`
	emptyYAMLObject := "--- {}"
	shallowSingleKeyJSON := `{"hello": "Hola"}`
	shallowSingleKeyYAML := `---
hello: Hello!`
	shallowMultiKeyJSON := `{"hello": "Hola", "bye": "Adios"}`
	shallowMultiKeyYAML := `---
hello: Hello!
bye: Good bye`
	nestedJSON := `{"greetings": { "hello":"Hola", "bye":{ "night": "Buenas noches", "afternoon":"Buenas tardes"}}}`
	nestedYAML := `---
greetings:
  hello: Hello!
  bye:
    night: Good night!
    afternoon: Good afternoon!`
	jsonWithArrays := `{"greetings": { "hello":"Hola", "bye":{ "night": "Buenas noches", "afternoon":"Buenas tardes"}},
 "numbers":["cero", "uno", "dos", "tres"]}`
	yamlWithArrays := `---
greetings:
  hello: Hello!
  bye:
    night: Good night!
    afternoon: Good afternoon!
numbers:
- zero
- one
- two
- three`

	testCases := []struct {
		yaml                 []byte
		translatorErrNotNil  bool
		key                  string
		result               string
		translationErrNotNil bool
		description          string
	}{
		{[]byte(empty), false, "my.key", "", true, "building a translator with an empty string"},
		{[]byte(emptyYAML), false, "my.key", "", true, "building a translator with just the YAML header"},
		{[]byte(emptyJSON), false, "my.key", "", true, "building a translator with an empty json"},
		{[]byte(invalidYAML), true, "my.key", "", false, "building a translator with invalid YAML"},
		{[]byte(invalidJSON), true, "my.key", "", false, "building a translator with invalid JSON"},
		{[]byte(emptyYAMLObject), false, "my.key", "", true, "building a translator with an empty YAML"},
		{[]byte(shallowSingleKeyJSON), false, "my.key", "", true, "building a translator with an shallow single-key json and translating with another key"},
		{[]byte(shallowSingleKeyJSON), false, "hello", "Hola", false, "building a translator with an shallow single-key json and translating with the key"},
		{[]byte(shallowSingleKeyYAML), false, "my.key", "", true, "building a translator with an shallow single-key yaml and translating with another key"},
		{[]byte(shallowSingleKeyYAML), false, "hello", "Hello!", false, "building a translator with an shallow single-key yaml and translating with the key"},
		{[]byte(shallowMultiKeyJSON), false, "my.key", "", true, "building a translator with an shallow multi-key json and translating with a non-contained key"},
		{[]byte(shallowMultiKeyJSON), false, "hello", "Hola", false, "building a translator with an shallow single-key json and translating with a contained key"},
		{[]byte(shallowMultiKeyYAML), false, "hello.my.key", "", true, "building a translator with an shallow single-key yaml and translating with a non-contained key"},
		{[]byte(shallowMultiKeyYAML), false, "bye", "Good bye", false, "building a translator with an shallow single-key yaml and translating with a contained key"},
		{[]byte(nestedJSON), false, "greeting.bye.morning", "", true, "building a translator with a nested json and translating with a non-contained key"},
		{[]byte(nestedJSON), false, "greetings.bye.night", "Buenas noches", false, "building a translator with a nested json and translating with a contained nested key"},
		{[]byte(nestedYAML), false, "greetings", "", true, "building a translator with a nested yaml and translating with a non-contained key"},
		{[]byte(nestedYAML), false, "greetings.hello", "Hello!", false, "building a translator with a nested yaml and translating with a contained nested key"},
		{[]byte(jsonWithArrays), false, "numbers.4", "", true, "building a translator with json with arrays and translating with a non-contained key"},
		{[]byte(jsonWithArrays), false, "numbers.2", "dos", false, "building a translator with a json with arrays and translating with a contained nested key"},
		{[]byte(yamlWithArrays), false, "numbers.4", "", true, "building a translator with a yaml with arrays and translating with a non-contained key"},
		{[]byte(yamlWithArrays), false, "numbers.0", "zero", false, "building a translator with a yaml with arrays and translating with a contained nested key"},
	}
	for _, tc := range testCases {
		translator, err := katolomb.NewYAMLTranslator(tc.yaml)
		errNotNil := err != nil
		if errNotNil != tc.translatorErrNotNil {
			if errNotNil {
				t.Errorf("expected NewYAMLTranslator not to return error when %v", tc.description)
			} else {
				t.Errorf("expected NewYAMLTranslator to return error when %v", tc.description)
			}
		}
		if err == nil {
			properties := katolomb.TranslationPropertiesFunc(func(property string) (string, error) {
				return "", nil
			})
			result, err := translator.Translate(tc.key, properties)
			errNotNil := err != nil
			if errNotNil != tc.translationErrNotNil {
				if errNotNil {
					t.Errorf("expected Translate not to return error when %v", tc.description)
				} else {
					t.Errorf("expected Translate to return error when %v", tc.description)
				}
			}
			if result != tc.result {
				t.Errorf("expected Translate to return %v when %v", strconv.Quote(tc.result), tc.description)
			}
		}
	}
}
