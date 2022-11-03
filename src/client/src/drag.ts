import {
  createFormStack ,
  rebuildRespondable ,
  convertRespondable ,
  setFieldsAsDisabled ,
  recursiveFormFieldRebuild ,
  convertContainerToJSON ,
} from "./builder"

var actively_dragged:HTMLElement;

function createZIndexModifier(id_target:string){
  let s = document.createElement("STYLE");
  s.innerHTML = `#${id_target} > div * {
    z-index: -1;
    position: relative;
  }`;
  s.id="zindexstyle";
  document.body.appendChild(s)
}
function removeZIndexModifier(){
  let s = document.getElementById("zindexstyle");
  document.body.removeChild(s)
}

export function handleContainerDragStartWithinParent(e){
  e.stopPropagation();
  if (this != e.target){
    e.preventDefault();
    return
  }
  actively_dragged = this.parentNode;
  this.parentNode.style.opacity = '0.4';
  this.style.cursor = "grabbing";
  createZIndexModifier((<HTMLElement>actively_dragged.parentNode).id);
}
export function handleContainerDragEndWithinParent(e){
  e.stopPropagation();
  e.preventDefault();
  actively_dragged = undefined;
  this.parentNode.style.opacity = "1.0";
  this.style.cursor = "";
  removeZIndexModifier();
}

export function handleContainerDragEnterWithinParent(e){
  if(e.currentTarget.parentNode != actively_dragged.parentNode || e.currentTarget == actively_dragged){
    e.preventDefault();
    return
  }
  e.currentTarget.style.border = '1px dashed red';
}
export function handleContainerDragLeaveWithinParent(e){
  // We want the lowest level listener for leaves
  if(e.currentTarget.parentNode != actively_dragged.parentNode || e.currentTarget == actively_dragged){
    e.preventDefault();
    return
  }
  e.currentTarget.style.border = '';
}
export function handleContainerDropWithinParent(e){
  if(e.currentTarget.parentNode != actively_dragged.parentNode  || e.currentTarget == actively_dragged){
    e.preventDefault();
    return
  }

  let cjson:any;
  let feedback_list = this.parentNode
  let is_subgroup = this.className.indexOf("sub-group") != -1
  if(is_subgroup){

    let current_group = [this , ...Array.from(this.childNodes)]

    let response_object = {
      FormFields: [ ]
    }
    let form_stack = createFormStack(current_group , response_object)
    cjson = convertContainerToJSON(response_object , form_stack)

    let field_structure = { SubGroups: response_object.FormFields }

    recursiveFormFieldRebuild(field_structure , <HTMLButtonElement>feedback_list.getElementsByTagName("BUTTON").item(0) , actively_dragged )
    setFieldsAsDisabled(false)
  } else{
    cjson = convertRespondable( this )
    rebuildRespondable(cjson , <HTMLButtonElement>feedback_list.parentNode.getElementsByTagName("BUTTON").item(0) , actively_dragged )
  }
  feedback_list.insertBefore(actively_dragged, this);
  feedback_list.removeChild(this);
}
