package katolomb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Interpolator is the interface that wraps the basic interpolate method.
//
// Interpolate takes a string and a TranslationProperties and returns the string
// with any interpolations declared replaced with the values of the properties
// in the TranslationProperties parameter.
//
// Interpolation in-text declaration format must be decided by the
// implementation.
type Interpolator interface {
	Interpolate(string, TranslationProperties) (string, error)
}

// InterpolatorFunc wraps a function with the Interpolator's Interpolate method
// signature to satisfy the Interpolator interface.
type InterpolatorFunc func(string, TranslationProperties) (string, error)

type interpolator struct {
	regexp *regexp.Regexp
}

type interpolation struct {
	literal         string
	property        string
	hasDefaultValue bool
	defaultValue    string
}

var defaultVanillaInterpolatorRegexp = regexp.MustCompile(`%\{(?P<name>[^\}\|]+)(?P<default>\|[^\}]*)?\}`)

// NewNoErrorInterpolator returns an Interpolator that wraps another
// Interpolator passed as parameter to avoid returning errors. When the wrapped
// Interpolator's interpolation returns an error, the original text is returned
// instead.
func NewNoErrorInterpolator(i Interpolator) Interpolator {
	return InterpolatorFunc(func(text string, props TranslationProperties) (string, error) {
		interpText, err := i.Interpolate(text, props)
		if err != nil {
			interpText = text
		}
		return interpText, nil
	})
}

// NewInterpolator returns an Interpolator that can detect and interpolate
// interpolation declarations with the following format:
//   %{<property name>|<default value>}
// where:
//   * <property name> is the name of the property to interpolate.
//   * <default value> is the value to interpolate when the property is not
//   available on the TranslationProperties
//   * the |<default value> part  is optional
func NewInterpolator() Interpolator {
	return &interpolator{defaultVanillaInterpolatorRegexp}
}

// Interpolate takes a string and a TranslationProperties and returns the string
// with any interpolations declared replaced with the right properties in the
// TranslationProperties parameter or a default value in the interpolation
// declaration if available. If a default value is not provided and the property
// to interpolate is not available on the TranslationProperties parameter, an
// error will be returned.
//
// The interpolation in-text declaration format is
//  %{<property name>|<default value>}
// where:
//   * <property name> is the name of the property to interpolate.
//   * <default value> is the value to interpolate when the property is not
//   available on the TranslationProperties
//   * the |<default value> part  is optional
func (i *interpolator) Interpolate(text string, properties TranslationProperties) (string, error) {
	for _, interpol := range i.findInterpolations(text) {
		value, err := properties.Property(interpol.property)
		if err != nil {
			if !interpol.hasDefaultValue {
				return "", fmt.Errorf("interpolating %v: %v", strconv.Quote(text), err)
			}
			value = interpol.defaultValue
		}
		text = strings.Replace(text, interpol.literal, value, -1)
	}
	return text, nil
}

// Interpolate calls the function with the received text and properties
// parameters and returns the result.
func (inF InterpolatorFunc) Interpolate(text string, properties TranslationProperties) (string, error) {
	return inF(text, properties)
}

func (i *interpolator) findInterpolations(text string) []*interpolation {
	interpolations := []*interpolation{}
	for _, interpolationSlice := range i.regexp.FindAllStringSubmatch(text, -1) {
		var defaultValue string
		hasDefaultValue := len(interpolationSlice) > 2 && len(interpolationSlice[2]) > 0
		if hasDefaultValue {
			defaultValue = interpolationSlice[2][1:]
		}
		interpol := &interpolation{
			literal:         interpolationSlice[0],
			property:        interpolationSlice[1],
			hasDefaultValue: hasDefaultValue,
			defaultValue:    defaultValue,
		}
		interpolations = append(interpolations, interpol)
	}
	return interpolations
}
