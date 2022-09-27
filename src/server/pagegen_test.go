package main

import (
  "testing"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
  "github.com/tyler-sommer/stick"
)

func TestGenericGeneration(t *testing.T) {
  var template_path string = "../../test/templates/test-views/hello-world-test.twig"
  // Write an html file created from twig templates into the proper directory
  env := templater.ReturnTemplateHandler()
  value_map  := map[string]stick.Value {
    "a_string": "A title field",
    "an_int": 10,
    "a_float": -0.001,
  }

  parsed_template, err := templater.ReturnFilledTemplate(env , template_path, value_map)
  if parsed_template != "A title field10-0.001\n"{
    if err != nil {
      t.Fatal("Templater produced an error: " , err)
    } else{
      t.Fatal("Template did not output expected result: \n\tExpected: A title field10-0.001\\n\n\tGot: " , parsed_template)
    }
  }
}
