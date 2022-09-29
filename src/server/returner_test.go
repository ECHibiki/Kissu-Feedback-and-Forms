package main

import (
    "testing"
    "strings"
    "encoding/json"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
)

func TestRetrieval(t *testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  demo_form_id_check_name := "../Test form 1"
  tools.DoFormInitialization(demo_form_id_check_name , "a-simple-identifier" , db , cfg)

  demo_form_name_check_assumed_storage_name := "__alternative_test_form_1"
  demo_form_name_check_name := "../alternative test form 1"
  tools.DoFormInitialization(demo_form_name_check_name , "a-simple-identifier" , db , cfg)

// ---- Forget the initialization of fields

  var insertable_form_id , insertable_form_name former.FormConstruct

  insertable_form_id_db , err := tools.GetFormOfID(db, 1)
  if err != nil {
    panic(err)
  }
  json.Unmarshal([]byte(insertable_form_id_db.FieldJSON) , &insertable_form_id)
  insertable_form_name_db, err := tools.GetFormOfName(db, demo_form_name_check_assumed_storage_name)
  if err != nil {
    panic(err)
  }
  json.Unmarshal([]byte(insertable_form_name_db.FieldJSON) , &insertable_form_name)

  env := templater.ReturnTemplateHandler()
  // generics for outputting a template depending on the form ID
  // Output should be an html page replicating the effect of marshal on the struct
  template_id , err := returner.RenderTestingTemplate(db, env, initialization_folder, int64(1))
  if err != nil{
    panic(err)
  }
  form_id_marshal, _ := json.Marshal(insertable_form_id)
  template_id = strings.ReplaceAll(template_id , "\n" , "")
  tilen := len(template_id)
  if template_id != string(form_id_marshal)[:tilen] {
    t.Error("Test template render by ID failed--\nCreation:" , string(template_id) , "\nAssmpton:", string(form_id_marshal)[:tilen]  )
  }

  template_name , err := returner.RenderTestingTemplate(db, env, initialization_folder , demo_form_name_check_assumed_storage_name)
  if err != nil{
    panic(err)
  }
  form_name_marshal, _ := json.Marshal(insertable_form_name)
  template_name = strings.ReplaceAll(template_name , "\n" , "")
  tiname := len(template_name)
  if template_name != string(form_name_marshal)[:tiname] {
    t.Error("Test template render by Name failed--\nCreation:" , string(template_name) , "\nAssmpton:", string(form_name_marshal)[:tiname] )
  }
}
