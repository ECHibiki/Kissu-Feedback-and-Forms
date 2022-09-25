package main

import (
    "testing"
    "os"
    "encoding/json"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/builder"
    "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
)

// TODO ONE: Correct any existing syntax errors
// TODO Two: Restory invalid structure tests, new func
// TODO Three: Create invalid character tests, new func
// TODO Four: Create duplicate field tests, new func

// TODO FIVE: Create duplicate forms of the same name



func TestConversionBetweenTypesAndInput(t *testing.T){
  var potential_inputs_InputType []string = []string{"color", "Date", "EMAIL", "nUMBER", "it's-invalid"}
  var correct_outputs_InputType []former.InputType = []former.InputType{former.Color , former.Date , former.Email, former.Number, former.Text}
  var valid_check_InputType []bool = []bool{  true ,  true ,   true,   true,   false}
  for index, _ := range valid_check_InputType {
    input_type , valid := former.InputTypeFromString(potential_inputs_InputType[index])
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
    input_type , valid := former.SelectionCategoryFromString(potential_inputs_SelectionCategory[index])
    if valid != valid_check_SelectionCategory[index]{
      t.Error("Validity of check was incorect" , potential_inputs_SelectionCategory[index] , valid , valid_check_SelectionCategory[index])
    } else if input_type != correct_outputs_SelectionCategory[index]{
      t.Error("Output type of check was incorect" , potential_inputs_SelectionCategory[index] , input_type , correct_outputs_SelectionCategory[index])
    }
  }

  var potential_inputs_FormObjectTag []string = []string{ "textarea", "genericinput", "fileinput", "selectiongroup", "optiongroup" , "invalid"}
  var correct_outputs_FormObjectTag []former.FormObjectTag = []former.FormObjectTag{ former.TextAreaTag , former.GenericInputTag , former.FileInputTag , former.SelectionGroupTag , former.OptionGroupTag ,  former.GenericInputTag }
  var valid_check_FormObjectTag []bool = []bool{  true ,  true ,   true,   true, true,  false}
  for index, _ := range valid_check_FormObjectTag {
    input_type , valid := former.FormObjectTagFromString(potential_inputs_FormObjectTag[index])
    if valid != valid_check_FormObjectTag[index]{
      t.Error("Validity of check was incorect" , potential_inputs_FormObjectTag[index] , valid , valid_check_FormObjectTag[index])
    } else if input_type != correct_outputs_FormObjectTag[index]{
      t.Error("Output type of check was incorect" , potential_inputs_FormObjectTag[index] , input_type , correct_outputs_FormObjectTag[index])
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
  demo_form_assumed_storage_name := "__Test_form_1"
  demo_form_name := "../Test form 1"
  var demo_form former.FormConstruct = former.FormConstruct{
      FormName: demo_form_name ,
      ID: "character-safe-form1",
      Description: "First test form",
      AnonOption: false,
      FormFields:[]former.FormGroup{
        {
          Label:"test-group1",
          ID: "test-group1",
          Description: "Groups and subgroups may have a description, when set it does not need respondables",
          // SubGroup: []former.FormGroup{},
          Respondables:[]former.UnmarshalerFormObject{
              {
                Type: former.TextAreaTag ,
                Object: former.TextArea{
                  Field: former.Field{
                    Label:"Test-Text-Area",
                    Name:"Test-TA",
                    Required:false,
                  },
                  Placeholder:"This is a test TA",
                },
              } ,
              {
                Type: former.GenericInputTag ,
                Object: former.GenericInput{
                  Field: former.Field{
                    Label:"Test-GenericInput",
                    Name:"Test-GI",
                    Required:true,
                  },
                  Placeholder:"This is a test GI",
                  Type:former.Text, // former.InputType
                },
              } ,
              {
                Type: former.FileInputTag ,
                Object: former.FileInput{
                  Field: former.Field{
                    Label:"Test-FileInput",
                    Name:"Test-FI",
                    Required:false,
                  },
                  AllowedExtRegex:".*",
                  MaxSize:10000000,
                },
              } ,
              {
                Type: former.SelectionGroupTag ,
                Object: former.SelectionGroup{
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
              },
              {
                Type: former.SelectionGroupTag,
                Object: former.SelectionGroup{
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
              },
              {
                Type: former.OptionGroupTag,
                Object: former.OptionGroup{
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
      },
  }

  issue_array := former.ValidateForm(demo_form)
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

  // do it again to garuntee multiple identical forms will not go through
  err = tools.StoreFormToDB(db, insertable_form)
  if err == nil{
    t.Fatal("Form was inserted twice, with same name, and passed without error")
  }

  var returned_form types.FormDBFields
  var rebuild_group former.FormConstruct
  returned_form , err = tools.GetFormOfID(db, 1)
  if returned_form.FieldJSON != string(marshaled_form_for_tests){
    t.Fatal("Fields returns from DB are not same as marshaled");
  }

  _ , err = tools.GetFormOfName(db, demo_form.StorageName())

  err = json.Unmarshal([]byte(returned_form.FieldJSON), &rebuild_group)
  if err != nil{
    panic(err)
  }
  marshaled_form_for_verify , err := json.Marshal(rebuild_group)
  if err != nil{
    panic(err)
  }
  if string(marshaled_form_for_verify) != string(marshaled_form_for_tests) {
    t.Fatal("unmarshalling process did not preserve data\n\n" , string(marshaled_form_for_verify) , "\n\n" , string(marshaled_form_for_tests));
  }




  err = builder.CreateFormDirectory(demo_form , cfg)
  if err != nil{
    panic(err)
  }
  safe_name := demo_form_assumed_storage_name
  _, err = os.Stat(initialization_folder + "/data/" + safe_name + "/")
  if err != nil {
    t.Fatal(initialization_folder + "/data/" + safe_name + "/" , err)
  }
}


func TestInvalidStructureForms(t *testing.T){
  var err error

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  var failing_demo_form former.FormConstruct = former.FormConstruct{
      AnonOption: false,
      FormName: "this should fail",
      ID: "invalid-structure-ID",
      Description: "The name of the form should be changed to remove the problem charactesr. It will fail for lacking form fields",
      FormFields:[]former.FormGroup{   },
  }

  failure_object := former.ValidateForm(failing_demo_form)
  if len(failure_object) != 1 {
    t.Fatal("The number of errors is incorrect" , failure_object)
  }
  if failure_object[0].FailType != former.HeadMissingMessage {
    t.Error("Error message is not recorded correctly" , failure_object)
  }
  if failure_object[0].FailCode != former.HeadMissingCode {
    t.Error("Error code is not recorded correctly" , failure_object)
  }

  if failure_object[0].FailPosition != failing_demo_form.ID {
    t.Error("Head fail position not correct. Fail position is unexpected" , failure_object)
  }

  first_failing_id := "test-formgroup-Fail1"
  second_failing_id := "test-subgroup-Fail"
  var an_alternative_failing_demo_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "../allowed-name",
    ID: "invalid-structure-form2",
    Description: "Another failing form, the invalid name should be cleaned",
    FormFields:[]former.FormGroup{
      {
        Label:"test-formgroup1",
        ID:"test-formgroup1",
        Respondables:[]former.UnmarshalerFormObject{
            {
              Type: former.GenericInputTag ,
              Object: former.GenericInput{
                Field: former.Field{
                  Label:"Test-GenericInput-1",
                  Name:"Test-GI-1",
                  Required:true,
                },
                Placeholder:"This is a test GI-1",
                Type:former.Text, // former.InputType
              } ,
          },
       },
        SubGroup: []former.FormGroup{
          {
            Label:"test-subgroup1",
            ID: "test-subgroup1",
            Respondables:[]former.UnmarshalerFormObject{
              {
                Type: former.GenericInputTag ,
                Object:former.GenericInput{
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
      },
      {
        // Fail here because no respondables
       Label:"test-formgroup-Fail1",
       ID:first_failing_id,
       Respondables:[]former.UnmarshalerFormObject{        },
       SubGroup: []former.FormGroup{
         {
           // Also Fail here because no respondables
           Label:"test-subgroup-Fail",
           ID:second_failing_id,
           Respondables:[]former.UnmarshalerFormObject{        },
           SubGroup: []former.FormGroup{   },
         },
       },
     },
     {
        // This should pass because they have a description
       Label:"2test-formgroup-Fail2",
       ID:"test-formgroup-Fail2",
       Description: "A description lets a field get away with no items",
     },
    },
  }

  failure_object = former.ValidateForm(an_alternative_failing_demo_form)
  if len(failure_object) != 2 {
    t.Error("The number of errors is incorrect" , failure_object)
  } else{
    if failure_object[0].FailType != former.GroupMissingMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailCode != former.GroupMissingCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailPosition != first_failing_id {
      t.Error("Error fail is not in correct location"  , failure_object[0].FailPosition , first_failing_id)
    }

    if failure_object[1].FailType != former.GroupMissingMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[1].FailCode != former.GroupMissingCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[1].FailPosition != second_failing_id {
      t.Error("Error fail is not in correct location"  , failure_object[1].FailPosition , second_failing_id)
    }
  }


}

func TestDuplicateIDNameFailure(t *testing.T){
  var err error

  // TESTS HEAD ID SAME AS SUBGROUP ID
  // TESTS SUBGROUP ID SAME AS SUBGROUP ID
  // TESTS FIELD NAME SAME AS FIELD NAME

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  var failing_duplicates_headtosub_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "this should fail",
    ID: "this-conflicts",
    Description: "This has an ID conflict. Potentially this could create issues with Javascript DOM methods",
    FormFields:[]former.FormGroup{
      {
         // This should pass because they have a description
        Label:"Forms with no entries can exist if they have a description, though this could allow for functionless forms there could be some uses...",
        ID:"this-conflicts",
        Description: "A description lets a field get away with no items",
      },
    },
  }

  failure_object := former.ValidateForm(failing_duplicates_headtosub_form)
  if len(failure_object) != 1 {
    t.Error("A form with an  duplicate ID is throwing incorrect number of errors" , failure_object)
  } else{
    if failure_object[0].FailType != former.DuplicateIDMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailCode != former.DuplicateIDCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailPosition != failing_duplicates_headtosub_form.ID {
      t.Error("Error fail is not in correct location"  , failure_object[0].FailPosition , failing_duplicates_headtosub_form.ID)
    }
  }


  duplicate_id_A := "matching-group-to-sub-id"
  duplicate_id_B := "matching-sub-to-sub-id"
  duplicate_id_C := "same-level-matching-sub-to-sub-id"
  var failing_duplicates_subtosub_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "this should fail",
    ID: "no-problem",
    Description: "This has an ID conflict. Potentially this could create issues with Javascript DOM methods",
    FormFields:[]former.FormGroup{
      {
         // This should pass because they have a description
        Label:"Forms with no entries can exist if they have a description, though this could allow for functionless forms there could be some uses...",
        ID:duplicate_id_A,
        Description: "A description lets a field get away with no items",
        SubGroup: []former.FormGroup{
          {
            // This should pass because they have a description
           Label:"...",
           ID:duplicate_id_B,
           Description: "Cross Subgroup check",
         },
          {
            // This should pass because they have a description
           Label:"...",
           ID:duplicate_id_C,
           Description: "same Subgroup check",
         },
          {
            // This should pass because they have a description
           Label:"...",
           ID:duplicate_id_C,
           Description: "same Subgroup check",
         },
        },
      },
      {
        // This should pass because they have a description
       Label:"...",
       ID:"safe-id",
       Description: "A bit more isolated",
       SubGroup: []former.FormGroup{
         {
           // This should pass because they have a description
          Label:"...",
          ID:duplicate_id_A,
          Description: "Cross Subgroup check",
        },
         {
           // This should pass because they have a description
          Label:"...",
          ID:duplicate_id_B,
          Description: "Cross Subgroup check",
        },
       },
     },
    },
  }

  // Follows outer leafs of search tree
  failure_object = former.ValidateForm(failing_duplicates_subtosub_form)
  if len(failure_object) != 3 {
    t.Error("The number of errors is incorrect" , failure_object)
  } else{
    if failure_object[0].FailType != former.DuplicateIDMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailCode != former.DuplicateIDCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailPosition !=   duplicate_id_A{
      t.Error("Error fail is not in correct location"  , failure_object[0].FailPosition , duplicate_id_A)
    }

    if failure_object[1].FailType != former.DuplicateIDMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[1].FailCode != former.DuplicateIDCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[1].FailPosition !=   duplicate_id_C{
      t.Error("Error fail is not in correct location"  , failure_object[1].FailPosition , duplicate_id_C)
    }

    if failure_object[2].FailType != former.DuplicateIDMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[2].FailCode != former.DuplicateIDCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[2].FailPosition !=   duplicate_id_B{
      t.Error("Error fail is not in correct location"  , failure_object[2].FailPosition , duplicate_id_B)
    }
  }



  duplicate_name_field := "duplicate-name"
  var failing_duplicates_fieldtofield_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "this should fail",
    ID: duplicate_name_field,
    Description: "There should be no issue with a name and ID being shared",
    FormFields:[]former.FormGroup{
      {
         // This should pass because they have a description
        Label:"Fields can't have the same name",
        ID:"no-problem",
        Description: "A description lets a field get away with no items",
        Respondables: []former.UnmarshalerFormObject{
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-1",
                Name:duplicate_name_field,
                Required:true,
              },
              Placeholder:"This is a test GI-1",
              Type:former.Text, // former.InputType
            } ,
          },
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-2",
                Name:duplicate_name_field,
                Required:false,
              },
              Placeholder:"This is a test GI-2",
              Type:former.Text, // former.InputType
            } ,
          },
        },
      },
    },
  }

  failure_object = former.ValidateForm(failing_duplicates_fieldtofield_form)
  if len(failure_object) != 1 {
    t.Error("A respondable with duplicate Name is throwing incorrect error numbers" , failure_object)
  } else{
    if failure_object[0].FailType != former.DuplicateNameMessage {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailCode != former.DuplicateNameCode {
      t.Error("Error message is not recorded correctly" , err)
    }
    if failure_object[0].FailPosition != duplicate_name_field {
      t.Error("Error fail is not in correct location"  , failure_object[0].FailPosition , duplicate_name_field)
    }
  }
}

func TestInvalidCharIDNameFailure(t *testing.T){
  // var err error

  // TESTS Construct and Group ID USES INVALID CHARACTESR FOR HTML ID
  // TESTS FIELD USES NAME WITH INVALID CHARACTESR FOR Name

  // Another Gin function builds the struct so that these functions can read it
  // function won't be tested because I don't want to mock HTTP requests at this time
  var failing_bad_char_form former.FormConstruct = former.FormConstruct{
    AnonOption: false,
    FormName: "this form name is altered to be placed into SQL and files",
    ID: "an id cannot have white spaces or characters that are not - or _ or : or .",
    Description: "ID and NAME tokens must begin with a letter ([A-Za-z]) and may be followed by any number of letters, digits ([0-9]), hyphens (\"-\"), underscores (\"_\"), colons (\":\"), and periods (\".\")",
    FormFields:[]former.FormGroup{
      {
         // This should pass because they have a description
        Label:"ID validity determined by alphanumeric checks",
        ID:"99starting-with-number-is-bad",
        Description: "",
        Respondables: []former.UnmarshalerFormObject{
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-1",
                Name: ":-names follow the same rules",
                Required:true,
              },
              Placeholder:"This is a test GI-1",
              Type:former.Text, // former.InputType
            } ,
          },
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-2",
                Name:"this:works:for:some:reason...why-though...",
                Required:false,
              },
              Placeholder:"This is a test GI-2",
              Type:former.Text, // former.InputType
            } ,
          },
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-2",
                Name:"this of course fails...",
                Required:false,
              },
              Placeholder:"This is a test GI-2",
              Type:former.Text, // former.InputType
            } ,
          },
          {
            Type: former.GenericInputTag ,
            Object: former.GenericInput{
              Field: former.Field{
                Label:"Test-GenericInput-2",
                Name:"::might as well show a compound error",
                Required:false,
              },
              Placeholder:"This is a test GI-2",
              Type:former.Text, // former.InputType
            } ,
          },
        },
        SubGroup: []former.FormGroup{
          {
            // This should pass because they have a description
           Label:"...",
           ID:"::one more for luck",
           Description: "Cross Subgroup check",
         },
       },
      },
    },
  }
  failure_object := former.ValidateForm(failing_bad_char_form)
  if len(failure_object) != 9 {
    t.Error("The number of errors is incorrect for inccorect character checks\n" , failure_object , len(failure_object))
  } else{
    if failure_object[0].FailType != former.InvalidIDCharactersMessage {
      t.Error("Error message is not recorded correctly" , failure_object[0])
    }
    if failure_object[0].FailCode != former.InvalidIDCharactersCode {
      t.Error("Error message is not recorded correctly" , failure_object[0])
    }
    if failure_object[0].FailPosition != "an id cannot have white spaces or characters that are not - or _ or : or ."{
      t.Error("Error fail is not in correct location"  , failure_object[0].FailPosition , "an id cannot have white spaces or characters that are not - or _ or : or .")
    }

    if failure_object[1].FailType != former.InvalidIDStarterMessage {
      t.Error("Error message is not recorded correctly" , failure_object[1])
    }
    if failure_object[1].FailCode != former.InvalidIDStarterCode {
      t.Error("Error message is not recorded correctly" , failure_object[1])
    }
    if failure_object[1].FailPosition != "99starting-with-number-is-bad"{
      t.Error("Error fail is not in correct location"  , failure_object[1].FailPosition , "99starting-with-number-is-bad")
    }

    if failure_object[2].FailType != former.InvalidIDStarterMessage {
      t.Error("Error message is not recorded correctly" , failure_object[2])
    }
    if failure_object[2].FailCode != former.InvalidIDStarterCode {
      t.Error("Error message is not recorded correctly" , failure_object[2])
    }
    if failure_object[2].FailPosition != "::one more for luck"{
      t.Error("Error fail is not in correct location"  , failure_object[2].FailPosition , "::one more for luck")
    }
    if failure_object[3].FailType != former.InvalidIDCharactersMessage {
      t.Error("Error message is not recorded correctly" , failure_object[3])
    }
    if failure_object[3].FailCode != former.InvalidIDCharactersCode {
      t.Error("Error message is not recorded correctly" , failure_object[3])
    }
    if failure_object[3].FailPosition != "::one more for luck"{
      t.Error("Error fail is not in correct location"  , failure_object[3].FailPosition , "::one more for luck")
    }

    if failure_object[4].FailType != former.InvalidNameStarterMessage {
      t.Error("Error message is not recorded correctly" , failure_object[4])
    }
    if failure_object[4].FailCode != former.InvalidNameStarterCode {
      t.Error("Error message is not recorded correctly" , failure_object[4])
    }
    if failure_object[4].FailPosition !=  ":-names follow the same rules"{
      t.Error("Error fail is not in correct location"  , failure_object[4].FailPosition , ":-names follow the same rules")
    }

    if failure_object[5].FailType != former.InvalidNameCharactersMessage {
      t.Error("Error message is not recorded correctly" , failure_object[5])
    }
    if failure_object[5].FailCode != former.InvalidNameCharactersCode {
      t.Error("Error message is not recorded correctly" , failure_object[5])
    }
    if failure_object[5].FailPosition !=  ":-names follow the same rules"{
      t.Error("Error fail is not in correct location"  , failure_object[4].FailPosition , ":-names follow the same rules")
    }

    if failure_object[6].FailType != former.InvalidNameCharactersMessage{
      t.Error("Error message is not recorded correctly" , failure_object[6])
    }
    if failure_object[6].FailCode != former.InvalidNameCharactersCode {
      t.Error("Error message is not recorded correctly" , failure_object[6])
    }
    if failure_object[6].FailPosition !=  "this of course fails..."{
      t.Error("Error fail is not in correct location"  , failure_object[6].FailPosition , "this of course fails...")
    }

    if failure_object[7].FailType != former.InvalidNameStarterMessage {
      t.Error("Error message is not recorded correctly" , failure_object[7])
    }
    if failure_object[7].FailCode != former.InvalidNameStarterCode {
      t.Error("Error message is not recorded correctly" , failure_object[7])
    }
    if failure_object[7].FailPosition !=  "::might as well show a compound error"{
      t.Error("Error fail is not in correct location"  , failure_object[7].FailPosition , "::might as well show a compound error")
    }

    if failure_object[8].FailType != former.InvalidNameCharactersMessage{
      t.Error("Error message is not recorded correctly" , failure_object[8])
    }
    if failure_object[8].FailCode != former.InvalidNameCharactersCode {
      t.Error("Error message is not recorded correctly" , failure_object[8])
    }
    if failure_object[8].FailPosition !=  "::might as well show a compound error"{
      t.Error("Error fail is not in correct location"  , failure_object[8].FailPosition , "::might as well show a compound error")
    }
  }



}
