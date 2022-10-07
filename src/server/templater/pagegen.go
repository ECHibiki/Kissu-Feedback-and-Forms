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

func ReturnFileSystemTemplateHandler(root_dir string) *stick.Env{
  return stick.New(stick.NewFilesystemLoader(root_dir))
}

func ReturnFilledTemplate(env *stick.Env, template_path string, value_map map[string]stick.Value) (string , error){
  _ , err := env.Loader.Load("./a-non-existent-file")
  var template_string string
  if err != nil{
    template_string = template_path
  } else{
    template_bytes , err := ioutil.ReadFile(template_path)
    template_string = string(template_bytes)
    if err != nil{
      return "" , err
    }
  }
  template_buffer := new(bytes.Buffer)
  if err := env.Execute(template_string, template_buffer, value_map ); err != nil {
    return "", err
  }
  return template_buffer.String() , nil
}
