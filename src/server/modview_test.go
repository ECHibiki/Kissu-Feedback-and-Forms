package main

import (
	"encoding/json"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	"testing"
	// "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/responder"
	// "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	// "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
)

func TestListAllForms(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	demo_form_id_check_name := "Test form 1"
	prebuilder.DoFormInitialization(demo_form_id_check_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	second_name := "Test form 2"
	prebuilder.DoFormInitialization(second_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}

	form1, err := tools.GetFormOfID(db, 1)
	form1.FieldJSON = ""
	form2, err := tools.GetFormOfID(db, 2)
	form2.FieldJSON = ""
	forms_1and2 := []types.FormDBFields{form1, form2}

	forms, err := returner.GetAllForms(db)
	if err != nil {
		t.Fatal(err)
	}

	forms_test_json, err := json.Marshal(forms_1and2)
	if err != nil {
		t.Fatal(err)
	}
	forms_json, err := json.Marshal(forms)
	if err != nil {
		t.Fatal(err)
	}

	if string(forms_test_json) != string(forms_json) {
		t.Fatal("Combined forms is lacking information")
	}
}

func TestListResponsesToForm(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	demo_form_id_check_name := "Test form 1"
	demo_form_assumed_storage_name := "Test_form_1"
	prebuilder.DoFormInitialization(demo_form_id_check_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	second_name := "Test form 2"
	second_store := "Test_form_2"
	prebuilder.DoFormInitialization(second_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	prebuilder.ReplyToForm(1, second_store, "192.168.1.3", db, initialization_folder)
	prebuilder.ReplyToForm(2, demo_form_assumed_storage_name, "192.168.1.1", db, initialization_folder)
	prebuilder.ReplyToForm(2, demo_form_assumed_storage_name, "192.168.1.2", db, initialization_folder)

	db_reply_list, err := returner.GetRepliesToForm(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	db_reply_list_json, err := json.Marshal(db_reply_list)
	if err != nil {
		t.Fatal(err)
	}

	reply_1, _ := tools.GetResponseByID(db, 2)
	reply_2, _ := tools.GetResponseByID(db, 3)
	replies_test := []types.ResponseDBFields{reply_2, reply_1}
	replies_test_json, err := json.Marshal(replies_test)

	if string(db_reply_list_json) != string(replies_test_json) {
		t.Fatal("Combined replies is lacking information\n", string(db_reply_list_json), "\n", string(replies_test_json))
	}
}

func TestDisplaySingleResponse(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	demo_form_id_check_name := "Test form 1"
	demo_form_assumed_storage_name := "Test_form_1"
	prebuilder.DoFormInitialization(demo_form_id_check_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	second_name := "Test form 2"
	second_store := "Test_form_2"
	prebuilder.DoFormInitialization(second_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	prebuilder.ReplyToForm(1, second_store, "192.168.1.3", db, initialization_folder)
	prebuilder.ReplyToForm(2, demo_form_assumed_storage_name, "192.168.1.1", db, initialization_folder)
	prebuilder.ReplyToForm(2, demo_form_assumed_storage_name, "192.168.1.2", db, initialization_folder)

	// add in something important related to form responses to make it different from tools.GetResponseByID...
	reply, err := returner.GetResponseByID(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	r, err := json.Marshal(reply)
	if err != nil {
		t.Fatal(err)
	}
	reply_test, err := tools.GetResponseByID(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	r_test, err := json.Marshal(reply_test)
	if err != nil {
		t.Fatal(err)
	}
	if string(r) != string(r_test) {
		t.Fatal("reply is lacking information")
	}
}
