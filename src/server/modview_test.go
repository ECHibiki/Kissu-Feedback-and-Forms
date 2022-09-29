package main
import(
  "testing"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/returner"
  // "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
)

func TestListAllForms( t *testing.T ){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  demo_form_id_check_name := "Test form 1"
  tools.DoFormInitialization(demo_form_id_check_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  second_name := "Test form 2"
  tools.DoFormInitialization(second_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }

  form1, err := tool.GetFormOfID(1)
  form2, err := tool.GetFormOfID(2)
  forms_1and2 := []types.FormDBFields{form1, form2}

  var forms []types.FormDBFields = returner.GetAllForms()

  forms_test_json , err  := json.Marshal(forms_1and2)
  if err != nil {
    t.Fatal(err)
  }
  forms_json , err := json.Marshal(forms)
  if err != nil {
    t.Fatal(err)
  }

  if strings(forms_test_json) != strings(forms_json){
    t.Fatal("Combined forms is lacking information")
  }
}

func TestListResponsesToForm(t *testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  demo_form_id_check_name := "Test form 1"
  demo_form_assumed_storage_name := "Test_form_1"
  tools.DoFormInitialization(demo_form_id_check_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  second_name := "Test form 2"
  second_store := "Test_form_2"
  tools.DoFormInitialization(second_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  testing.ReplyToForm(db , 1 , second_store, "192.168.1.3", db , initialization_folder)
  testing.ReplyToForm(db , 2 , demo_form_assumed_storage_name , "192.168.1.1" , db , initialization_folder)
  testing.ReplyToForm(db , 2 , demo_form_assumed_storage_name, "192.168.1.2", db , initialization_folder)

  reply_list := returner.GetRepliesToForm( 2 )
  list_json, err := json.Marshal(reply_list)

  reply_1 , _ := tools.GetResponseByID(db  , 2 )
  reply_2 , _ := tools.GetResponseByID(db  , 3 )
  replies_test := []types.ResponseDBFields{reply_1 , reply_2}
  list_test_json, err := json.Marshal(list_test)

  if strings(list_test_json) != strings(list_json){
    t.Fatal("Combined replies is lacking information")
  }
}

func TestDisplaySingleResponse(t *testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  demo_form_id_check_name := "Test form 1"
  demo_form_assumed_storage_name := "Test_form_1"
  tools.DoFormInitialization(demo_form_id_check_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  second_name := "Test form 2"
  second_store := "Test_form_2"
  tools.DoFormInitialization(second_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  testing.ReplyToForm(db , 1 , second_store, "192.168.1.3", db , initialization_folder)
  testing.ReplyToForm(db , 2 , demo_form_assumed_storage_name , "192.168.1.1" , db , initialization_folder)
  testing.ReplyToForm(db , 2 , demo_form_assumed_storage_name, "192.168.1.2", db , initialization_folder)

  // add in something important related to form responses to make it different from tool.GetResponseByID...
  reply := resonder.GetResponse(db  , 2)
  r, err := json.Marshal(reply)
  reply_test , _ := tools.GetResponseByID(db  , 2 )
  r_test, err := json.Marshal(reply_test)

  // Retrival should list full information  if strings(r) != strings(r_test){
    t.Fatal("Reply not identical")
  }
}
