package templater

import (
  "io/ioutil"
  "bytes"

  "github.com/tyler-sommer/stick"
  "github.com/tyler-sommer/stick/twig"
)

func ReturnTemplateHandler() *stick.Env{
   return twig.New(nil)
}

func ReturnFilledTemplate(env *stick.Env, template_path string, value_map map[string]stick.Value) (string , error){
  template_bytes, err := ioutil.ReadFile(template_path)
  if err != nil{
    return "" , err
  }
  // value_typecast := stringMapToStickValue(value_map)
  template_buffer := new(bytes.Buffer)
  if err := env.Execute(string(template_bytes), template_buffer, value_map ); err != nil {
    return "", err
  }
  return template_buffer.String() , nil
}
