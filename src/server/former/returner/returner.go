package returner

import (
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"github.com/tyler-sommer/stick"
	"path/filepath"
	"encoding/json"
	"compress/gzip"
	"encoding/csv"
	"database/sql"
	"archive/tar"
	"io/ioutil"
	"strconv"
	"errors"
	"bytes"
	"fmt"
	"os"
	"io"
)

func RenderTestingTemplate[T int64 | string](db *sql.DB, env *stick.Env, root_dir string, db_key T) (string, error) {

	var returned_form types.FormDBFields
	var rebuild_group former.FormConstruct
	var err error

	var i interface{} = db_key
	switch i.(type) {
	case int64:
		returned_form, err = GetFormOfID(db, i.(int64))
	case string:
		returned_form, err = GetFormOfName(db, i.(string))
	}
	if err != nil {
		fmt.Printf("%v\n", db_key)
		panic(err)
	}

	err = json.Unmarshal([]byte(returned_form.FieldJSON), &rebuild_group)
	if err != nil {
		return "", err
	}
	// Turn rebuild_group into a templatable format
	var construction_variables map[string]stick.Value = map[string]stick.Value{"form": rebuild_group}

	// Render a form only used for testing
	testing_form_render, err := templater.ReturnFilledTemplate(env, root_dir+"/templates/test-views/render-test.twig", construction_variables)
	return testing_form_render, err
}

func GetAllForms(db *sql.DB) (parsed_row_list []types.FormDBFields, err error) {
	row_list, err := db.Query("SELECT id, name, updated_at FROM forms ORDER BY updated_at DESC")
	if err != nil {
		return
	}
	defer row_list.Close()
	for row_list.Next() {
		var parsed_row types.FormDBFields
		err = row_list.Scan(&parsed_row.ID, &parsed_row.Name, &parsed_row.UpdatedAt)
		if err != nil {
			return
		}
		parsed_row_list = append(parsed_row_list, parsed_row)
	}
	return
}

func GetRepliesToForm(db *sql.DB, id int64) (parsed_row_list []types.ResponseDBFields, err error) {
	row_list, err := db.Query("SELECT id, fk_id, identifier, response_json, submitted_at FROM responses WHERE fk_id = ? ORDER BY id DESC", id)
	if err != nil {
		return
	}
	defer row_list.Close()
	for row_list.Next() {
		var parsed_row types.ResponseDBFields
		err = row_list.Scan(&parsed_row.ID, &parsed_row.FK_ID, &parsed_row.Identifier, &parsed_row.ResponseJSON, &parsed_row.SubmittedAt)
		if err != nil {
			return
		}
		parsed_row_list = append(parsed_row_list, parsed_row)
	}
	return
}

func CreateInstancedCSVForGivenForm(db *sql.DB, id int64, initialization_folder string) error {
	form_data, err := GetFormOfID(db, id)
	if err != nil {
		return err
	}
	var form_construct former.FormConstruct
	err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
	if err != nil {
		return err
	}
	var csv_list [][]string
	var field_list []string
	var field_map map[string]int = make(map[string]int)

	field_list = append(field_list, "Identifier")
	field_map["Identifier"] = 0

	fields := GetFieldsOfFormConstruct(form_construct)
	for field_index, field := range fields {
		if field.Type == former.SelectionGroupTag {
			sg := field.Object.(former.SelectionGroup)
			if sg.SelectionCategory == former.Checkbox {
				for chk_index := 0; chk_index < len(sg.CheckableItems); chk_index++ {
					chk_str_index := strconv.Itoa(chk_index + 1)
					// + 1 because it's possitioned based on the identifier being set
					field_map[field.Object.GetName()+"-"+chk_str_index] = field_index + chk_index + 1
					field_list = append(field_list, field.Object.GetName()+"-"+chk_str_index)
				}
			}
		} else {
			// + 1 because it's possitioned based on the identifier being set
			field_map[field.Object.GetName()] = field_index + 1
			field_list = append(field_list, field.Object.GetName())
		}

	}

	field_list = append(field_list, "SubmittedAt")
	field_map["SubmittedAt"] = len(field_list) - 1

	csv_list = append(csv_list, field_list)

	responses, err := GetRepliesToForm(db, id)
	for _, r := range responses {
		responses_list := make([]string, len(field_list))
		responses_list[field_map["Identifier"]] = r.Identifier
		responses_list[field_map["SubmittedAt"]] = strconv.Itoa(int(r.SubmittedAt))
		var response map[string]string = make(map[string]string)
		err = json.Unmarshal([]byte(r.ResponseJSON), &response)
		if err != nil {
			return err
		}
		for k, v := range response {
			if _, exists := field_map[k]; !exists {
				fmt.Printf("%s %s Does not exist on field list", k, v)
				continue
			}

			responses_list[field_map[k]] = v
		}
		csv_list = append(csv_list, responses_list)
	}
	err = writeCSVToDir(initialization_folder+"/data/"+form_data.Name+"/data.csv", csv_list)
	return err
}

func CreateReadmeForGivenForm(db *sql.DB, id int64, initialization_folder string) error {
	form_data, err := GetFormOfID(db, id)
	if err != nil {
		return err
	}
	var form_construct former.FormConstruct
	err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
	if err != nil {
		return err
	}

	fields := GetFieldsOfFormConstruct(form_construct)

	var field_map map[string]string = make(map[string]string)
	field_map["FormName"] = form_data.Name
	field_map["ID"] = strconv.Itoa(int(form_data.ID))
	field_map["FormDescription"] = form_construct.Description
	field_map["AnonOption"] = strconv.FormatBool(form_construct.AnonOption)
	for _, field := range fields {
		field_map[field.Object.GetName()] = field.Object.GetDescription()
	}
	err = writeJSONReadmeToDir(initialization_folder+"/data/"+form_data.Name+"/field-descriptors.json", field_map)
	return err
}

// borrowing from https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
func CreateDownloadableForGivenForm( form_name string , initialization_folder string) error {
	form_dir := initialization_folder + "/data/" + form_name + "/"
	file_path := initialization_folder + "/data/" + form_name + "/downloadable.tar.gz"

	os.Remove(file_path)

	var buf bytes.Buffer
	zr := gzip.NewWriter(&buf)
	tw := tar.NewWriter(zr)
	// walk through every file in the folder
	err := filepath.Walk(form_dir, func(file string, fi os.FileInfo, err error) error {
		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}

	compressed_file, err := os.OpenFile(file_path, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return err
	}
	if _, err := io.Copy(compressed_file, &buf); err != nil {
		return err
	}
	return nil
}

func writeJSONReadmeToDir(filename string, field_map map[string]string) error {
	nice_marshal, err := json.MarshalIndent(field_map, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, nice_marshal, 0644)
}

func writeCSVToDir(filename string, csv_data [][]string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f)
	return writer.WriteAll(csv_data)
}

func GetFieldsOfFormConstruct(form former.FormConstruct) (field_list []former.UnmarshalerFormObject) {
	if len(form.FormFields) == 0 {
		return
	}
	var subgroup_stack []former.FormGroup
	subgroup_stack = append(subgroup_stack, form.FormFields...)
	// fail location identified by an ID
	for len(subgroup_stack) > 0 {
		item := subgroup_stack[0]
		subgroup_stack = subgroup_stack[1:]
		if len(item.Respondables) != 0 {
			for _, r := range item.Respondables {
				name := r.Object.GetName()
				name_found := false
				for _, v := range field_list {
					if v.Object.GetName() == name {
						name_found = true
						break
					}
				}
				if !name_found {
					field_list = append(field_list, r)
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

func GetResponseByID(db *sql.DB, id int64) (types.ResponseDBFields, error) {
	q := db.QueryRow("SELECT * FROM responses WHERE id=?", id)
	var db_response types.ResponseDBFields
	err := q.Scan(&db_response.ID, &db_response.FK_ID, &db_response.Identifier, &db_response.ResponseJSON, &db_response.SubmittedAt)
	if err != nil {
		return db_response, err
	}
	if db_response.FK_ID == 0 {
		return db_response, errors.New("Database has no row for ID")
	}
	return db_response, nil
}

func GetFormByNameAndID(db *sql.DB, name string, id int64) (types.FormDBFields, error) {
	data := db.QueryRow("SELECT * FROM forms WHERE name = ? AND id = ?", name, id)

	var rtn types.FormDBFields
	err := data.Scan(&rtn.ID, &rtn.Name, &rtn.FieldJSON, &rtn.UpdatedAt)
	if err != nil {
		return types.FormDBFields{}, err
	}
	return rtn, nil

}

func GetFormOfID(db *sql.DB, id int64) (types.FormDBFields, error) {
	q := db.QueryRow("SELECT * FROM forms WHERE id=?", id)
	var db_form types.FormDBFields
	err := q.Scan(&db_form.ID, &db_form.Name, &db_form.FieldJSON, &db_form.UpdatedAt)
	return db_form, err
}
func GetFormOfName(db *sql.DB, name string) (types.FormDBFields, error) {
	q := db.QueryRow("SELECT * FROM forms WHERE name=?", name)
	var db_form types.FormDBFields
	err := q.Scan(&db_form.ID, &db_form.Name, &db_form.FieldJSON, &db_form.UpdatedAt)
	return db_form, err
}
