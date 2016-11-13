package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pbanos/katolomb"
)

func main() {
	locale, err := ioutil.ReadFile("locales/en.yml")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	interpolator := katolomb.NewInterpolator()
	translator, err := katolomb.NewYAMLTranslator(locale)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(2)
	}
	translator = katolomb.NewInterpolatedTranslator(translator, interpolator)
	props := katolomb.NewTranslationProperties(map[string]string{
		"cacahuete": "hola",
		"timestamp": "today",
		// "name":      "Mr. Darcy",
	})
	translation, err := translator.Translate("en.my.message", props)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(3)
	} else {
		fmt.Printf("%v\n", translation)
	}
}
