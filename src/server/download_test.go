package main

import (
	"os"
	"testing"
	// "fmt"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
)

func TestCompressionOfGivenForm(t *testing.T) {
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
	prebuilder.ReplyToForm(1, second_store, "192.168.1.3", db, initialization_folder)
	prebuilder.ReplyToForm(2, first_store, "192.168.1.1", db, initialization_folder)
	prebuilder.ReplyToForm(2, first_store, "192.168.1.2", db, initialization_folder)

	// a CSV of all replies for given form data, allowing for Excel viewing
	// Placed into the form's directory
	// Steps: Get form for all Columns from FormConstruct and initialize a map[string][]string, get all rows and fill into the map
	//        In order defined by FormConstruct, createa a [][]sring for CSV creation

	err = returner.CreateInstancedCSVForGivenForm(db, 2, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store + "/data.csv")
	if err != nil {
		t.Fatal(err)
	}
	err = returner.CreateReadmeForGivenForm(db, 2, initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store + "/field-descriptors.json")
	if err != nil {
		t.Fatal(err)
	}
	// A tar.gz file containing the CSV, as it has zipped the entire form directory together
	err = returner.CreateDownloadableForGivenForm( second_store , initialization_folder)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(initialization_folder + "/data/" + second_store + "/" + "downloadable.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	return
	// serving of files by http done without tests..
}
