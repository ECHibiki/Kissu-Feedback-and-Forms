package former

import (
  "encoding/json"
  "errors"
  "strings"
  "regexp"

  // "fmt"
)

type InputType string
type SelectionCategory string
type FormValidationError string
type FormErrorCodes int64
type FormObjectTag string

const (
    Text InputType  = "text"
    Color           = "color"
    Date            = "date"
    Email           = "email"
    Number          = "number"
    Password        = "password"
    Range           = "range"
    Time            = "time"
    URL             = "url"
)

const (
  Checkbox SelectionCategory = "checkbox"
  Radio                      = "radio"
)

const (
  GroupMissingMessage   FormValidationError    = "Form subgroup does not contain any items or a description. Add one of these to submit."
  HeadMissingMessage                           = "Form group does not contain any items. Add one to submit."
  DuplicateNameMessage                         = "Form field names must be unique. Correct the duplicates."
  DuplicateIDMessage                           = "Form section IDs must be unique. Correct the duplicates."
  EmptyIDMessage                               = "IDs must have a value."
  InvalidIDStarterMessage                      = "IDs must start with letters."
  InvalidIDCharactersMessage                   = "IDs can only have the '-', '_', '.' or ':' characters."
  EmptyNameMessage                             = "Names must have a value."
  InvalidNameStarterMessage                    = "Names must start with letters."
  InvalidNameCharactersMessage                 = "Names can only have the '-', '_', '.' or ':' characters."
)
const (
  HeadMissingCode  FormErrorCodes = iota
  GroupMissingCode

  DuplicateNameCode
  DuplicateIDCode

  EmptyIDCode
  InvalidIDCharactersCode
  InvalidIDStarterCode
  EmptyNameCode
  InvalidNameCharactersCode
  InvalidNameStarterCode
)

const (
  TextAreaTag FormObjectTag = "textarea"
  GenericInputTag           = "genericinput"
  FileInputTag              = "fileinput"
  SelectionGroupTag         = "selectiongroup"
  OptionGroupTag            = "optiongroup"

)

func InputTypeFromString(raw_input string) (InputType , bool) {
  raw_input = strings.ToLower(raw_input)
  switch raw_input {
  case "text":
    return Text , true
  case "color":
    return Color , true
  case "date":
    return Date , true
  case "email":
    return Email , true
  case "number":
    return Number , true
  case "password":
    return Password , true
  case "range":
    return Range , true
  case "time":
    return Time , true
  case "url":
    return URL , true
  }
  return Text , false
}

func SelectionCategoryFromString(raw_input string) (SelectionCategory , bool) {
  raw_input = strings.ToLower(raw_input)
  switch raw_input {
  case "checkbox":
    return Checkbox , true
  case "radio":
    return Radio , true
  }
  return Checkbox , false
}

func FormObjectTagFromString(raw_input string) (FormObjectTag , bool ){
  raw_input = strings.ToLower(raw_input)
  switch raw_input {
  case "textarea":
    return TextAreaTag , true
  case "genericinput":
    return GenericInputTag , true
  case "fileinput":
    return FileInputTag , true
  case "selectiongroup":
    return SelectionGroupTag , true
  case "optiongroup":
    return OptionGroupTag , true
  }
  return GenericInputTag , false
}

type FailureObject struct {
  FailType FormValidationError
  FailCode FormErrorCodes
  FailPosition string
}

type FormConstruct struct {
  FormName string
  ID string
  Description string
  AnonOption bool
  FormFields []FormGroup
}

func (fc *FormConstruct) StorageName() string{
  remover := regexp.MustCompile("[^a-zA-Z0-9 \\-\\.]")

  safe_name := string(remover.ReplaceAll([]byte(fc.FormName) , []byte("")))
  safe_name = strings.ReplaceAll(safe_name , "." , "_")
  safe_name = strings.ReplaceAll(safe_name , "-" , "_")
  safe_name = strings.ReplaceAll(safe_name , " " , "_")
  if len(safe_name) > 250 {
    safe_name = safe_name[0:250]
  }
  return safe_name
}

type FormGroup struct {
  Label string
  ID string
  Description string
  SubGroup []FormGroup
  Respondables []UnmarshalerFormObject
}

//currently a placeholder to gain polymorphic properties
type FormObject interface{
    ElementType() string
    GetName() string
}

// implement Unmarshaler interface
type UnmarshalerFormObject struct {
  Type FormObjectTag
  Object FormObject
}

func (ufo *UnmarshalerFormObject) UnmarshalJSON(data []byte) error{
  var rfo struct {
    Type FormObjectTag
    Object json.RawMessage
  }
  err := json.Unmarshal(data, &rfo)
  if err != nil {
    return err

  }

  var fo FormObject
  switch rfo.Type {
    case TextAreaTag:
      ta  := TextArea{}
      err := json.Unmarshal(rfo.Object , &ta)
      if err != nil {
        return err
      }
      fo = ta
    case GenericInputTag:
      gi  := GenericInput{}
      err := json.Unmarshal(rfo.Object , &gi)
      if err != nil {
        return err
      }
      fo = gi
    case FileInputTag:
      fi  := FileInput{}
      err := json.Unmarshal(rfo.Object , &fi)
      if err != nil {
        return err
      }
      fo = fi
    case SelectionGroupTag:
      sg  := SelectionGroup{}
      err := json.Unmarshal(rfo.Object , &sg)
      if err != nil {
        return err
      }
      fo = sg
    case OptionGroupTag:
      og  := OptionGroup{}
      err := json.Unmarshal(rfo.Object , &og)
      if err != nil {
        return err
      }
      fo = og
    default:
      return errors.New("Unset OptionGroup type")
  }
  ufo.Type = rfo.Type
  ufo.Object = fo
  return nil
}

type Field struct {
  Label string
  Name string
  Required bool
}

type TextArea struct{
  Field Field
  Placeholder string
}
func (ta TextArea) ElementType() string {
  return "TEXTAREA"
}
func (ta TextArea) GetName() string {
  return ta.Field.Name
}

type GenericInput struct{
  Field Field
  Placeholder string
  Type InputType
}
func (gi GenericInput) ElementType() string {
  return "INPUT"
}
func (gi GenericInput) GetName() string {
  return gi.Field.Name
}

type FileInput struct{
  Field Field
  AllowedExtRegex string
  MaxSize int64
}
func (fi FileInput) ElementType() string {
  return "INPUT"
}
func (fi FileInput) GetName() string {
  return fi.Field.Name
}

type SelectionGroup struct{
  Field Field
  SelectionCategory SelectionCategory
  CheckableItems []Checkable
}
func (sg SelectionGroup) ElementType() string {
  return "INPUT"
}
func (sg SelectionGroup) GetName() string {
  return sg.Field.Name
}

type Checkable struct {
  Label string
}

type OptionGroup struct{
  Field Field
  Options []OptionItem
}

func (sg OptionGroup) ElementType() string {
  return "OPTION"
}
func (sg OptionGroup) GetName() string {
  return sg.Field.Name
}

type OptionItem struct {
  Label string
  Value string
}

func ValidateForm(form FormConstruct) (error_list []FailureObject) {
  uniqueness_errors := checkNameAndIDUniqueness(form)
  if len(uniqueness_errors) > 0 {
    error_list = append(error_list , uniqueness_errors...)
  }

  character_errors := checkNameAndIDPropperCharacters(form)
  if len(character_errors) > 0 {
    error_list = append(error_list , character_errors...)
  }

  // uniqueness errors through all other chekcs into confusion
  // allowing for other errors to display isn't important... not even golang shows all error types at once
  // this limitation, I guess, issue will apply to situations where a submits without the client
  if len(uniqueness_errors) > 0 {
    return error_list
  }

  struct_errors := checkValidFormStructure(form)
  if len(struct_errors) > 0 {
    error_list = append(error_list , struct_errors...)
  }
  return error_list
}

func checkValidFormStructure(form FormConstruct) (error_list []FailureObject)  {
  if len(form.FormFields) == 0 {
    return []FailureObject{ { HeadMissingMessage , HeadMissingCode, form.ID } }
  }
  var subgroup_stack []FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // fail location identified by an ID
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    fail_location  := item.ID
    if fail_location != "" && len(item.Respondables) == 0 && item.Description == "" {
      error_list = append(error_list , FailureObject{ GroupMissingMessage , GroupMissingCode, fail_location } )
    }
    // verify it has validity to it
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }
  return error_list
}

// For the next few, create a []struct {isID:bool, name:string }.
// From this struct perform the checks
func checkNameAndIDUniqueness(form FormConstruct) (error_list []FailureObject)  {
  if len(form.FormFields) == 0 {
    return []FailureObject{}
  }
  var id_checklist  = map[string]uint{ form.ID : 1 }
  var failing_ids []string
  var name_checklist  = make(map[string]uint)
  var failing_names []string
  var subgroup_stack []FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // fail location identified by an ID
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    if _, ok := id_checklist[item.ID] ; !ok {
      id_checklist[item.ID] = 1;
    } else if id_checklist[item.ID] != 2 {
      id_checklist[item.ID] = 2;
      failing_ids = append(failing_ids , item.ID )
    }
    if len(item.Respondables) != 0 {
      for _ , r := range item.Respondables {
        name := r.Object.GetName()
        if _, ok := name_checklist[name] ; !ok {
          name_checklist[name] = 1;
        } else if name_checklist[name] != 2 {
          name_checklist[name] = 2;
          failing_names = append( failing_names , name )
        }
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }

  for _ , id := range failing_ids {
    error_list = append(error_list , FailureObject{ DuplicateIDMessage , DuplicateIDCode, id } )
  }
  for _ , name := range failing_names {
    error_list = append(error_list , FailureObject{ DuplicateNameMessage , DuplicateNameCode, name } )
  }


  return error_list
}

func checkNameAndIDPropperCharacters(form FormConstruct) (error_list []FailureObject)  {
  if len(form.FormFields) == 0 {
    return []FailureObject{}
  }

  var ids []string = []string{form.ID}
  var names []string
  var subgroup_stack []FormGroup
  subgroup_stack = append(subgroup_stack , form.FormFields...)
  // fail location identified by an ID
  for len(subgroup_stack) > 0  {
    item := subgroup_stack[len(subgroup_stack) - 1]
    subgroup_stack = subgroup_stack[:len(subgroup_stack) - 1]
    ids = append(ids , item.ID )
    if len(item.Respondables) != 0 {
      for _ , r := range item.Respondables {
        name := r.Object.GetName()
        names = append( names , name )
      }
    }
    if len(item.SubGroup) != 0 {
      // add children to the stack
      subgroup_stack = append(subgroup_stack , item.SubGroup...)
    }
  }


  invalid_entry := regexp.MustCompile("^[^a-zA-Z]")
  invalid_body := regexp.MustCompile("[^a-zA-Z0-9\\-_:\\.]")
  for _ , id := range ids {
    if len(id) == 0 {
      error_list = append(error_list , FailureObject{ EmptyIDMessage , EmptyIDCode, id } )
    }
    if invalid_entry.Match([]byte(id)) {
      error_list = append(error_list , FailureObject{ InvalidIDStarterMessage , InvalidIDStarterCode, id } )
    }
    if invalid_body.Match([]byte(id)) {
      error_list = append(error_list , FailureObject{ InvalidIDCharactersMessage , InvalidIDCharactersCode, id } )
    }
  }
  for _ , name := range names {
    if len(name) == 0 {
      error_list = append(error_list , FailureObject{ EmptyNameMessage , EmptyNameCode, name } )
    }
    if invalid_entry.Match([]byte(name)) {
      error_list = append(error_list , FailureObject{ InvalidNameStarterMessage , InvalidNameStarterCode, name } )
    }
    if invalid_body.Match([]byte(name)) {
      error_list = append(error_list , FailureObject{ InvalidNameCharactersMessage , InvalidNameCharactersCode, name } )
    }
  }
  return error_list
}
