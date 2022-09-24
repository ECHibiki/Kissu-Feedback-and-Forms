package main

import (
    "testing"
    "os"
    "strings"
    "encoding/json"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/builder"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
)


func TestConversionBetweenTypesAndInput(t *testing.T){
  var potential_inputs_InputType []string = []string{"color", "Date", "EMAIL", "nUMBER", "it's-invalid"}
  var correct_outputs_InputType []former.InputType = []former.InputType{former.Color , former.Date , former.Email, former.Number, former.Text}
  var valid_check_InputType []bool = []bool{  true ,  true ,   true,   true,   false}
  for index, _ := range valid_check_InputType {
    input_type , valid := builder.InputTypeFromString(potential_inputs_InputType[index])
    if valid != valid_check_InputType[index]{
      t.Error("Validity of check was incorect" , potential_inputs_InputType[index] , valid , valid_check_InputType[index])
    } else if input_type != correct_outputs_InputType[index]{
      t.Error("Output type of check was incorect" , potential_inputs_InputType[index] , input_type , correct_outputs_InputType[index])
    }
  }

  var potential_inputs_SelectionCategory []string = []string{ "checkbox", "radio", "Checkbox", "rADIO", "invalid"}
  var correct_outputs_SelectionCategory []former.SelectionCategory = []former.SelectionCategory{former.Checkbox , former.Radio , former.Checkbox, former.Radio, former.Checkbox}
  var valid_check_SelectionCategory []bool = []bool{  true ,  true ,   true,   true,   false}
  for index, _ := range valid_check_SelectionCategory {
    input_type , valid := builder.SelectionCategoryFromString(potential_inputs_SelectionCategory[index])
    if valid != valid_check_SelectionCategory[index]{
      t.Error("Validity of check was incorect" , potential_inputs_SelectionCategory[index] , valid , valid_check_SelectionCategory[index])
    } else if input_type != correct_outputs_SelectionCategory[index]{
      t.Error("Output type of check was incorect" , potential_inputs_SelectionCategory[index] , input_type , correct_outputs_SelectionCategory[index])
    }
  }
}

func TestFormMake(t *testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  var demo_form former.FormConstruct = former.FormConstruct{
      FormName: "../Test form 1",
      Description: "First test form",
      AnonOption: false,
      FormFields:[]former.FormGroup{
        {
          Label:"test-group1",
          Description: "Groups and subgroups may have a description, when set it does not need respondables",
          SubGroup: []former.FormGroup{},
          Respondables:[]former.FormObject{
              former.TextArea{
                Field: former.Field{
                  Label:"Test-Text-Area",
                  Name:"Test-TA",
                  Required:false,
                },
                Placeholder:"This is a test TA",
              } ,
              former.GenericInput{
                Field: former.Field{
                  Label:"Test-GenericInput",
                  Name:"Test-GI",
                  Required:true,
                },
                Placeholder:"This is a test GI",
                Type:former.Text, // former.InputType
              } ,
              former.FileInput{
                Field: former.Field{
                  Label:"Test-FileInput",
                  Name:"Test-FI",
                  Required:false,
                },
                AllowedExtRegex:".*",
                MaxSize:10000000,
              } ,
              former.SelectionGroup{
                Field: former.Field{
                  Label:"Test-Chk-SelectGroup",
                  Name:"Test-Chk-SG",
                  Required:true,
                },
                SelectionCategory: former.Checkbox,
                CheckableItems:[]former.Checkable{
                  {Label:"A check Item"},
                  {Label:"Another check Item"},
                },
              },
              former.SelectionGroup{
                Field: former.Field{
                  Label:"Test-rdo-SelectGroup",
                  Name:"Test-rdo-SG",
                  Required:true,
                },
                SelectionCategory: former.Radio,
                CheckableItems:[]former.Checkable{
                  {Label:"A radio Item"},
                  {Label:"Another radio Item"},
                },
              },
              former.OptionGroup{
                Field: former.Field{
                  Label:"Test-optGrp",
                  Name:"Test-optGrp",
                  Required:true,
                },
                Options:[]former.OptionItem{
                  {
                    Label:"Item 1",
                    Value: "item-1",
                  } ,
                  {
                    Label:"Item 2",
                    Value: "item-2",
                  } ,
                },
              },
          },
        },
      },
  }

  issue_array := builder.ValidateForm(demo_form)
  if len(issue_array) != 0 {
    t.Fatal(issue_array)
  }

  // see use-case doc for other things to add

  marshaled_form_for_tests , err := json.Marshal(demo_form)
  if err != nil{
    panic(err)
  }
  var insertable_form types.FormDBFields
  insertable_form , err =  builder.MakeFormWritable(demo_form)
  if err != nil{
    panic(err)
  }
  if string(marshaled_form_for_tests) != insertable_form.FieldJSON{
    t.Fatal("Form marshaling did not render an expected result")
  }
  if insertable_form.ID != 0{
    t.Fatal("Form construction set a value for ID")
  }
  if insertable_form.UpdatedAt == 0{
    t.Fatal("Form construction did not set an updated time")
  }
  err = tools.StoreFormToDB(db, insertable_form)
  if err != nil{
    panic(err)
  }
  var returned_form types.FormDBFields
  var rebuild_group former.FormConstruct
  returned_form , err = tools.GetFormOfID(db, 1)
  if returned_form.FieldJSON != string(marshaled_form_for_tests){
    t.Fatal("Fields returns from DB are not same as marshaled");
  }
  err = json.Unmarshal([]byte(returned_form.FieldJSON), &rebuild_group)
  if err != nil{
    panic(err)
  }

  err = builder.CreateFormDirectory(demo_form , cfg)
  if err != nil{
    panic(err)
  }
  safe_name := strings.ReplaceAll(demo_form.FormName , "." , "-")
  safe_name = strings.ReplaceAll(safe_name , "/" , "_")

  _, err = os.Stat(initialization_folder + "/data/" + demo_form.FormName + "/")
  if err != nil {
    t.Fatal(initialization_folder + "/data/" + demo_form.FormName + "/")
  }
  _, err = os.Stat(initialization_folder + "/data/" + demo_form.FormName + "/files")
  if err != nil {
    t.Fatal(initialization_folder + "/data/" + demo_form.FormName + "/files")
  }
}


func TestInvalidForms(t *testing.T){
  var err error

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  var failing_demo_form former.FormConstruct = former.FormConstruct{
      AnonOption: false,
      FormName: "../../www/this should fail",
      Description: "The name of the form should be changed to remove the problem charactesr. It will fail for lacking form fields",
      FormFields:[]former.FormGroup{   },
  }

  failure_object := builder.ValidateForm(failing_demo_form)
  if len(failure_object) == 0 {
    t.Fatal("A form with errors is passing")
  }
  if len(failure_object) < 1 {
    t.Fatal("Some errors are missing" , failure_object)
  }
  if failure_object[0].FailType != former.GroupMissingError {
    t.Error("Error message is not recorded correctly" , failure_object)
  }
  if failure_object[0].FailPosition != 0 {
    t.Error("Fail position is unexpected" , failure_object)
  }

  var an_alternative_failing_demo_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "../invalid-name",
    Description: "Another failing form, the invalid name should be cleaned",
    FormFields:[]former.FormGroup{
      {
        Label:"test-formgroup1",
        Respondables:[]former.FormObject{
            former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-1",
                Name:"Test-GI-1",
                Required:true,
              },
              Placeholder:"This is a test GI-1",
              Type:former.Text, // former.InputType
            } ,
       },
        SubGroup: []former.FormGroup{
          {
            Label:"test-subgroup1",
            Respondables:[]former.FormObject{
              former.GenericInput{
                Field: former.Field{
                  Label:"Test-GenericInput-2",
                  Name:"Test-GI-2",
                  Required:true,
                },
                Placeholder:"This is a test GI-2",
                Type:former.Text, // former.InputType
              } ,
            },
          },
        },
      },
      {
        // Fail here because no respondables
       Label:"test-formgroup-Fail",
       Respondables:[]former.FormObject{        },
       SubGroup: []former.FormGroup{
         {
           // Also Fail here because no respondables
           Label:"test-subgroup-Fail",
           Respondables:[]former.FormObject{        },
           SubGroup: []former.FormGroup{   },
         },
       },
     },
     {
        // This should pass because they have a description
       Label:"test-formgroup-Fail",
       Description: "A description lets a field get away with no items",
     },
    },
  }

  failure_object = builder.ValidateForm(an_alternative_failing_demo_form)
  if len(failure_object) == 0 {
    t.Error("A form with an empty subgroup is passing")
  }
  if len(failure_object) < 2 {
    t.Error("Some error messages are missing" , err)
  }
  if failure_object[0].FailType != former.GroupMissingError {
    t.Error("Error message is not recorded correctly" , err)
  }
  if failure_object[0].FailPosition != 3 {
    t.Error("Error fail is not in correct location" , err)
  }
  if failure_object[0].FailType != former.GroupMissingError {
    t.Error("Error message is not recorded correctly" , err)
  }
  if failure_object[0].FailPosition != 4 {
    t.Error("Error fail is not in correct location" , err)
  }

}
