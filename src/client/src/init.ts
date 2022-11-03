import {recursiveFormFieldRebuild , createNewSubgroup} from "./builder"
import {submitForm} from "./submit"

export function helloWorld() {
  console.log("hello client");
}
export function attatchCreators(settings:any){
  let root_group_button = document.getElementById("sub-create")
  root_group_button.onclick = () => createNewSubgroup(root_group_button as HTMLButtonElement)

  let submit_button = document.getElementById("form-submit-button")
  if(submit_button) {
    submit_button.onclick = () =>  submitForm(submit_button as HTMLButtonElement , "/mod/create" )
  }

  let edit_button = document.getElementById("form-edit-button")
  if (edit_button){
    edit_button.onclick = () =>  submitForm(edit_button as HTMLButtonElement , "/mod/edit/" + settings.form_number)
  }
}

export function rebuildFromRaw(raw_json:string){
  try {
    raw_json = raw_json.replace(/\n/g, "\\n")
    raw_json = raw_json.replace(/\t/g, "\\t")
    raw_json = raw_json.replace(/[\r\n\t\f\v]/g, "")
    let form_construct = JSON.parse(raw_json)
    let field_structure = { SubGroups: form_construct.FormFields }
    let reference_button = document.getElementById("sub-create");
    recursiveFormFieldRebuild(field_structure , <HTMLButtonElement>reference_button)
  } catch (error) {
    console.error(error)
    alert("Issue with rebuilding")
  }
}
