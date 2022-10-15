package main

import (
	"encoding/json"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/destroyer"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"os"
	"testing"
)

func TestDeleteOfForm(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	first_name := "Test form 1"
	first_store := "Test_form_1"
	prebuilder.DoFormInitialization(first_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	second_name := "Test form 2"
	second_store := "Test_form_2"
	prebuilder.DoFormInitialization(second_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	prebuilder.ReplyToForm(1, first_store, "192.168.1.1", db, initialization_folder)
	prebuilder.ReplyToForm(2, second_store, "192.168.1.2", db, initialization_folder)
	prebuilder.ReplyToForm(2, second_store, "192.168.1.3", db, initialization_folder)

	destroyer.DeleteForm(db, second_store, 2)

	r, err := returner.GetResponseByID(db, 3)
	if err == nil {
		t.Fatal("A response to a delte form lingers", r)
	}

	f, err := returner.GetFormOfID(db, 2)
	if err == nil {
		t.Fatal("deleted form still exists", f)
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store)
	if err != nil {
		t.Fatal("Form deletes should retain old files in case of mistakes ", second_store, " missing")
	}

	f, err = returner.GetFormOfID(db, 1)
	if err != nil {
		t.Fatal("Error on form that should still exist", r)
	}
	_, err = os.Stat(initialization_folder + "/data/" + first_store)
	if err != nil {
		t.Fatal("An unrelated form directory was removed", first_store)
	}
}

func TestDeleteOfResponse(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	first_name := "Test form 1"
	first_store := "Test_form_1"
	prebuilder.DoFormInitialization(first_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	second_name := "Test form 2"
	second_store := "Test_form_2"
	prebuilder.DoFormInitialization(second_name, "a-simple-identifier", db, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	prebuilder.ReplyToForm(1, first_store, "192.168.1.1", db, initialization_folder)
	prebuilder.ReplyToForm(2, second_store, "192.168.1.2", db, initialization_folder)
	prebuilder.ReplyToForm(2, second_store, "192.168.1.3", db, initialization_folder)

	destroyer.DeleteResponse(db, initialization_folder, 2, second_store, "192.168.1.2")

	_, err = returner.GetFormOfID(db, int64(1))
	if err != nil {
		t.Error("Error on form that should still exist")
	}
	_, err = os.Stat(initialization_folder + "/data/" + first_store)
	if err != nil {
		t.Error("An unrelated form directory was removed")
	}
	_, err = returner.GetFormOfID(db, int64(2))
	if err != nil {
		t.Error("Error on form that should still exist")
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store)
	if err != nil {
		t.Error("An unrelated form directory was removed")
	}

	_, err = returner.GetResponseByID(db, 1)
	if err != nil {
		t.Error("Error on response that should still exist")
	}
	_, err = os.Stat(initialization_folder + "/data/" + first_store + "/192.168.1.1/")
	if err != nil {
		t.Error("An unrelated Response directory was removed")
	}
	_, err = returner.GetResponseByID(db, 2)
	if err == nil {
		t.Error("Error on response that should not exist")
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store + "/192.168.1.2/")
	if err == nil {
		t.Error("A Response directory was not removed")
	}
	_, err = returner.GetResponseByID(db, 3)
	if err != nil {
		t.Error("Error on response that should still exist")
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store + "/192.168.1.3/")
	if err != nil {
		t.Error("An unrelated Response directory was removed")
	}

}

func UndoTest(t *testing.T) {
	var initialization_folder string = "../../test"
	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)
	first_name := "Test form 1"
	first_store := "Test_form_1"
	prebuilder.DoFormInitialization(first_name, "a-simple-identifier", db, initialization_folder)

	destroyer.UndoForm(db, first_store, initialization_folder)
	_, err := returner.GetFormOfID(db, int64(1))
	if err == nil {
		t.Error("Form should be removed")
	}
	_, err = os.Stat(initialization_folder + "/data/" + first_store + "/")
	if err == nil {
		t.Error("Directory should be removed")
	}
}
func DirectoryUndoTest(t *testing.T) {
	var initialization_folder string = "../../test"
	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)
	first_name := "Test form 1"
	first_store := "Test_form_1"
	prebuilder.DoFormInitialization(first_name, "a-simple-identifier", db, initialization_folder)
	f, err := returner.GetFormOfID(db, int64(1))
	var form_construct former.FormConstruct
	json.Unmarshal([]byte(f.FieldJSON), &form_construct)

	destroyer.UndoFormDirectory(form_construct, initialization_folder)

	_, err = returner.GetFormOfID(db, int64(1))
	if err != nil {
		t.Error("Form should not be removed")
	}
	_, err = os.Stat(initialization_folder + "/data/" + first_store + "/")
	if err == nil {
		t.Error("Directory should be removed")
	}
}

func ResponseUndoTest(t *testing.T) {
	t.Fatal("unimplemented")
}
