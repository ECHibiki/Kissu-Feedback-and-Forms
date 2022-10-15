package main

import (
	"encoding/json"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"strings"
	"testing"
)

func TestRetrieval(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	// Another Gin function builds the struct so that these functions can read it
	// function won't be tested because I don't want to mock HTTP requests at this time
	demo_form_id_check_name := "../Test form 1"
	prebuilder.DoFormInitialization(demo_form_id_check_name, "a-simple-identifier", db, initialization_folder)

	demo_form_name_check_assumed_storage_name := "__alternative_test_form_1"
	demo_form_name_check_name := "../alternative test form 1"
	prebuilder.DoFormInitialization(demo_form_name_check_name, "a-simple-identifier", db, initialization_folder)

	// ---- Forget the initialization of fields

	var insertable_form_id, insertable_form_name former.FormConstruct

	insertable_form_id_db, err := returner.GetFormOfID(db, 1)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(insertable_form_id_db.FieldJSON), &insertable_form_id)
	insertable_form_name_db, err := returner.GetFormOfName(db, demo_form_name_check_assumed_storage_name)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(insertable_form_name_db.FieldJSON), &insertable_form_name)

	env := templater.ReturnTemplateHandler()
	// generics for outputting a template depending on the form ID
	// Output should be an html page replicating the effect of marshal on the struct
	template_id, err := returner.RenderTestingTemplate(db, env, initialization_folder, int64(1))
	if err != nil {
		panic(err)
	}
	form_id_marshal, _ := json.Marshal(insertable_form_id)
	template_id = strings.ReplaceAll(template_id, "\n", "")
	tilen := len(template_id) - len(insertable_form_id.FormFields[0].Respondables[0].GetType()+insertable_form_id.FormFields[0].Respondables[0].Object.GetName())
	form_id_str := (string(form_id_marshal)[:tilen] + insertable_form_id.FormFields[0].Respondables[0].GetType() + insertable_form_id.FormFields[0].Respondables[0].Object.GetName())

	if template_id != form_id_str {
		t.Error("Test template render by ID failed--\nCreation:", string(template_id), "\nAssmpton:", form_id_str)
	}

	template_name, err := returner.RenderTestingTemplate(db, env, initialization_folder, demo_form_name_check_assumed_storage_name)
	if err != nil {
		panic(err)
	}
	form_name_marshal, _ := json.Marshal(insertable_form_name)
	template_name = strings.ReplaceAll(template_name, "\n", "")
	tiname := len(template_name) - len(insertable_form_id.FormFields[0].Respondables[0].GetType()+insertable_form_id.FormFields[0].Respondables[0].Object.GetName())
	form_name_str := (string(form_name_marshal)[:tiname] + insertable_form_id.FormFields[0].Respondables[0].GetType() + insertable_form_id.FormFields[0].Respondables[0].Object.GetName())
	if template_name != form_name_str {
		t.Error("Test template render by Name failed--\nCreation:", string(template_name), "\nAssmpton:", form_name_str)
	}
}

func TestGetForm(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	demo_form_name := "Test form 1"
	demo_form_store_name := "Test_form_1"
	prebuilder.DoFormInitialization(demo_form_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	_, err = returner.GetFormByNameAndID(db, demo_form_store_name, 1)
	if err != nil {
		t.Fatal(err)
	}
}
