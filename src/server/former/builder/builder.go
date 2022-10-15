package builder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"os"
	"regexp"
	"strconv"
	"time"
)

func ConstructFormObject(inputs map[string]string) former.FormConstruct {
	return former.FormConstruct{}
}

func MakeFormWritable(form former.FormConstruct) (types.FormDBFields, error) {
	marshaled_form, err := json.Marshal(form)
	if err != nil {
		return types.FormDBFields{}, err
	}
	return types.FormDBFields{
		ID:        0,
		Name:      form.StorageName(),
		FieldJSON: string(marshaled_form),
		UpdatedAt: time.Now().Unix(),
	}, nil
}

func CreateFormDirectory(form former.FormConstruct, root_dir string) error {
	safe_name := form.StorageName()
	err := os.Mkdir(root_dir+"/data/"+safe_name+"/", 0775)
	fmt.Printf("Writting new form to %s/data/%s/ \n", root_dir, safe_name)
	if err != nil {
		return err
	}
	return nil
}

func CheckFormDirectoryExists(form former.FormConstruct, root_dir string) error {
	safe_name := form.StorageName()
	_, err := os.Stat(root_dir + "/data/" + safe_name + "/")
	return err
}

func checkNameUnique(db *sql.DB, form former.FormConstruct) (error_list former.FailureObject) {
	row := db.QueryRow("SELECT * FROM forms WHERE name = ?", form.StorageName())
	var db_field types.FormDBFields
	err := row.Scan(&db_field.ID, &db_field.Name, &db_field.FieldJSON, &db_field.UpdatedAt)
	if err == nil {
		return former.FailureObject{former.DuplicateFormNameMessage, former.DuplicateFormNameCode, "FormName"}
	}
	return former.FailureObject{}
}

func ValidateForm(db *sql.DB, form former.FormConstruct) (error_list []former.FailureObject) {
	if form.FormName == "" {
		error_list = append(error_list, former.FailureObject{former.HeadMissingFormNameMessage, former.HeadMissingFormNameCode, form.FormName})
	}

	dupe_err := checkNameUnique(db, form)
	if dupe_err.FailType != "" {
		error_list = append(error_list, dupe_err)
	}

	uniqueness_errors := checkNameAndIDUniqueness(form)
	if len(uniqueness_errors) > 0 {
		error_list = append(error_list, uniqueness_errors...)
	}

	character_errors := checkNameAndIDPropperCharacters(form)
	if len(character_errors) > 0 {
		error_list = append(error_list, character_errors...)
	}

	// uniqueness errors through all other chekcs into confusion
	// allowing for other errors to display isn't important... not even golang shows all error types at once
	// this limitation, I guess, issue will apply to situations where a submits without the client
	if len(uniqueness_errors) > 0 {
		return error_list
	}

	struct_errors := checkValidFormStructure(form)
	if len(struct_errors) > 0 {
		error_list = append(error_list, struct_errors...)
	}
	return error_list
}

func checkValidFormStructure(form former.FormConstruct) (error_list []former.FailureObject) {
	if len(form.FormFields) == 0 {
		return []former.FailureObject{{former.HeadMissingMessage, former.HeadMissingCode, form.ID}}
	}
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		fail_location := item.ID
		if fail_location != "" && len(item.Respondables) == 0 && item.Description == "" {
			error_list = append(error_list, former.FailureObject{former.GroupMissingMessage, former.GroupMissingCode, fail_location})
		}
		// verify it has validity to it
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}
	return error_list
}

// For the next few, create a []struct {isID:bool, name:string }.
// From this struct perform the checks
func checkNameAndIDUniqueness(form former.FormConstruct) (error_list []former.FailureObject) {
	if len(form.FormFields) == 0 {
		return []former.FailureObject{}
	}
	var id_checklist = map[string]uint{form.ID: 1}
	var name_checklist = make(map[string]uint)
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		if _, ok := id_checklist[item.ID]; !ok {
			id_checklist[item.ID] = 1
		} else if id_checklist[item.ID] != 2 {
			id_checklist[item.ID] = 2
			error_list = append(error_list, former.FailureObject{former.DuplicateIDMessage, former.DuplicateIDCode, item.ID})
		}
		if len(item.Respondables) != 0 {
			for _, r := range item.Respondables {
				sg, is_sj := r.Object.(former.SelectionGroup)
				if is_sj && sg.SelectionCategory == former.Checkbox {
					for i := range sg.CheckableItems {
						name := r.Object.GetName() + "-" + strconv.Itoa(i+1)
						if _, ok := name_checklist[name]; !ok {
							name_checklist[name] = 1
						} else if name_checklist[name] != 2 {
							name_checklist[name] = 2
							error_list = append(error_list, former.FailureObject{former.InvalidCheckboxMessage, former.InvalidCheckboxCode, r.Object.GetName()})
						}
					}
				} else {
					name := r.Object.GetName()
					if _, ok := name_checklist[name]; !ok {
						name_checklist[name] = 1
					} else if name_checklist[name] != 2 {
						name_checklist[name] = 2
						error_list = append(error_list, former.FailureObject{former.DuplicateNameMessage, former.DuplicateNameCode, name})
					}
				}
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}

	return error_list
}

func checkNameAndIDPropperCharacters(form former.FormConstruct) (error_list []former.FailureObject) {

	var ids []string = []string{form.ID}
	var names []string
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[len(subgroup_stack)-1]
		subgroup_stack = subgroup_stack[:len(subgroup_stack)-1]
		ids = append(ids, item.ID)
		if len(item.Respondables) != 0 {
			for _, r := range item.Respondables {
				name := r.Object.GetName()
				names = append(names, name)
			}
		}
		if len(item.SubGroups) != 0 {
			// add children to the stack
			subgroup_stack = append(subgroup_stack, item.SubGroups...)
		}
	}

	invalid_entry := regexp.MustCompile("^[^a-zA-Z]")
	invalid_body := regexp.MustCompile("[^a-zA-Z0-9\\-_:\\.]")

	for _, id := range ids {
		if len(id) == 0 {
			error_list = append(error_list, former.FailureObject{former.EmptyIDMessage, former.EmptyIDCode, id})
		}
		if invalid_entry.Match([]byte(id)) {
			error_list = append(error_list, former.FailureObject{former.InvalidIDStarterMessage, former.InvalidIDStarterCode, id})
		}
		if invalid_body.Match([]byte(id)) {
			error_list = append(error_list, former.FailureObject{former.InvalidIDCharactersMessage, former.InvalidIDCharactersCode, id})
		}
	}
	for _, name := range names {
		if len(name) == 0 {
			error_list = append(error_list, former.FailureObject{former.EmptyNameMessage, former.EmptyNameCode, name})
		}
		if invalid_entry.Match([]byte(name)) {
			error_list = append(error_list, former.FailureObject{former.InvalidNameStarterMessage, former.InvalidNameStarterCode, name})
		}
		if invalid_body.Match([]byte(name)) {
			error_list = append(error_list, former.FailureObject{former.InvalidNameCharactersMessage, former.InvalidNameCharactersCode, name})
		}
	}
	return error_list
}

func ValidateFormEdit(replace former.FormConstruct, base former.FormConstruct) (error_list []former.FailureObject) {
	if replace.FormName != base.FormName {
		error_list = append(error_list, former.FailureObject{former.EditNameChangeMessage, former.EditNameChangeCode, "FormName"})
	}

	// base is valid of course

	uniqueness_errors := checkNameAndIDUniqueness(replace)
	if len(uniqueness_errors) > 0 {
		error_list = append(error_list, uniqueness_errors...)
	}

	character_errors := checkNameAndIDPropperCharacters(replace)
	if len(character_errors) > 0 {
		error_list = append(error_list, character_errors...)
	}

	// uniqueness errors through all other chekcs into confusion
	// allowing for other errors to display isn't important... not even golang shows all error types at once
	// this limitation, I guess, issue will apply to situations where a submits without the client
	if len(uniqueness_errors) > 0 {
		return error_list
	}

	struct_errors := checkValidFormStructure(replace)
	if len(struct_errors) > 0 {
		error_list = append(error_list, struct_errors...)
	}
	return
}

func StoreForm(db *sql.DB, db_form types.FormDBFields) (int, error) {
	err := StoreFormToDB(db, db_form)
	if err != nil {
		return 0, err
	}
	return tools.GetLastIndex(db, "forms")
}
func UpdateForm(db *sql.DB, index int64, db_form types.FormDBFields) error {
	_, err := db.Exec("UPDATE forms SET field_json = ? , updated_at = ? WHERE id = ?", db_form.FieldJSON, db_form.UpdatedAt, index)
	return err
}
func StoreFormToDB(db *sql.DB, db_form types.FormDBFields) error {
	// Unique name prevents duplicate entries while auto-incremented ID makes for easy foreign keys
	_, err := db.Exec("INSERT INTO forms VALUES( NULL , ? , ? , ? )", db_form.Name, db_form.FieldJSON, db_form.UpdatedAt)
	return err
}
