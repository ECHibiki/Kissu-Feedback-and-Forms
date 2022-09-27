package returner

import (
  "github.com/tyler-sommer/stick"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "database/sql"
  "encoding/json"
  "fmt"
)

func RenderTestingTemplate[T int64 | string](db *sql.DB, env *stick.Env, root_dir string,  db_key T) (string , error){

  var returned_form types.FormDBFields
  var rebuild_group former.FormConstruct
  var err error

  var i interface{} = db_key
  switch i.(type) {
    case int64:
      returned_form , err = tools.GetFormOfID(db, i.(int64))
    case string:
      returned_form , err = tools.GetFormOfName(db, i.(string))
  }
  if err != nil{
    fmt.Println(db_key)
    panic(err)
  }

  err = json.Unmarshal([]byte(returned_form.FieldJSON), &rebuild_group)
  if err != nil{
    return "" , err
  }

  // Turn rebuild_group into a templatable format
  var construction_variables map[string]stick.Value = map[string]stick.Value{"form" : rebuild_group}

  // Render a form only used for testing
  testing_form_render, err := templater.ReturnFilledTemplate(env , root_dir + "/templates/test-views/render-test.twig" , construction_variables)
  return testing_form_render , err
}
