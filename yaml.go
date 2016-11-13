package katolomb

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type yamlTranslator struct {
	separator    string
	translations yamlTranslations
}

type yamlTranslations map[string]interface{}

// NewYAMLTranslator returns a Translator that looks for translations in the
// YAML passed as a  byte-slice parameter.
//
// The translation's YAML will be deserialized and its keys treated as strings.
// Lists in the YAML will be treated as maps with string-formatted integers as
// keys
//
// The result's Translate method will use the "." string as separator to split
// the key parameter into a tree route to a translation.
func NewYAMLTranslator(yml []byte) (Translator, error) {
	return NewYAMLTranslatorWithSeparator(yml, ".")
}

// NewYAMLTranslatorWithSeparator returns a Translator that looks for
// translations in the YAML in the byte slice passed as parameter.
//
// The translation's YAML will be deserialized and its keys treated as strings.
// Lists in the YAML will be treated as maps with string-formatted integers as
// keys
//
// The result's Translate method will use the given separator string to split
// the key parameter into a tree route to a translation.
func NewYAMLTranslatorWithSeparator(yml []byte, separator string) (Translator, error) {
	ts := make(map[interface{}]interface{})
	err := yaml.Unmarshal(yml, ts)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling yaml translations: %v", err)
	}
	yts := yamlTranslationizeMap(ts)
	yt := &yamlTranslator{
		separator:    separator,
		translations: yts,
	}
	return yt, nil
}

func (t *yamlTranslator) Translate(key string, properties TranslationProperties) (string, error) {
	var path []string
	if t.separator == "" {
		path = append(path, key)
	} else {
		path = strings.Split(key, t.separator)
	}
	translation, err := t.translations.find(path)
	if err != nil {
		return translation, fmt.Errorf("translating %v: %v", strconv.Quote(key), err)
	}
	return translation, nil
}

func (yts yamlTranslations) find(path []string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("incomplete path")
	}
	k := path[0]
	path = path[1:]
	v, ok := yts[k]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	switch v := v.(type) {
	case yamlTranslations:
		return v.find(path)
	case string:
		if len(path) != 0 {
			return "", fmt.Errorf("not found")
		}
		return v, nil
	default:
		panic(fmt.Sprintf("unexpected value of type %T in yamlTranslations", v))
	}
}

func yamlTranslationizeMap(ts map[interface{}]interface{}) yamlTranslations {
	yts := make(yamlTranslations)
	for k, v := range ts {
		yts[fmt.Sprintf("%v", k)] = yamlTranslationizeValue(v)
	}
	return yts
}

func yamlTranslationizeValue(v interface{}) interface{} {
	switch v := v.(type) {
	case map[interface{}]interface{}:
		return yamlTranslationizeMap(v)
	case []interface{}:
		return yamlTranslationizeSlice(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func yamlTranslationizeSlice(ts []interface{}) yamlTranslations {
	yts := make(yamlTranslations)
	for k, v := range ts {
		yts[fmt.Sprintf("%v", k)] = yamlTranslationizeValue(v)
	}
	return yts
}
