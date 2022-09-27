package responder

import (
  "os"
  "io"
  "io/ioutil"
  "strings"
  "errors"
  "time"
  "encoding/json"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"

)
// c.Request.FromFile
// http.DetectContentType

// Using header information can be problematic, just check the file itself
// Filename not important and would rather not deal with any issues relating to it
// Add to the templates that filenames are not preserved

// The input names of checkgroups are ignored for what the form says they should be

/*
  ResponseMissingMessage                       = "A field is required yet has no response."
  InvalidInputMessage                     ...
  InvalidOptionValueMessage                    = "The value of an options group does not exist on the server."
  InvalidFileExtMessage                        = "The extention of a file is not permitted on the server."
  InvalidFileSizeMessage                       = "The size of a file is too large."
*/

func ValidateTextResponsesAgainstForm(text_responses map[string]string , form former.FormConstruct)  (error_list []former.FailureObject){
  r_fo := validateRequiredTextFields( text_responses , form )
  i_fo := validateResponseTextFields( text_responses , form )
  error_list = append(error_list , r_fo...)
  error_list = append(error_list , i_fo...)
  return
}

func validateRequiredTextFields(text_responses map[string]string , form former.FormConstruct) (error_list []former.FailureObject){
  var subgroup_stack []former.FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)

  // Validate required fields and select/option group first pass verification
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if len(item.Respondables) != 0 {
      for _ , respondable := range item.Respondables {
        fail := false
        required := respondable.Object.GetRequired()
        if !required {
          continue
        }

        name := respondable.Object.GetName()
        select_group , is_selection := respondable.Object.(former.SelectionGroup)
        if is_selection && select_group.SelectionCategory == former.Checkbox {
          answer_found := false
          for i , _ := range select_group.CheckableItems {
            response , ok := text_responses[ name + "-" + strconv.Itoa(i + 1) ]
            if ok {
              answer_found = true
            }
          }
          if !answer_found {
            error_list = append(error_list , former.FailureObject{ former.ResponseMissingMessage , former.ResponseMissingCode, name  })
          }
        } else{
          response , ok := text_responses[name]
          if !ok{
            fail = true
          } else if required && ok {
            trimed := strings.TrimSpace(response)
            if len(trimed) == 0 {
              fail = true
            }
          }
          if fail
            error_list = append(error_list , former.FailureObject{ former.ResponseMissingMessage , former.ResponseMissingCode, name  })
          }
        }
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }

  // Next validate that the responses actually make sense against the form
  return
}

func validateResponseTextFields(text_responses map[string]string , form former.FormConstruct) (error_list []former.FailureObject){
  var subgroup_stack []former.FormGroup
  var field_list map[string]former.UnmarshalerFormObject

  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // Validate required fields and select/option group first pass verification
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if len(item.Respondables) != 0 {
      for _ , respondable := range item.Respondables {
        name := respondable.Object.GetName()

        select_group , is_selection := respondable.Object.(former.SelectionGroup)
        if is_selection && select_group.SelectionCategory == former.Checkbox {
          for i , v := range select_group.CheckableItems {
            field_list[ name + "-" + strconv.Itoa(i + 1) ] = respondable
          }
        } else{
          field_list[name] = respondable
        }
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }

  for field , response := range text_responses{
    respondable , exists := field_list[field]
    if !exists{
      error_list = append(error_list , former.FailureObject{ former.InvalidInputMessage , former.InvalidInputCode, name  })
      continue
    }
    options_respondable , is_options := respondable.Object.type(former.OptionGroup)
    if is_options {
      found_value := false
      for _ , opt := range options_respondable.Options {
        if response == opt.Value {
          found_value = true
          break
        }
      }
      if !found_value {
        error_list = append(error_list , former.FailureObject{ former.InvalidOptionValueMessage , former.InvalidOptionValueCode, name  })
      }
    }
  }

  return
}

func ValidateFileObjectsAgainstForm(file_tags map[string]former.MultipartFile , form former.FormConstruct) (error_list []former.FailureObject){
  r_fo := validateRequiredFileFields( file_tags , form )
  i_fo := validateResponseFileFields( file_tags , form )
  error_list = append(error_list , r_fo...)
  error_list = append(error_list , i_fo...)
  return
}

func validateRequiredFileFields(file_tags map[string]former.MultipartFile , form former.FormConstruct) (error_list []former.FailureObject){

    var subgroup_stack []former.FormGroup
    subgroup_stack = append(subgroup_stack , form.FormFields...)

    // fail location identified by an ID
    for len(subgroup_stack) > 0  {
      item := subgroup_stack[len(subgroup_stack) - 1]
      subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
      if len(item.Respondables) != 0 {
        for _ , respondable := range item.Respondables {
          fail := false
          name := respondable.Object.GetName()
          required := respondable.Object.GetRequired()
          response , ok := file_tags[name]
          if required && !ok{
            fail = true
          } else if ok {
            // begin file verification
            var FileInput former.MultipartFile = respondable.Object.(former.FileInput)
            

          }
          if fail {
            fo := former.FailureObject{ former.ResponseMissingMessage , former.ResponseMissingCode, name  }
            error_list = append(error_list , fo)
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

func validateResponseFileFields(file_tags map[string]former.MultipartFile , form former.FormConstruct) (error_list []former.FailureObject) {
  var subgroup_stack []former.FormGroup
  var field_list map[string]former.UnmarshalerFormObject

  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // Validate required fields and select/option group first pass verification
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if len(item.Respondables) != 0 {
      for _ , respondable := range item.Respondables {
        name := respondable.Object.GetName()
        field_list[name] = respondable
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }

  for field , response := range text_responses{
    respondable , exists := field_list[field]
    if !exists{
      error_list = append(error_list , former.FailureObject{ former.InvalidInputMessage , former.InvalidInputCode, name  })
      continue
    }
  }

  return
}


//
func CreateResponderFolder(root_dir string , response_struct former.FormResponse) error{
  err := os.Mkdir(root_dir + "/" + response_struct.FormName + "/" + response_struct.ResponderID + "/" , 0755)
  if err != nil {
    return err
  }
  if len(response_struct.FileObjects) > 0 {
    err = os.Mkdir(root_dir + "/" + response_struct.FormName + "/" + response_struct.ResponderID + "/files/" , 0755)
  }
  return err
}

func WriteResponsesToJSONFile(root_dir string , resp former.FormResponse) error {
  storage_dir := root_dir + "/" + resp.FormName + "/" + resp.ResponderID + "/"

  json_resp := ConvertFormResponseToJSONFormResponse(root_dir , resp)

  json_bytes , err := json.MarshalIndent(json_resp , "" , " ")
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(storage_dir + "responses.json" , json_bytes , 0644)
  return err
}

func ConvertFormResponseToJSONFormResponse(root_dir string, resp former.FormResponse) former.JSONFormResponse {
  json_resp := former.JSONFormResponse{}
  json_resp.FormName= resp.FormName
  json_resp.RelationalID = resp.RelationalID
  json_resp.ResponderID = resp.ResponderID
  json_resp.Responses = resp.Responses
  storage_dir := root_dir + "/" + resp.FormName + "/" + resp.ResponderID + "/"
  for k, v := range(resp.FileObjects) {
    json_resp.FilePaths[k] = storage_dir + "files/" + v.Header.Filename
  }

  return json_resp
}

func WriteFilesFromMultipart(root_dir string , response_struct former.FormResponse) []error{
  storage_dir := root_dir + "/" + response_struct.FormName + "/" + response_struct.ResponderID + "/files/"
  var err_list []error = []error{}
  for _, v := range response_struct.FileObjects {
    fname := v.Header.Filename
    if strings.Contains(fname , "/"){
      tools.LogError( storage_dir , storage_dir +  fname )
      err_list = append(err_list , errors.New("File " + fname +  " contained illegal characters"))
      continue
    }
    handler, err := os.OpenFile(storage_dir +  fname , os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
      tools.LogError( storage_dir , storage_dir +  fname )
      err_list = append(err_list , err)
      continue
    }
    defer handler.Close()
    _ , err = io.Copy(handler , v.File )
     if err != nil {
       tools.LogError( storage_dir , storage_dir +  fname )
       err_list = append(err_list , err)
       continue
     }
  }
  return err_list
}

func FormResponseToDBFormat(response former.FormResponse) (types.ResponseDBFields , error){
  for k , v := range(response.FileObjects) {
    response.Responses[k] = v.Header.Filename
  }

  resp_bytes , err := json.Marshal( response.Responses )
  return types.ResponseDBFields {
    ID: 0,
    FK_ID: response.RelationalID,
    Identifier: response.ResponderID,
    ResponseJSON: string(resp_bytes),
    SubmittedAt: time.Now().Unix(),
  } , err
}
