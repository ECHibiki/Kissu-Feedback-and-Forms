package former

type InputType string
type SelectionCategory string
type FormValidationError string

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
  GroupMissingError   FormValidationError    = "Form group does not contain any items or a description. Add one of these to submit."
)

type FailureObject struct {
  FailType FormValidationError
  FailPosition int64
}


type FormConstruct struct {
  FormName string
  Description string
  AnonOption bool
  FormFields []FormGroup
}

type FormGroup struct {
  Label string
  Description string
  SubGroup []FormGroup
  Respondables []FormObject
}

//currently a placeholder to gain polymorphic properties
type FormObject interface{
    ElementType() string
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
type GenericInput struct{
  Field Field
  Placeholder string
  Type InputType
}

type FileInput struct{
  Field Field
  AllowedExtRegex string
  MaxSize int64
}
type SelectionGroup struct{
  Field Field
  SelectionCategory SelectionCategory
  CheckableItems []Checkable
}
type Checkable struct {
  Label string
}

type OptionGroup struct{
  Field Field
  Options []OptionItem
}

type OptionItem struct {
  Label string
  Value string
}

func (ta TextArea) ElementType() string {
  return "TEXTAREA"
}
func (gi GenericInput) ElementType() string {
  return "INPUT"
}
func (fi FileInput) ElementType() string {
  return "INPUT"
}
func (sg SelectionGroup) ElementType() string {
  return "INPUT"
}
func (sg OptionGroup) ElementType() string {
  return "OPTION"
}
