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
  "strconv"
  "errors"
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
  var construction_variables map[string]stick.Value = map[string]stick.Value{"form" : rebuild_group }

  // Render a form only used for testing
  // fmt.Println(construction_variables["responables"])
  testing_form_render, err := templater.ReturnFilledTemplate(env , root_dir + "/templates/test-views/render-test.twig" , construction_variables)
  return testing_form_render , err
}

func GetAllForms(db *sql.DB) (parsed_row_list []types.FormDBFields , err error){
  row_list , err := db.Query("SELECT id, name, updated_at FROM forms ORDER BY updated_at DESC")
  if err != nil {
    return
  }
  defer row_list.Close()
  for row_list.Next(){
    var parsed_row types.FormDBFields
    err = row_list.Scan( &parsed_row.ID ,  &parsed_row.Name , &parsed_row.UpdatedAt )
    if err != nil{
      return
    }
    parsed_row_list = append(parsed_row_list , parsed_row)
  }
  return
}

func GetRepliesToForm(db *sql.DB , id int64)  (parsed_row_list []types.ResponseDBFields , err error){
  row_list , err := db.Query("SELECT id, fk_id, identifier, response_json, submitted_at FROM responses WHERE fk_id = ? ORDER BY id ASC" , id)
  if err != nil {
    return
  }
  defer row_list.Close()
  for row_list.Next(){
    var parsed_row types.ResponseDBFields
    err = row_list.Scan( &parsed_row.ID , &parsed_row.FK_ID , &parsed_row.Identifier ,  &parsed_row.ResponseJSON , &parsed_row.SubmittedAt )
    if err != nil{
      return
    }
    parsed_row_list = append(parsed_row_list , parsed_row)
  }
  return
}

func CreateInstancedCSVForGivenForm(db *sql.DB , id int64 , initialization_folder string) error{
    form_data , err := tools.GetFormOfID(db , 2)
    if err != nil {
      return err
    }
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form_data.FieldJSON) , &form_construct)
    if err != nil {
      return err
    }
    var csv_list [][]string
    var field_list []string
    field_list = append(field_list , "Identifier")
    fields := GetFieldsOfFormConstruct(form_construct)
    field_list = append(field_list , "SubmittedAt")
    var field_map map[string]int = make(map[string]int)
    for i , field := range fields {
      field_map[field.GetName()] = i
    }
    csv_list = append(csv_list , field_list)


    responses, err := GetRepliesToForm(db , 2)
    for _ , r := range responses {
      responses_list := make([]string , len(fields) + 2)
      responses_list[field_map["Identifier"]] =  r.Identifier
      responses_list[field_map["SubmittedAt"]] = strconv.Itoa(int(r.SubmittedAt))
      var response map[string]string = make(map[string]string)
      err = json.Unmarshal([]byte(r.ResponseJSON) , &response)
      if err != nil {
        return err
      }
      for _ , v := range response {
        responses_list[field_map[v]] = v
      }
      csv_list = append(csv_list , responses_list)
    }
    err = tools.WriteCSVToDir(initialization_folder + "/data/" + form_data.Name + "/data.csv" , csv_list)
    return err
}

func CreateReadmeForGivenForm(db *sql.DB , id int64 , initialization_folder string) error{
  form_data , err := tools.GetFormOfID(db , 2)
  if err != nil {
    return err
  }
  var form_construct former.FormConstruct
  err = json.Unmarshal([]byte(form_data.FieldJSON) , &form_construct)
  if err != nil {
    return err
  }

  fields := GetFieldsOfFormConstruct(form_construct)

  var field_map map[string]string  = make(map[string]string)
  field_map["FormName"] = form_data.Name
  field_map["ID"] = strconv.Itoa(int(form_data.ID))
  field_map["FormDescription"] = form_construct.Description
  field_map["AnonOption"] = strconv.FormatBool(form_construct.AnonOption)
  for _ , field := range fields {
    field_map[field.GetName()] = field.GetDescription()
  }
  err = tools.WriteJSONReadmeToDir(initialization_folder + "/data/" + form_data.Name + "/field-descriptors.json" , field_map)
  return err
}

func GetFieldsOfFormConstruct(form former.FormConstruct) (field_list []former.FormObject){
  if len(form.FormFields) == 0 {
    return
  }
  var subgroup_stack []former.FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // fail location identified by an ID
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if len(item.Respondables) != 0 {
      for _ , r := range item.Respondables {
        name := r.Object.GetName()
        name_found := false
        for _, v := range(field_list){
          if v.GetName() == name {
            name_found = true
            break;
          }
        }
        if !name_found {
          field_list = append(field_list , r.Object)
        }
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }
  return
}


func GetResponseByID(db *sql.DB , id int64) (types.ResponseDBFields , error) {
  data , err := tools.GetResponseByID(db,id)
  if err != nil{
    return types.ResponseDBFields{} , err
  }
  if data.FK_ID == 0 {
    return types.ResponseDBFields{} , errors.New("Database has no row for ID")
  }
  return data, nil
}

func GetFormByNameAndID(db *sql.DB , name string, id int64) (types.FormDBFields , error) {
  data := db.QueryRow("SELECT * FROM forms WHERE name = ? AND id = ?" , name , id)

  var rtn types.FormDBFields
  err := data.Scan(&rtn.ID, &rtn.Name , &rtn.FieldJSON , &rtn.UpdatedAt)
  if err != nil{
    return types.FormDBFields{} , err
  }
  return rtn, nil

}
