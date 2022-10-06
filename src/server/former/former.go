package former

import (
  "encoding/json"
  "errors"
  "strings"
  "regexp"
  "mime/multipart"
  "crypto/md5"
  // "io"
  "fmt"
)

type InputType string
type SelectionCategory string
type FormValidationError string
type FormErrorCode int64
type FormObjectTag string

func (fot *FormObjectTag) String() string{
  return fmt.Sprintf("working? %v", fot)
}
func (fot *FormObjectTag) string() string{
  return fmt.Sprintf("working? %v", fot)
}

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

  DuplicateFormNameMessage                     = "The form name conflicts with another form. Change the name."
  DuplicateNameMessage                         = "Form field names must be unique. Correct the duplicates."
  DuplicateIDMessage                           = "Form section IDs must be unique. Correct the duplicates."

  EmptyIDMessage                               = "IDs must have a value."
  InvalidIDStarterMessage                      = "IDs must start with letters."
  InvalidIDCharactersMessage                   = "IDs can only have the '-', '_', '.' or ':' characters."
  EmptyNameMessage                             = "Names must have a value."
  InvalidNameStarterMessage                    = "Names must start with letters."
  InvalidNameCharactersMessage                 = "Names can only have the '-', '_', '.' or ':' characters."
  InvalidCheckboxMessage                       = "A checkbox creates fields of the given name followed by 'hypen number'(eg. name-3). A checkbox conflicts with other fields that end with 'hypen number'."

  ResponseMissingMessage                       = "A field is required yet has no response."
  InvalidInputMessage                          = "A field filled out does not actually exist on the server."
  InvalidSelectionIndexMessage                 = "A selection group's position does not make sense."
  InvalidSelectionValueMessage                 = "The value of a selection group does not exist on the server."
  InvalidOptionValueMessage                    = "The value of an options group does not exist on the server."
  InvalidFileExtMessage                        = "The extention of a file is not permitted on the server."
  InvalidFileSizeMessage                       = "The size of a file is too large."
  InvalidExtRegexMessage                       = "The form's regex is invalid."

  DangerousPathMessage                         = "Path contains illegal characters"

  EditNameChangeMessage                        = "The name can not change from an edit"
)
const (
  HeadMissingCode  FormErrorCode = iota
  GroupMissingCode

  DuplicateFormNameCode
  DuplicateNameCode
  DuplicateIDCode

  EmptyIDCode
  InvalidIDCharactersCode
  InvalidIDStarterCode
  EmptyNameCode
  InvalidNameCharactersCode
  InvalidNameStarterCode
  InvalidCheckboxCode

  ResponseMissingCode
  InvalidInputCode
  InvalidSelectionIndexCode
  InvalidSelectionValueCode
  InvalidOptionValueCode
  InvalidFileExtCode
  InvalidFileSizeCode
  InvalidExtRegexCode

  DangerousPathCode

  EditNameChangeCode
)

const (
  TextAreaTag FormObjectTag             = "textarea"
  GenericInputTag                       = "genericinput"
  FileInputTag                          = "fileinput"
  SelectionGroupTag                     = "selectiongroup"
  OptionGroupTag                        = "optiongroup"

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
  FailCode FormErrorCode
  FailPosition string
}

type FormConstruct struct {
  FormName string
  ID string
  Description string
  // With anon option set to true there is an ability for users to flag themselves as anonymous
  // Under other conditions there is no anonymity
  AnonOption bool
  FormFields []FormGroup
}

type MultipartFile struct{
  File multipart.File
  Header *multipart.FileHeader
}

type FormResponse struct {
  FormName string
  // The DB ID of the form
  RelationalID int64
  // IP or hash of IP
  ResponderID string
  Responses map[string]string
  FileObjects map[string]MultipartFile
}

type JSONFormResponse struct {
  FormName string
  // The DB ID of the form
  RelationalID int64
  // IP or hash of IP
  ResponderID string
  Responses map[string]string
  FilePaths map[string]string
}

func (fr *FormResponse)ScrambleResponderID() {
  fr.ResponderID = fmt.Sprintf("%x", md5.Sum([]byte(fr.ResponderID)))
}
func (fr *FormResponse)GetScrambledID() string {
  return fmt.Sprintf("%x", md5.Sum([]byte(fr.ResponderID)))
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
    GetDescription() string
    GetRequired() bool
}

// implement Unmarshaler interface
type UnmarshalerFormObject struct {
  Type FormObjectTag
  Object FormObject
}

func (ufo *UnmarshalerFormObject) GetType() string{
  return string(ufo.Type)
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
func (ta TextArea) GetRequired() bool {
  return ta.Field.Required
}
func (ta TextArea) GetDescription() string {
  return ta.Field.Label
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
func (gi GenericInput) GetRequired() bool {
  return gi.Field.Required
}
func (gi GenericInput) GetDescription() string {
  return gi.Field.Label
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
func (fi FileInput) GetRequired() bool {
  return fi.Field.Required
}
func (fi FileInput) GetDescription() string {
  return fi.Field.Label
}

type SelectionGroup struct{
  // name is a base name
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
func (sg SelectionGroup) GetRequired() bool {
  return sg.Field.Required
}
func (sg SelectionGroup) GetDescription() string {
  return sg.Field.Label
}

// On the question of giving the Checkable field a Name for checkbox...
// I've decided that the UI will system will associate
type Checkable struct {
  Label string
  Value string
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
func (sg OptionGroup) GetRequired() bool {
  return sg.Field.Required
}
func (sg OptionGroup) GetDescription() string {
  return sg.Field.Label
}

type OptionItem struct {
  Label string
  Value string
}
