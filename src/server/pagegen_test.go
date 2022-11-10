package main

import (
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	"github.com/flosch/pongo2"
	"testing"
	"fmt"
)

func TestGenericGeneration(t *testing.T) {
	var template_path string = "../../test/templates/test-views/hello-world-test.html"
	// Write an html file created from twig templates into the proper directory
	value_map := pongo2.Context{
		"a_string": "A title field",
		"an_int":   10,
		"a_float":  -0.001,
	}

	parsed_template, err := templater.ReturnFilledTemplate(template_path, value_map)
	if parsed_template != "A title field10-0.001000\n" {
		if err != nil {
			t.Fatal("Templater produced an error: ", err)
		} else {
			t.Fatal("Template did not output expected result: \n\tExpected: \"A title field10-0.001000\\n\"\n\tGot: \t", fmt.Sprintf("%#v" , parsed_template))
		}
	}
}
