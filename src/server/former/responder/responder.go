package responder

import (
	// "fmt"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"database/sql"
	"io/ioutil"
	"strconv"
	"strings"
	"regexp"
	"errors"
	"time"
	"os"
)

// c.Request.FromFile
// http.DetectContentType

// Using header information can be problematic, just check the file itself
// Filename not important and would rather not deal with any issues relating to it
// Add to the templates that filenames are not preserved

// The input names of checkgroups are ignored for what the form says they should be

/*
  ResponseMissingMessage                       = "A field is required yet has no response."
  InvalidInputMessage                     ...
  InvalidOptionValueMessage                    = "The value of an options group does not exist on the server."
  InvalidFileExtMessage                        = "The extention of a file is not permitted on the server."
  InvalidFileSizeMessage                       = "The size of a file is too large."
*/

func ValidateTextResponsesAgainstForm(text_responses map[string]string, form former.FormConstruct) (error_list []former.FailureObject) {
	r_fo := validateRequiredTextFields(text_responses, form)
	i_fo := validateResponseTextFields(text_responses, form)
	error_list = append(error_list, r_fo...)
	error_list = append(error_list, i_fo...)
	return
}

func validateRequiredTextFields(text_responses map[string]string, form former.FormConstruct) (error_list []former.FailureObject) {

	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)

	// Validate required fields and select/option group first pass verification
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, respondable := range item.Respondables {
				fail := false
				required := respondable.Object.GetRequired()
				if !required {
					continue
				}

				name := respondable.Object.GetName()
				select_group, is_selection := respondable.Object.(former.SelectionGroup)
				if is_selection && select_group.SelectionCategory == former.Checkbox {
					answer_found := false
					for i := range select_group.CheckableItems {
						_, ok := text_responses[name+"-"+strconv.Itoa(i+1)]
						if ok {
							answer_found = true
						}
					}
					if !answer_found {
						error_list = append(error_list, former.FailureObject{former.ResponseMissingMessage, former.ResponseMissingCode, name})
					}
				} else {
					response, ok := text_responses[name]
					if !ok {
						fail = true
					} else if required && ok {
						trimmed := strings.TrimSpace(response)
						if len(trimmed) == 0 {
							fail = true
						}
					}
					if fail {
						error_list = append(error_list, former.FailureObject{former.ResponseMissingMessage, former.ResponseMissingCode, name})
					}
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}

	// Next validate that the responses actually make sense against the form
	return
}

func validateResponseTextFields(text_responses map[string]string, form former.FormConstruct) (error_list []former.FailureObject) {
	var subgroup_stack []former.FormGroup
	field_list := make(map[string]former.UnmarshalerFormObject)

	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// Validate required fields and select/option group first pass verification
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, respondable := range item.Respondables {
				name := respondable.Object.GetName()

				select_group, is_selection := respondable.Object.(former.SelectionGroup)
				if is_selection && select_group.SelectionCategory == former.Checkbox {
					for i := range select_group.CheckableItems {
						field_list[name+"-"+strconv.Itoa(i+1)] = respondable
					}
				} else {
					field_list[name] = respondable
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}

	for field, response := range text_responses {
		respondable, exists := field_list[field]
		if !exists && field != "anon-option" {
			error_list = append(error_list, former.FailureObject{former.InvalidInputMessage, former.InvalidInputCode, field})
			continue
		}
		if response == "" {
			continue
		}
		if len(response) > globals.MaxInputTextLen {
			error_list = append(error_list, former.FailureObject{former.InvalidTextLengthMessage, former.InvalidTextLengthCode, field})
		}

		options_respondable, is_options := respondable.Object.(former.OptionGroup)
		if is_options {
			found_value := false
			for _, opt := range options_respondable.Options {
				if response == opt.Value {
					found_value = true
					break
				}
			}
			if !found_value {
				error_list = append(error_list, former.FailureObject{former.InvalidOptionValueMessage, former.InvalidOptionValueCode, field})
			}
		}

		selection_group, is_select := respondable.Object.(former.SelectionGroup)
		if is_select {
			found_value := false
			for i, chk := range selection_group.CheckableItems {
				if response == chk.Value {
					found_value = true
					if selection_group.SelectionCategory == former.Checkbox {
						field_index, err := strconv.Atoi(field[strings.LastIndex(field, "-")+1:])
						if err != nil {
							error_list = append(error_list, former.FailureObject{former.InvalidSelectionIndexMessage, former.InvalidSelectionIndexCode, field})
							continue
						}
						if field_index != i+1 {
							error_list = append(error_list, former.FailureObject{former.InvalidSelectionIndexMessage, former.InvalidSelectionIndexCode, field})
						}
					}
					break
				}
			}
			if !found_value {
				error_list = append(error_list, former.FailureObject{former.InvalidSelectionValueMessage, former.InvalidSelectionValueCode, field})
			}
		}
	}

	return
}

func ValidateFileObjectsAgainstForm(file_tags map[string]former.MultipartFile, form former.FormConstruct) (error_list []former.FailureObject) {
	r_fo := validateRequiredFileFields(file_tags, form)
	i_fo := validateResponseFileFields(file_tags, form)
	error_list = append(error_list, r_fo...)
	error_list = append(error_list, i_fo...)
	return
}

func validateRequiredFileFields(file_tags map[string]former.MultipartFile, form former.FormConstruct) (error_list []former.FailureObject) {

	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)

	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, respondable := range item.Respondables {
				input, is_file := respondable.Object.(former.FileInput)
				if !is_file {
					continue
				}
				fail := false
				name := respondable.Object.GetName()
				required := respondable.Object.GetRequired()
				file, ok := file_tags[name]
				if required && !ok {
					fail = true
				} else if ok {
					// begin file verification
					if file.Header.Size > input.MaxSize {
						error_list = append(error_list, former.FailureObject{former.InvalidFileSizeMessage, former.InvalidFileSizeCode, name})
					}
					r, e := regexp.Compile(input.AllowedExtRegex)
					if e != nil {
						error_list = append(error_list, former.FailureObject{former.InvalidExtRegexMessage, former.InvalidExtRegexCode, name})
						continue
					}
					if !r.Match([]byte(file.Header.Filename)) {
						error_list = append(error_list, former.FailureObject{former.InvalidFileExtMessage, former.InvalidFileExtCode, name})
					}
					if strings.Contains(file.Header.Filename, "/") {
						error_list = append(error_list, former.FailureObject{former.DangerousPathMessage, former.DangerousPathCode, name})
					}
				}
				if fail {
					error_list = append(error_list, former.FailureObject{former.ResponseMissingMessage, former.ResponseMissingCode, name})
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}
	return
}

func validateResponseFileFields(file_tags map[string]former.MultipartFile, form former.FormConstruct) (error_list []former.FailureObject) {
	var subgroup_stack []former.FormGroup
	field_list := make(map[string]former.UnmarshalerFormObject)

	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// Validate required fields and select/option group first pass verification
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, respondable := range item.Respondables {
				name := respondable.Object.GetName()
				field_list[name] = respondable
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}

	for field := range file_tags {
		_, exists := field_list[field]
		if !exists {
			error_list = append(error_list, former.FailureObject{former.InvalidInputMessage, former.InvalidInputCode, field})
			continue
		}
	}

	return
}

//
func CreateResponderFolder(root_dir string, response_struct former.FormResponse) error {
	var err error

	if strings.Contains(response_struct.FormName, "/") || strings.Contains(response_struct.ResponderID, "/") {
		return errors.New("ResponderID contains illegal character '/'")
	}

	err = os.Mkdir(root_dir+"/data/"+response_struct.FormName+"/"+response_struct.ResponderID+"/", 0755)
	if err != nil {
		return err
	}
	if len(response_struct.FileObjects) > 0 {
		err = os.Mkdir(root_dir+"/data/"+response_struct.FormName+"/"+response_struct.ResponderID+"/files/", 0755)
	}
	return err
}

func FormResponseToDBFormat(response former.FormResponse) (types.ResponseDBFields, error) {
	for k, v := range response.FileObjects {
		response.Responses[k] = v.Header.Filename
	}

	resp_bytes, err := json.Marshal(response.Responses)
	return types.ResponseDBFields{
		ID:           0,
		FK_ID:        response.RelationalID,
		Identifier:   response.ResponderID,
		ResponseJSON: string(resp_bytes),
		SubmittedAt:  time.Now().Unix(),
	}, err
}

func CheckIfEdit(db *sql.DB, new_response former.FormResponse) (bool, string, error) {
	scrambled_id := new_response.GetScrambledID()
	r := db.QueryRow("SELECT identifier FROM responses WHERE fk_id = ? AND (identifier = ? OR identifier = ?)", new_response.RelationalID, new_response.ResponderID, scrambled_id)
	var id string
	err := r.Scan(&id)
	if err != nil {
		return false, "", err
	} else {
		return true, id, nil
	}
}

func FillMapWithPostParams(c *gin.Context, resp_map map[string]string, form former.FormConstruct) {
	if len(form.FormFields) == 0 {
		return
	}
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, r := range item.Respondables {
				if r.Type == former.FileInputTag {
					continue
				}
				if r.Type == former.SelectionGroupTag && r.Object.(former.SelectionGroup).SelectionCategory == former.Checkbox {
					for i := range r.Object.(former.SelectionGroup).CheckableItems {
						if _, exists := c.GetPostForm(r.Object.GetName() + "-" + strconv.Itoa(i+1)); exists {
							resp_map[r.Object.GetName()+"-"+strconv.Itoa(i+1)] = c.PostForm(r.Object.GetName() + "-" + strconv.Itoa(i+1))
						}
					}
				} else {
					resp_map[r.Object.GetName()] = c.PostForm(r.Object.GetName())
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}
}
func FillMapWithPostFiles(c *gin.Context, file_map map[string]former.MultipartFile, form former.FormConstruct) {
	if len(form.FormFields) == 0 {
		return
	}
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if len(item.Respondables) != 0 {
			for _, r := range item.Respondables {
				if r.Type == former.FileInputTag {
					image, header, _ := c.Request.FormFile(r.Object.GetName())
					mp := former.MultipartFile{
						File:   image,
						Header: header,
					}
					file_map[r.Object.GetName()] = mp
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}
}

func WriteResponsesToJSONFile(root_dir string, resp former.FormResponse) error {
	storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"

	json_resp := ConvertFormResponseToJSONFormResponse(root_dir, resp)

	json_bytes, err := json.MarshalIndent(json_resp, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(storage_dir+"responses.json", json_bytes, 0644)
	return err
}

func ConvertFormResponseToJSONFormResponse(root_dir string, resp former.FormResponse) former.JSONFormResponse {
	json_resp := former.JSONFormResponse{}
	json_resp.FormName = resp.FormName
	json_resp.RelationalID = resp.RelationalID
	json_resp.ResponderID = resp.ResponderID
	json_resp.Responses = resp.Responses
	json_resp.FilePaths = make(map[string]string)
	storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"
	for k, v := range resp.FileObjects {
		json_resp.FilePaths[k] = storage_dir + "files/" + v.Header.Filename
	}

	return json_resp
}

func StoreResponseToDB(db *sql.DB, db_response types.ResponseDBFields) error {
	_, err := db.Exec("INSERT INTO responses VALUES( NULL , ? , ? , ? , ? )", db_response.FK_ID, db_response.Identifier, db_response.ResponseJSON, db_response.SubmittedAt)
	return err
}
