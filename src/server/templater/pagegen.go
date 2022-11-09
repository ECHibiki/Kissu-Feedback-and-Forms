package templater

import (
	"github.com/flosch/pongo2"
	"fmt"
)

var tpl_in_memory map[string]*pongo2.Template = map[string]*pongo2.Template{}
var root_dir string

func init(){
	fmt.Printf("Pongo2 Autoescape set to false\n")
	pongo2.SetAutoescape(false)
	err := pongo2.RegisterFilter("keyvalue", FilterGetValueByKey)
	if err != nil{
		panic(err)
	}
	err = pongo2.RegisterFilter("tostring", ToString)
	if err != nil{
		panic(err)
	}
}

// from https://github.com/flosch/pongo2/issues/162
func FilterGetValueByKey(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
        m := in.Interface().(map[string]string)
        return pongo2.AsValue(m[param.String()]), nil
}

func ToString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error){
		return pongo2.AsValue(in.String()) , nil
}


func SetRootDir(root string){
	fmt.Printf("Root template directory: %s\n" , root)
	root_dir = root
}

func ReturnFilledTemplate(template_path string, ctx_values pongo2.Context) (string, error) {
	template , exists := tpl_in_memory[template_path]
	var err error
	if true || !exists{
		template , err = pongo2.FromFile(root_dir + template_path)
		if err != nil {
			return "", err
		}
		tpl_in_memory[template_path] = template
	}
	fmt.Println(template , ctx_values)
	return template.Execute(ctx_values)
}
