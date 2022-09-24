package builder

import (
  "encoding/json"
  "strings"
  "time"
  "os"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
)

func InputTypeFromString(raw_input string) (former.InputType , bool) {
  raw_input = strings.ToLower(raw_input)
  switch raw_input {
  case "text":
    return former.Text , true
  case "color":
    return former.Color , true
  case "date":
    return former.Date , true
  case "email":
    return former.Email , true
  case "number":
    return former.Number , true
  case "password":
    return former.Password , true
  case "range":
    return former.Range , true
  case "time":
    return former.Time , true
  case "url":
    return former.URL , true
  }
  return former.Text , false
}

func SelectionCategoryFromString(raw_input string) (former.SelectionCategory , bool) {
  raw_input = strings.ToLower(raw_input)
  switch raw_input {
  case "checkbox":
    return former.Checkbox , true
  case "radio":
    return former.Radio , true
  }
  return former.Checkbox , false
}

func ValidateForm(form former.FormConstruct) ([]former.FailureObject ) {
  var error_list []former.FailureObject

  var subgroup_stack []former.FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)
  for i := int64(0) ; len(subgroup_stack) > 0 ; i++ {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if len(item.Respondables) == 0 && item.Description == "" {
      error_list = append(error_list , former.FailureObject{ former.GroupMissingError , i } )
    }
    // verify it has validity to it
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }
  return error_list
}

func MakeFormWritable(form former.FormConstruct) (types.FormDBFields  , error){
  marshaled_form , err := json.Marshal(form)
  if err != nil{
    return types.FormDBFields{}, err
  }
  return types.FormDBFields{
    ID: 0 ,
    FieldJSON: string(marshaled_form),
    UpdatedAt: time.Now().Unix(),
  } , nil
}

func CreateFormDirectory(form former.FormConstruct , cfg types.ConfigurationSettings) error{
  safe_name := strings.ReplaceAll(form.FormName , "." , "-")
  safe_name = strings.ReplaceAll(safe_name , "/" , "_")
  err := os.Mkdir(cfg.ResoruceDirectory + "/data/" + safe_name + "/" , 0664 )
  if err != nil {
    return err
  }
  err = os.Mkdir(cfg.ResoruceDirectory + "/data/" + safe_name + "/files" , 0664)
  return err
}
