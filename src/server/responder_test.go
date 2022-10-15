package main

import (
	"encoding/json"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/destroyer"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/responder"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"io/ioutil"
	"mime/multipart"
	"os"
	"testing"
)

func TestInputStorage(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	// Another Gin function builds the struct so that these functions can read it
	// function won't be tested because I don't want to mock HTTP requests at this time
	demo_form_assumed_storage_name := "__Test_form_1"
	demo_form_name := "../Test form 1"
	prebuilder.DoFormInitialization(demo_form_name, "a-simple-identifier", db, initialization_folder)

	var files map[string]former.MultipartFile = prebuilder.CopyTestFilesToMemory(initialization_folder, map[string]string{"Test-FI": "test-file-1.jpg"})

	// populate a response struct that would be filled out in a route
	//
	demo_response := former.FormResponse{
		FormName:     demo_form_assumed_storage_name,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses: map[string]string{
			// Fill them out here
			"anon-option":   "true",
			"Test-TA":       "../some text\n\n\tasdf",
			"Test-GI":       "../some text",
			"Test-Chk-SG-1": "ck1", // in this case the check group has been assigned ...-1 from the algorithm
			"Test-Chk-SG-3": "ck3", // in this case the check group has been assigned ...-3 from the algorithm
			"Test-rdo-SG":   "rd1", // In this case the radio group is just called
			"Test-optGrp":   "item-2",
		},
		// File paths are distinct in that their data effects file storage.
		// While text based fields may be passed around and later inserted somewhere,
		// Files must be moved around the OS
		// This means giving them a unique identifier, in this case a JSON column
		// It does not need a database column because after validation the files are written to a predictable location
		FileObjects: files,
	}

	//Get a form
	form, err := tools.GetFormOfID(db, demo_response.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}

	// Validate responses
	// Check files are valid
	var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(demo_response.Responses, rebuild_group)
	var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(demo_response.FileObjects, rebuild_group)
	issue_array := append(text_issue_array, file_issue_array...)
	// It should pass so write to the propper locations
	if len(issue_array) > 0 {
		t.Fatal("There should be no errors here", issue_array)
	}

	// check to scramble IP
	demo_response.ScrambleResponderID()
	if demo_response.ResponderID == "192.168.1.1" {
		t.Fatal("ResponderID was unchanged even though set to anonymous")
	}
	if len(demo_response.ResponderID) <= len("192.168.1.1") {
		t.Fatal("Unexpected ResponderID", demo_response.ResponderID)
	}

	err = responder.CreateResponderFolder(initialization_folder, demo_response)
	if err != nil {
		t.Fatal("Failed to create response folder", err)
	}
	_, err = os.Stat(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/")
	if err != nil {
		t.Error(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/ is missing")
	}
	_, err = os.Stat(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/files/")
	if err != nil {
		t.Error(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/files/ is missing")
	}

	error_list := tools.WriteFilesFromMultipart(initialization_folder, demo_response)
	if len(error_list) != 0 {
		t.Fatal(error_list)
	}
	_, err = os.Stat(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/files/Test-FI-test-file-1.jpg")
	if err != nil {
		t.Error(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/files/Test-FI-test-file-1.jpg")
	}

	// A combination of Responses and File Locations listing a URL for file download where it will be served
	err = tools.WriteResponsesToJSONFile(initialization_folder, demo_response)
	if err != nil {
		t.Fatal("Writting JSON failed")
	}
	json_file, err := ioutil.ReadFile(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/responses.json")
	if err != nil {
		t.Fatal(initialization_folder + "/data/" + demo_response.FormName + "/" + demo_response.ResponderID + "/responses.json is missing")
	}

	var json_resp former.JSONFormResponse
	err = json.Unmarshal(json_file, &json_resp)
	if err != nil {
		t.Fatal("responses.json error during parse")
	}
	original_json_resp := tools.ConvertFormResponseToJSONFormResponse(initialization_folder, demo_response)
	testing_json_r, err := json.Marshal(json_resp)
	if err != nil {
		t.Fatal(err)
	}
	original_testing_json_r, err := json.Marshal(original_json_resp)
	if err != nil {
		t.Fatal(err)
	}
	if string(testing_json_r) != string(original_testing_json_r) {
		t.Error("Data was lost when writting json file")
	}

	demo_response_db_fields, err := responder.FormResponseToDBFormat(demo_response)
	if err != nil {
		t.Fatal(err)
	}
	// A combination of Responses and File Locations listing a URL for file download where it will be served
	err = tools.StoreResponseToDB(db, demo_response_db_fields)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := returner.GetResponseByID(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	resp.ID = demo_response_db_fields.ID
	a, _ := json.Marshal(demo_response_db_fields)
	b, _ := json.Marshal(resp)
	if string(a) != string(b) {
		t.Fatal("GetResponseByID: Retrieval from DB did not produce expected results\n\t", string(a), "\n\t", string(b))
	}
}

func TestInputRejection(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	// Another Gin function builds the struct so that these functions can read it
	// function won't be tested because I don't want to mock HTTP requests at this time
	demo_form_assumed_storage_name := "__Test_form_1"
	demo_form_name := "../Test form 1"
	prebuilder.DoFormInitialization(demo_form_name, "a-simple-identifier", db, initialization_folder)

	// The Value type should be that which FormFile produces
	var files map[string]former.MultipartFile = prebuilder.CopyTestFilesToMemory(initialization_folder, map[string]string{"Test-FI": "invalid-file.png"})

	// populate a response struct that would be filled out in a route
	// Every item here should fail

	demo_response_fail := former.FormResponse{
		FormName:     demo_form_assumed_storage_name,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses: map[string]string{
			"Test-TA": "   \n\n\t", // Required field is empty
			// "Test-GI": "", // Required field is missing from POST message
			"Test-Chk-SG-9999": "invalid",                          // Checkgroup numbering does not exist
			"Test-rdo-SG":      "asfdasdfasdfasdf1324123412341324", // radio group
			"Test-optGrp":      "[][][][][][][][][][][][][]",
		},
		FileObjects: files,
	}

	//Get a form
	form, err := tools.GetFormOfID(db, demo_response_fail.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}
	// Validate responses
	// Check files are valid
	var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(demo_response_fail.Responses, rebuild_group)
	// Won't be tested here as yet
	var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(demo_response_fail.FileObjects, rebuild_group)
	issue_array := append(text_issue_array, file_issue_array...)
	// It should pass so write to the propper locations
	if len(issue_array) != 8 {
		t.Error("Error count mismatched, wants 8 has ")
	}
	for i, e := range issue_array {
		var fm former.FormValidationError
		var fc former.FormErrorCode
		if e.FailPosition == "Test-TA" {
			fm = former.ResponseMissingMessage
			fc = former.ResponseMissingCode
		} else if e.FailPosition == "Test-GI" {
			fm = former.ResponseMissingMessage
			fc = former.ResponseMissingCode
		} else if e.FailPosition == "Test-Chk-SG" {
			fm = former.ResponseMissingMessage
			fc = former.ResponseMissingCode
		} else if e.FailPosition == "Test-Chk-SG-9999" {
			fm = former.InvalidInputMessage
			fc = former.InvalidInputCode
		} else if e.FailPosition == "Test-rdo-SG" {
			fm = former.InvalidSelectionValueMessage
			fc = former.InvalidSelectionValueCode

		} else if e.FailPosition == "Test-optGrp" {
			fm = former.InvalidOptionValueMessage
			fc = former.InvalidOptionValueCode
		} else if e.FailPosition == "Test-FI" {
			// AllowedExtRegex string
			// MaxSize int64
			if i == 6 {
				fm = former.InvalidFileSizeMessage
				fc = former.InvalidFileSizeCode
			} else if i == 7 {
				fm = former.InvalidFileExtMessage
				fc = former.InvalidFileExtCode
			} else {
				t.Error("Fail FI in wrong location")
			}
		}
		// Anon option is an optional flag
		/* else if e.FailPosition == "anon-option"{
		   if i != 8{
		     t.Error("Fail anon-option in wrong location")
		   }
		   fm = former.ResponseMissingMessage
		   fc = former.ResponseMissingCode
		 } */

		if fm != e.FailType {
			t.Error("Index", i, "Failtype missmatch on ", e, ", is ", fm)
		}
		if fc != e.FailCode {
			t.Error("Failcode missmatch on", e, ", is ", fc)
		}
	}
}

func TestIllegalFName(t *testing.T) {

	var initialization_folder string = "../../test"
	var err error

	db, _, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	// Another Gin function builds the struct so that these functions can read it
	// function won't be tested because I don't want to mock HTTP requests at this time
	demo_form_assumed_storage_name := "__Test_form_1"
	demo_form_name := "../Test form 1"
	prebuilder.DoFormInitialization(demo_form_name, "a-simple-identifier", db, initialization_folder)

	demo_response_fail := former.FormResponse{
		FormName:     demo_form_assumed_storage_name,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses:    map[string]string{},
		FileObjects: map[string]former.MultipartFile{
			"Test-FI": {
				File: nil,
				Header: &multipart.FileHeader{
					Filename: "/definetly/not/allowed/jpg",
					Header:   nil, // we won't use this unless we must
					Size:     0,
				},
			},
		},
	}
	//Get a form
	form, err := tools.GetFormOfID(db, demo_response_fail.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}

	var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(demo_response_fail.FileObjects, rebuild_group)
	if len(file_issue_array) == 0 {
		t.Fatal("File issue missing", file_issue_array)
	}
	if file_issue_array[0].FailType != former.DangerousPathMessage {
		t.Fatal("Fail message missmatch for illegal fname", file_issue_array[0], former.DangerousPathMessage)

	}
	if file_issue_array[0].FailCode != former.DangerousPathCode {
		t.Fatal("Fail code missmatch for illegal fname", file_issue_array[0], former.DangerousPathCode)
	}
}

func TestSamePersonReplyingTwice(t *testing.T) {
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

	prebuilder.ReplyToForm(1, first_store, "192.168.1.1", db, initialization_folder)
	var files map[string]former.MultipartFile = prebuilder.CopyTestFilesToMemory(initialization_folder, map[string]string{"Test-FI": "test-file-1.jpg"})
	new_response := former.FormResponse{
		FormName:     first_store,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses: map[string]string{
			// Fill them out here
			"anon-option":   "false",
			"Test-TA":       "../some text\n\n\tasdf",
			"Test-GI":       "new text",
			"Test-Chk-SG-1": "ck1", // in this case the check group has been assigned ...-1 from the algorithm
			"Test-Chk-SG-3": "ck3", // in this case the check group has been assigned ...-3 from the algorithm
			"Test-rdo-SG":   "rd1", // In this case the radio group is just called
			"Test-optGrp":   "item-2",
		},
		// File paths are distinct in that their data effects file storage.
		// While text based fields may be passed around and later inserted somewhere,
		// Files must be moved around the OS
		// This means giving them a unique identifier, in this case a JSON column
		// It does not need a database column because after validation the files are written to a predictable location
		FileObjects: files,
	}

	form, err := tools.GetFormOfID(db, new_response.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}
	// Validate responses
	// Check files are valid
	var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(new_response.Responses, rebuild_group)
	// Won't be tested here as yet
	var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(new_response.FileObjects, rebuild_group)
	issue_array := append(text_issue_array, file_issue_array...)
	if len(issue_array) != 0 {
		t.Fatal("There should be no errors")
	}

	// Edit mode should scramble the IP to see if it exists elsewhere
	edit_mode, old_id, err := responder.CheckIfEdit(db, new_response)

	if err != nil {
		t.Fatal(err)
	}
	if edit_mode == false {
		t.Fatal("Edit mode was not set")
	}
	if old_id != "192.168.1.1" {
		t.Fatal("We're not testing that yet")
	}

	oldresp, err := returner.GetResponseByID(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	old_resp_marshal, err := json.Marshal(oldresp)
	if err != nil {
		t.Fatal(err)
	}

	if edit_mode {
		// Deleting is important because the responder ID could be set to scramble
		// Also easier and since nothing relies on the data it can be done
		destroyer.UndoResponse(db, new_response, old_id, initialization_folder)
	}

	_, err = os.Stat(initialization_folder + "/data/" + first_store + "/" + new_response.ResponderID)
	if err == nil {
		t.Fatal("Responder folder should be deleted")
	}

	responder.CreateResponderFolder(initialization_folder, new_response)
	tools.WriteFilesFromMultipart(initialization_folder, new_response)
	tools.WriteResponsesToJSONFile(initialization_folder, new_response)
	new_response_db_fields, err := responder.FormResponseToDBFormat(new_response)
	if err != nil {
		t.Fatal(err)
	}
	// A combination of Responses and File Locations listing a URL for file download where it will be served
	new_resp_marshal, err := json.Marshal(new_response_db_fields)
	if err != nil {
		t.Fatal(err)
	}
	if string(old_resp_marshal) == string(new_resp_marshal) {
		t.Fatal("The was no update")
	}
	err = tools.StoreResponseToDB(db, new_response_db_fields)
	if err != nil {
		t.Fatal(err)
	}
}

func UpdateToAnonOption(t *testing.T) {
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

	prebuilder.ReplyToFormScrambled(1, first_store, "192.168.1.1", db, initialization_folder)
	var files map[string]former.MultipartFile = prebuilder.CopyTestFilesToMemory(initialization_folder, map[string]string{"Test-FI": "test-file-1.jpg"})
	new_response := former.FormResponse{
		FormName:     first_store,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses: map[string]string{
			// Fill them out here
			"anon-option":   "false",
			"Test-TA":       "../some text\n\n\tasdf",
			"Test-GI":       "new text",
			"Test-Chk-SG-1": "ck1", // in this case the check group has been assigned ...-1 from the algorithm
			"Test-Chk-SG-3": "ck3", // in this case the check group has been assigned ...-3 from the algorithm
			"Test-rdo-SG":   "rd1", // In this case the radio group is just called
			"Test-optGrp":   "item-2",
		},
		// File paths are distinct in that their data effects file storage.
		// While text based fields may be passed around and later inserted somewhere,
		// Files must be moved around the OS
		// This means giving them a unique identifier, in this case a JSON column
		// It does not need a database column because after validation the files are written to a predictable location
		FileObjects: files,
	}

	form, err := tools.GetFormOfID(db, new_response.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}
	// Validate responses
	// Check files are valid
	var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(new_response.Responses, rebuild_group)
	// Won't be tested here as yet
	var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(new_response.FileObjects, rebuild_group)
	issue_array := append(text_issue_array, file_issue_array...)
	if len(issue_array) != 0 {
		t.Fatal("There should be no errors")
	}

	// Edit mode should scramble the IP to see if it exists elsewhere
	edit_mode, old_id, err := responder.CheckIfEdit(db, new_response)

	if err != nil {
		t.Fatal(err)
	}
	if edit_mode == false {
		t.Fatal("Edit mode was not set")
	}
	new_response.ScrambleResponderID()
	assumed_old_id := new_response.ResponderID
	new_response.ResponderID = "192.168.1.1"
	if old_id != assumed_old_id {
		t.Fatal("This should be scrambled")
	}

	if edit_mode {
		// Deleting is important because the responder ID could be set to scramble
		// Also easier and since nothing relies on the data it can be done
		destroyer.UndoResponse(db, new_response, old_id, initialization_folder)
	}

	_, err = os.Stat(initialization_folder + "/data/" + first_store + "/" + assumed_old_id)
	if err == nil {
		t.Fatal("Responder folder should be deleted")
	}
}

func TestExcessiveInputLen(t *testing.T) {
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

	terror_str := ""
	for i := 0; i < globals.MaxInputTextLen; i++ {
		terror_str += "E"
	}
	terror_str += "E"

	new_response := former.FormResponse{
		FormName:     first_store,
		RelationalID: 1,
		ResponderID:  "192.168.1.1",
		Responses: map[string]string{
			// Fill them out here
			"anon-option":   "false",
			"Test-TA":       terror_str,
			"Test-GI":       "new text",
			"Test-Chk-SG-1": "ck1", // in this case the check group has been assigned ...-1 from the algorithm
			"Test-Chk-SG-3": "ck3", // in this case the check group has been assigned ...-3 from the algorithm
			"Test-rdo-SG":   "rd1", // In this case the radio group is just called
			"Test-optGrp":   "item-2",
		},
	}

	form, err := tools.GetFormOfID(db, new_response.RelationalID)
	if err != nil {
		t.Fatal(err)
	}
	var rebuild_group former.FormConstruct
	err = json.Unmarshal([]byte(form.FieldJSON), &rebuild_group)
	if err != nil {
		t.Fatal(err)
	}
	// Validate responses
	// Check files are valid
	var issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(new_response.Responses, rebuild_group)
	// It should pass so write to the propper locations
	if len(issue_array) != 1 {
		t.Fatal("There should be errors here", issue_array, len(new_response.Responses["Test-TA"]), "vs", globals.MaxInputTextLen)
	}
	if issue_array[0].FailType != former.InvalidTextLengthMessage {
		t.Error("Incorrect error msg", issue_array[0].FailType)

	}
	if issue_array[0].FailCode != former.InvalidTextLengthCode {
		t.Error("Incorrect error code", issue_array[0].FailCode)
	}
}
