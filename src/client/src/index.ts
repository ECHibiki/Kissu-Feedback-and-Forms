console.log("FormLibrary initialized.\nFeedback&Forms product of Kissu.moe");

import * as Submit from "./submit";
import * as Init from "./init";
import * as Builder from "./builder";

export function dropdownList(activator:HTMLAnchorElement, id:string) : boolean{
  return Submit.dropdownList(activator , id)
}

export function attatchCreators(settings){
  Init.attatchCreators(settings)
}

export function createNewResponseElement(button: HTMLButtonElement , insert_before_element?:HTMLElement ): ({base_id:string, respondable_container_id:string, container:HTMLDivElement }){
  return Builder.createNewResponseElement(button , insert_before_element)
}

export function createNewSubgroup(button: HTMLButtonElement , insert_before_element?:HTMLElement ): string{
  return Builder.createNewSubgroup(button , insert_before_element)
}

export function deleteContainer(base_container_id:string , sub_container_id:string){
  Builder.deleteContainer(base_container_id  , sub_container_id)
}

export function responseTypeSelected(respondable_container_id:string , button: HTMLButtonElement ){
  Builder.responseTypeSelected(respondable_container_id , button)
}

export function addListField(button : HTMLButtonElement , type:string) {
  Builder.addListField(button , type)
}

export function removeListItem(button : HTMLButtonElement , type:string){
  Builder.removeListItem(button , type)
}

export function rebuildFromRaw(raw_json:string){
  Init.rebuildFromRaw(raw_json)
}

export function setFieldsAsDisabled(disable_state:boolean){
  Builder.setFieldsAsDisabled(disable_state)
}

export function postDeleteForm(name: string , number:string) : boolean{
  return Submit.postDeleteForm(name , number)
}

export function postDeleteResponse(form_name: string , response_number:string) : boolean{
  return Submit.postDeleteResponse(form_name , response_number)
}

export function submitUserPost(form) : boolean{
  return Submit.submitUserPost(form)
}
