import {
  handleContainerDropWithinParent ,
  handleContainerDragEndWithinParent ,
  handleContainerDragStartWithinParent ,
  handleContainerDragLeaveWithinParent ,
  handleContainerDragEnterWithinParent ,
} from "./drag"

export function deleteContainer(base_container_id:string , sub_container_id:string){
  document.getElementById(base_container_id).removeChild(document.getElementById( sub_container_id));
}
export function responseTypeSelected(respondable_container_id:string , button: HTMLButtonElement ){
  let base_id = button.getAttribute("data-link-id")
  let select_obj = document.getElementById(base_id + "-type")
  let container = document.getElementById(base_id + "-fields") as HTMLDivElement
  let selection = (select_obj as HTMLInputElement).value;
  switch (selection) {
    case "textarea":
      createTextAreaInputs(base_id , respondable_container_id , container)
      break;
    case "input":
      createInputInputs(base_id , respondable_container_id , container)
      break;
    case "file":
      createFileInputs(base_id , respondable_container_id , container)
      break;
    case "select":
      createSelectGroup(base_id , respondable_container_id , container)
      break;
    case "option":
      createOptionsGroup(base_id , respondable_container_id , container)
      break;

    default:
      break;
  }
}

export function createNewSubgroup( button: HTMLButtonElement  , insert_before_element?:HTMLElement): string{
  let parent_id = button.getAttribute("data-link-id")
  let parent_container =  document.getElementById(parent_id + "-group")

  let group_id = "group" + ((Date.now() + Math.random())*10000)
  let container = document.createElement('DIV');
  container.setAttribute("style" , "")
  container.className =  "sub-group form-group feedback-group"
  container.id =  group_id + "-group"
  container.innerHTML = `<LABEL>Group Label : <INPUT ondragstart="return false" draggable="false" type="text" name="form-label" id="${group_id}-label"/> </LABEL> <br/>
  <LABEL>Group ID : <INPUT ondragstart="return false" draggable="false" type="text" name="id" id="${group_id}-id"/></LABEL> <br/>
  <LABEL>Group Descriptor : <br/> <TEXTAREA name="description" id="${group_id}-description"></TEXTAREA></LABEL> <br/>
  <BUTTON id="${group_id}" onclick="FormLibrary.createNewResponseElement(this)" data-link-id="${group_id}" >Create New Respondable Below Last Respondable</BUTTON><br/>
  <BUTTON  onclick="FormLibrary.createNewSubgroup(this)" data-link-id="${group_id}">Create New Group Below Last Respondable</BUTTON><br/>
  <BUTTON  onclick="FormLibrary.deleteContainer('${parent_id + "-group"}' , '${ group_id + "-group"}')">Delete Subgroup</BUTTON><br/>
  <SPAN  class="respondable-container" id="${group_id}-respondables"></SPAN><BR/>
  `;
  container.ondrop = handleContainerDropWithinParent
  container.ondragenter = handleContainerDragEnterWithinParent
  container.ondragleave = handleContainerDragLeaveWithinParent
  container.ondragover = function (e) {
    e.preventDefault()
  }

  appendDropIcon(container)

  if( insert_before_element == undefined){
    //insert before the existing drop-icon
    insert_before_element = <HTMLDivElement>parent_container.lastChild
  }
  parent_container.insertBefore(container , insert_before_element )
  return group_id
}

export function addListField(button : HTMLButtonElement , type:string) {
  let parent_id = button.getAttribute("data-link-id")
  let ol = document.getElementById(parent_id + "-" + type)
  let li = document.createElement("LI")
  let child_count = ol.children.length
  li.innerHTML = `
    <INPUT ondragstart="return false" draggable="false" placeholder="Label" data-list-item-no="${child_count}" data-field="1" name='Label' id="${parent_id}-${type}-label-${child_count}"/>
    <INPUT ondragstart="return false" draggable="false" placeholder="Value" data-list-item-no="${child_count}" data-field="1" name='Value' id="${parent_id}-${type}-value-${child_count}"/>
  `
  ol.appendChild(li)
}
export function removeListItem(button : HTMLButtonElement , type:string) {
  let parent_id = button.getAttribute("data-link-id")
  let ol = document.getElementById(parent_id + "-" + type)
  if (ol.childNodes.length <= 1) {
    return
  }
  ol.removeChild(ol.lastChild)
}


export function setFieldsAsDisabled(disable_state:boolean){
  let buttons = Array.from(document.getElementsByTagName("BUTTON"))
  let inputs = Array.from(document.getElementsByTagName("INPUT"))
  let selects = Array.from(document.getElementsByTagName("SELECT"))
  let disablables = [].concat(buttons, inputs, selects)
  for( let i = 0 ; i < disablables.length ; i++){
    disable_state ? disablables[i].setAttribute("disabled" , "" + disable_state) : disablables[i].removeAttribute("disabled");
  }
}




export function createNewResponseElement(button: HTMLButtonElement , insert_before_element?:HTMLElement ): ({base_id:string, respondable_container_id:string, container:HTMLDivElement }){
  let parent_id = button.getAttribute("data-link-id")
  let respondable_container =  document.getElementById(parent_id + "-respondables")

  let res_id = "response" + ((Date.now() + Math.random())*10000)
  let container = document.createElement('DIV');
  container.setAttribute("style", "width:400px;min-height:200px")
  container.setAttribute("data-type", "blank")
  container.className =  "creation-prompt respondable-group feedback-group"
  container.id =  res_id + "-fields"
  container.innerHTML = `<SPAN>Element Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        <LABEL>Item Type:</LABEL>
        <SELECT id="${res_id}-type" ondragstart="return false" draggable="false">
          <OPTION value="textarea">TextArea</OPTION>
          <OPTION value="input">Input</OPTION>
          <OPTION value="file">FileInput</OPTION>
          <OPTION value="select">SelectGroup</OPTION>
          <OPTION value="option">OptionGroup</OPTION>
        </SELECT>
      </LI>
      <LI>If we want any of the unimplemented features, then you'll have to ask me or wait until I personally require it</LI>
      <LI><BUTTON id="${res_id}" data-link-id="${res_id}" onclick="FormLibrary.responseTypeSelected('${parent_id + "-respondables"}' , this)">Next</BUTTON></LI>
      <LI><BUTTON  data-link-id="${res_id}" onclick="FormLibrary.deleteContainer('${parent_id + "-respondables"}' , '${res_id + "-fields"}')">Delete</BUTTON></LI>
    </UL>`;
    container.ondrop = handleContainerDropWithinParent
    container.ondragenter = handleContainerDragEnterWithinParent
    container.ondragleave = handleContainerDragLeaveWithinParent
    container.ondragover = function (e) {
      e.preventDefault()
    }

    appendDropIcon(container)

    if( insert_before_element == undefined){
      respondable_container.appendChild(container);
    } else{
      respondable_container.insertBefore(container , insert_before_element )
    }
    // CREATE NEW AT END
    // INSERT IN POSITION
    return { base_id: res_id, respondable_container_id: parent_id + "-respondables" , container: <HTMLDivElement>container  }
}

export function buildDropdownResponse(container:HTMLDivElement , json:any) : string{
  let html = `<UL class="item-ul">${
    (() => {
      if (!json.formatted_replies){
        return "<LI class='item-li'>Empty Set</LI>"
      }

      let list = "";
      json.formatted_replies.forEach((r) => {
        list += `<LI class='item-li'>${r}</LI>`;
      });
      return list
    })()
  }</UL>`;
  return html
}


export function recursiveFormFieldRebuild(field_group:any , base_button:HTMLButtonElement , insert_before_element?:HTMLElement){
  if(!field_group){
    return
  }

  let respondables = field_group.Respondables
  if(respondables){
    respondables.forEach(field => {
      rebuildRespondable(field  , base_button)
    });
  }
  setFieldsAsDisabled(true)

  let group_structures = field_group.SubGroups
  if(group_structures){
    group_structures.forEach(structure => {
      let button_ref = createNewSubgroup(base_button , insert_before_element);
      (<HTMLInputElement>document.getElementById(button_ref + "-label")).value = structure.Label;
      (<HTMLInputElement>document.getElementById(button_ref + "-id")).value = structure.ID;
      (<HTMLTextAreaElement>document.getElementById(button_ref + "-description")).value = structure.Description;

      let reference_button = document.getElementById(button_ref);
      recursiveFormFieldRebuild(structure , <HTMLButtonElement>reference_button)
    });
  }
}

export function rebuildRespondable(field , base_button, insert_before_element?:HTMLElement){
  let respondable_creation_object = createNewResponseElement(base_button , insert_before_element)
  let field_id = ""
  switch (field.Type) {
    case "textarea":
      field_id = createTextAreaInputs(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id ,respondable_creation_object.container);
      (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
      (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
      (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;
      (<HTMLInputElement>document.getElementById(field_id + "-placeholder")).value = field.Object.Placeholder;
      break;
    case "genericinput":
    field_id = createInputInputs(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id ,respondable_creation_object.container);
      (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
      (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
      (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;
      (<HTMLInputElement>document.getElementById(field_id + "-placeholder")).value = field.Object.Placeholder;
      (<HTMLInputElement>document.getElementById(field_id + "-type")).value = field.Object.Type;
      break;
    case "fileinput":
    field_id = createFileInputs(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id ,respondable_creation_object.container);
      (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
      (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
      (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;
      (<HTMLInputElement>document.getElementById(field_id + "-allowed-ext")).value = field.Object.AllowedExtRegex;
      (<HTMLInputElement>document.getElementById(field_id + "-max-size")).value = field.Object.MaxSize;
      break;
    case "selectiongroup":
      field_id = createSelectGroup(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id , respondable_creation_object.container);
      (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
      (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
      (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;

      for(let check_count = 1 ; check_count < field.Object.CheckableItems.length ; check_count++ ){
         addListField( <HTMLButtonElement>document.getElementById(field_id) , "checkable" )
      }
      field.Object.CheckableItems.forEach((check , index) => {
        (<HTMLInputElement>document.getElementById(field_id + "-checkable-label-" + index)).value = check.Label;
        (<HTMLInputElement>document.getElementById(field_id + "-checkable-value-" + index)).value = check.Value;
      });
      break;
    case "optiongroup":
      field_id = createOptionsGroup(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id , respondable_creation_object.container);
      (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
      (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
      (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;

      for(let opt_count = 1 ; opt_count < field.Object.Options.length ; opt_count++ ){
         addListField( <HTMLButtonElement>document.getElementById(field_id) , "option" )
      }
      field.Object.Options.forEach((opt , index) => {
        (<HTMLInputElement>document.getElementById(field_id + "-option-label-" + index)).value = opt.Label;
        (<HTMLInputElement>document.getElementById(field_id + "-option-value-" + index)).value = opt.Value;
      });
      break;
    default:
      break;
  }
}

function appendDropIcon(container){
  let drop_icon = document.createElement('DIV');
  drop_icon.setAttribute("draggable" , "true");
  drop_icon.className = "drop-icon";
  drop_icon.textContent = " . . . ";
  drop_icon.title = "Drag this to relocate item";
  drop_icon.ondragstart = handleContainerDragStartWithinParent
  drop_icon.ondragend = handleContainerDragEndWithinParent

  drop_icon.onmouseover = function(e){
    // (<HTMLElement>(<HTMLDivElement>e.target).parentNode).style.color = "red"
    (<HTMLElement>(<HTMLDivElement>e.target).parentNode).style.borderColor = "red"
  }
  drop_icon.onmouseleave = function(e){
    // (<HTMLElement>(<HTMLDivElement>e.target).parentNode).style.color = ""
    (<HTMLElement>(<HTMLDivElement>e.target).parentNode).style.borderColor = ""
  }

  container.appendChild(drop_icon);
}


export function convertContainerToJSON(base_json, form_stack){
  // 1. Current Group build fields
  // 2. Current Group attatch responses
  // 3. Collect subgroups
    // R. Return to 1 and send the JSON over to be filled out
  // 4. JSON is filled out, return JSON
  let depth = 0
  while( form_stack.length != 0 ) {
    let level_order = 0
    let subgroup = form_stack.shift()
    let routes = subgroup.paths
    let first_route = routes[0]
    let current_form_field = base_json["FormFields"][first_route]
    for (let route_position = 1 ; route_position < routes.length ; route_position++){
      current_form_field = current_form_field["SubGroups"][routes[route_position]];
    }
    for (let c = 0 ; c < subgroup.node.childNodes.length ; c++){
      let child = subgroup.node.childNodes[c];
      if(child.nodeType != Node.ELEMENT_NODE){
        continue
      }
      if (child.className.indexOf("form-group") != -1){
        let r = routes
        r.push(level_order)
        form_stack.unshift( {node: child , paths: r} )
        level_order += 1;
        // If response_object Fields are not passed by reference then we'll have to do a reconstruction of the response object using current_form_field
        // Probably would be slow,  but is the first thing that comes to mind
        current_form_field["SubGroups"].push({
          Label: "",
          ID: "",
          Description: "",
          SubGroups: [],
          Respondables: [],
        })

        // a respondable group
      } else if(child.className.indexOf("respondable-container") != -1) {
        if(child.childNodes.length == 0 ){
          continue;
        }
        let response_list =  child.childNodes
        let return_list = []
        let len = response_list.length
        for (let res_index = 0 ; res_index < len ; res_index++){
          let response_child = response_list[res_index]
          current_form_field["Respondables"].push(convertRespondable( response_child ))
        }
      } else if((child as HTMLLabelElement).tagName.toUpperCase() == "LABEL"){
        let input_node:any = (child  as HTMLLabelElement).getElementsByTagName("INPUT")[0] as HTMLInputElement
        if (!input_node){
          input_node = (child as HTMLLabelElement).getElementsByTagName("SELECT")[0] as HTMLSelectElement
        }
        if (!input_node){
          input_node = (child as HTMLLabelElement).getElementsByTagName("TEXTAREA")[0] as HTMLTextAreaElement
        }
        if(!input_node){
          continue;
        }
        switch (input_node.getAttribute("name")) {
          case "form-label":
          current_form_field["Label"] = input_node.value

            break;
          case "id":
          current_form_field["ID"] = input_node.value

            break;
          case "description":
          current_form_field["Description"] = input_node.value

            break;

          default:
            break;
        }
      }
    }
  }
}

export function convertRespondable( response_child ){
  let res = {Object: {Field: { Name: "", Label: "", Required: false} } , Type: response_child.getAttribute("data-type")}
  let inputs = response_child.getElementsByTagName("INPUT")
  let selects = response_child.getElementsByTagName("SELECT")
  if(selects.length){
    switch (response_child.getAttribute("data-type")) {
      case "selectiongroup":
        res.Object[selects[0].getAttribute("Name")] = selects[0].value
        break;
      case "genericinput":
        res.Object[selects[0].getAttribute("Name")] = selects[0].value
        break;
      default:
        break;
    }
  }
  for (let field_index = 0 ; field_index < inputs.length ; field_index++){
    let val:any = inputs[field_index].value
    if (inputs[field_index].getAttribute("type") == "checkbox"){
       val = inputs[field_index].checked
    } else if(inputs[field_index].getAttribute("type") == "number"){
      val = parseInt(inputs[field_index].value)
    }
    if(response_child.getAttribute("data-type") == "selectiongroup"){
      // Selection groups
      // val could be boolean or string value
      if(inputs[field_index].getAttribute("data-list-item-no")) {
        if(!res.Object["CheckableItems"]){
          res.Object["CheckableItems"] = []
        }
        let checkitem_index = parseInt(inputs[field_index].getAttribute("data-list-item-no"))
        if(!res.Object["CheckableItems"][checkitem_index]){
          res.Object["CheckableItems"][checkitem_index] = {}
        }
        res.Object["CheckableItems"][checkitem_index][inputs[field_index].getAttribute("name")] =  val
      } else{
        switch (inputs[field_index].getAttribute("name")) {
          case "Name":
          case "Label":
          case "Required":
            res.Object.Field[inputs[field_index].getAttribute("name")] =  val
            break;
          default:
            res.Object[inputs[field_index].getAttribute("name")] =  val
            break;
        }
      }
    } else if(response_child.getAttribute("data-type") == "optiongroup"){
      // Option groups
      if(inputs[field_index].getAttribute("data-list-item-no")) {
        if(!res.Object["Options"]){
          res.Object["Options"] = []
        }
        let optitem_index = parseInt(inputs[field_index].getAttribute("data-list-item-no"))
        if(!res.Object["Options"][optitem_index]){
          res.Object["Options"][optitem_index] = {}
        }
        res.Object["Options"][optitem_index][inputs[field_index].getAttribute("name")] =  val
      } else{
        switch (inputs[field_index].getAttribute("name")) {
          case "Name":
          case "Label":
          case "Required":
            res.Object.Field[inputs[field_index].getAttribute("name")] =  val
            break;
          default:
            res.Object[inputs[field_index].getAttribute("name")] =  val
            break;
        }
      }
    } else{
      // Textarea and genericinput
      switch (inputs[field_index].getAttribute("name")) {
        case "Name":
        case "Label":
        case "Required":
          res.Object.Field[inputs[field_index].getAttribute("name")] =  val
          break;
        default:
          res.Object[inputs[field_index].getAttribute("name")] =  val
          break;
      }
    }
  }
  return res;
  // defined by UnmarshalerFormObject where object will be a name as key ->value as value type object
  //
}

export function createFormStack(node_list , response_object){
  let form_stack = []
  let first_depth_no = 0;
  for (let index = 0 ; index < node_list.length ; index++ ){
    if(node_list[index].nodeType != Node.ELEMENT_NODE){
      continue
    }
    if ((node_list[index] as HTMLElement).className.indexOf("form-group") != -1){
      form_stack.unshift( {node: (node_list[index] as HTMLElement) , paths: [first_depth_no]} )
      first_depth_no += 1;
      response_object["FormFields"].push({
        Label: "",
        ID: "",
        Description: "",
        SubGroups: [],
        Respondables: [],
      })
    }
  }
  return form_stack
}

function createTextAreaInputs(base_id:string, respondable_container_id: string ,  container: HTMLDivElement): string{
  let field_id = "field" + base_id
  let ta_id = "text-area" + ((Date.now() + Math.random())*10000)
  container.className =  "respondable-group feedback-group"
  container.id = ta_id
  container.setAttribute("data-type", "textarea")
  container.setAttribute("style" ,"width:400px;min-height:200px")
  container.innerHTML = `<SPAN>TextArea Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${field_id}-required" type="checkbox"/><BR/>
      </LI>
      <LI>
        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${field_id}-placeholder"/><BR/>
      </LI>
      <LI><BUTTON id="${ta_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${ta_id}')">Delete</BUTTON></LI>
    </UL>`;
    appendDropIcon(container)

    return field_id
}

function createFileInputs(base_id:string, respondable_container_id: string ,  container: HTMLDivElement): string{
  let field_id = "field" + base_id
  let fi_id = "file" + ((Date.now() + Math.random())*10000)
  container.className =  "respondable-group feedback-group"
  container.id = fi_id
  container.setAttribute("data-type", "fileinput")
  container.setAttribute("style" ,"width:400px;min-height:200px")
  container.innerHTML = `<SPAN>FileInput Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${field_id}-required" type="checkbox"/><BR/>
      </LI>
      <LI>
        Allowed Extention Pattern : <INPUT ondragstart="return false" draggable="false" data-field="1" name='AllowedExtRegex' id="${field_id}-allowed-ext"/><BR/>
      </LI>
      <LI>
        Max Filesize(Bytes) : <INPUT ondragstart="return false" draggable="false" data-field="1" type="number" name='MaxSize' id="${field_id}-max-size"/><BR/>
      </LI>
      <LI><BUTTON id="${fi_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${fi_id}')">Delete</BUTTON></LI>
    </UL>`;

    appendDropIcon(container)

    return field_id
}

function createInputInputs(base_id:string, respondable_container_id: string ,  container: HTMLDivElement){
  let field_id = "field" + base_id
  let in_id = "input" + ((Date.now() + Math.random())*10000)
  container.className =  "respondable-group feedback-group"
  container.id = in_id
  container.setAttribute("data-type", "genericinput")
  container.setAttribute("style" ,"width:400px;min-height:200px")
  container.innerHTML = `<SPAN>Input Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${field_id}-required" type="checkbox"/><BR/>
      </LI>
      <LI>
        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${field_id}-placeholder"/><BR/>
      </LI>
      <LI>
        Type : <SELECT data-field="1" name='Type' id="${field_id}-type" ondragstart="return false" draggable="false">
          <OPTGROUP label="Text Types">
            <OPTION value="text">Text</OPTION>
            <OPTION value="email">Email</OPTION>
            <OPTION value="number">Number</OPTION>
            <OPTION value="password">Password</OPTION>
            <OPTION value="url">URL</OPTION>
          </OPTGROUP>
          <OPTGROUP label="Time Types">
            <OPTION value="time">Time</OPTION>
            <OPTION value="date">Date</OPTION>
          </OPTGROUP>
          <OPTGROUP label="Special Types">
            <OPTION value="color">Color Picker</OPTION>
            <OPTION value="range">Number Range</OPTION>
          </OPTGROUP>
        </SELECT>
      </LI>
      <LI><BUTTON id="${in_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${in_id}')">Delete</BUTTON></LI>
    </UL>`;

    appendDropIcon(container)

    return field_id
}

function createSelectGroup(base_id:string, respondable_container_id:string ,  container: HTMLDivElement): string{
  let field_id = "field" + ((Date.now() + Math.random())*10000)
  let select_id = "select" + ((Date.now() + Math.random())*10000)
  container.className = "respondable-group feedback-group"
  container.id =  select_id
  container.setAttribute("data-type", "selectiongroup")
  container.setAttribute("style", "width:400px;min-height:200px")
  container.innerHTML = `<SPAN>SelectGroup Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${field_id}-required"/><BR/>
      </LI>
      <LI>
        Select type : <SELECT  data-field="1" name='SelectionCategory' type="checkbox" id="${field_id}-selectioncatergory" ondragstart="return false" draggable="false">
          <OPTION value="checkbox">checkbox</OPTION>
          <OPTION value="radio">radio</OPTION>
        </SELECT><BR/>
      </LI>
      <LI>
        Group Items :
        <BUTTON id="${field_id}" onclick="FormLibrary.addListField(this , 'checkable')" data-link-id="${field_id}">+</BUTTON>
        <BUTTON onclick="FormLibrary.removeListItem(this, 'checkable')" data-link-id="${field_id}">-</BUTTON><BR/>
        <OL data-field="1" data-select="1" id="${field_id}-checkable">
          <LI>
            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${field_id}-checkable-label-0"/>
            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${field_id}-checkable-value-0"/>
          </LI>
        </OL>
      </LI>
      <LI><BUTTON id="${select_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${select_id}')">Delete</BUTTON></LI>
    </UL>`;

    appendDropIcon(container)

    return field_id
}

function createOptionsGroup(base_id:string, respondable_container_id:string ,  container: HTMLDivElement) : string {
  let field_id = "field" + ((Date.now() + Math.random())*10000)
  let opt_id = "option" + ((Date.now() + Math.random())*10000)
  container.className = "respondable-group feedback-group"
  container.id =  opt_id
  container.setAttribute("data-type", "optiongroup")
  container.setAttribute("style", "width:400px;min-height:200px")
  container.innerHTML = `<SPAN>OptionGroup Creation Info:</SPAN><BR/>
    <UL>
      <LI>
        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${field_id}-required"/><BR/>
      </LI>
      <LI>
      Option Fields :
        <BUTTON id="${field_id}" onclick="FormLibrary.addListField(this , 'option')" data-link-id="${field_id}">+</BUTTON>
        <BUTTON onclick="FormLibrary.removeListItem(this, 'option')" data-link-id="${field_id}">-</BUTTON><BR/>
        <OL data-field="1" data-select="1" id="${field_id}-option">
          <LI>
            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${field_id}-option-label-0"/>
            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${field_id}-option-value-0"/>
          </LI>
        </OL>
      </LI>
      <LI><BUTTON id="${opt_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${opt_id}')">Delete</BUTTON></LI>
    </UL>`;

    appendDropIcon(container)

    return field_id
}
