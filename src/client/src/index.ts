// What's left to do...:
//
// correct this file
// Submit a form
// hndle errors
//
// Render the form for a client
// Allow client to respond
//
// Allow mod to view responses(forms, form, replies)
//
// Upload beta version to site
// Finish implementing other fields, create and recieve
// Allow form edits
// Allow downloads
// Allow deletes on form and reply
//
// API for displaying data
// Remake this file with react components and a better mod view mode

// the library file for displaying various views
// I'm gonna have to not hardcode(eg. nextSibling) this later... consider this file essentially a draft

export function helloWorld() {
  console.log("hello client");
}

export function attatchCreators(settings:any){
  console.log(settings)
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
    let form_construct = JSON.parse(raw_json)
    let field_structure = { SubGroups: form_construct.FormFields }
    console.log("FC", field_structure)
    let reference_button = document.getElementById("sub-create");
    recursiveFormFieldRebuild(field_structure , <HTMLButtonElement>reference_button)
  } catch (error) {
    console.error(error)
    alert("Issue with rebuilding")
  }
}

function recursiveFormFieldRebuild(field_group:any , base_button:HTMLButtonElement){
  console.log(field_group , base_button)
  if(!field_group){
    return
  }

  let respondables = field_group.Respondables
  console.log("RES", respondables)
  if(respondables){
    respondables.forEach(field => {
      console.log(field)
      let respondable_creation_object = createNewResponseElement(base_button)
      let field_id = ""
      switch (field.Type) {
        case "textarea":
          field_id = createTextAreaInputs(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id ,respondable_creation_object.container);
          (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
          (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
          (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;
          (<HTMLInputElement>document.getElementById(field_id + "-placeholder")).value = field.Object.Placeholder;
          break;
        case "selectiongroup":
          field_id = createSelectGroup(respondable_creation_object.base_id, respondable_creation_object.respondable_container_id , respondable_creation_object.container);
          (<HTMLInputElement>document.getElementById(field_id + "-name")).value = field.Object.Field.Name;
          (<HTMLInputElement>document.getElementById(field_id + "-label")).value = field.Object.Field.Label;
          (<HTMLInputElement>document.getElementById(field_id + "-required")).checked = field.Object.Field.Required;

          for(let check_count = 1 ; check_count < field.Object.CheckableItems.length ; check_count++ ){
             addCheckable( <HTMLButtonElement>document.getElementById(field_id) )
          }
          field.Object.CheckableItems.forEach((check , index) => {
            (<HTMLInputElement>document.getElementById(field_id + "-checkable-label-" + index)).value = check.Label;
            (<HTMLInputElement>document.getElementById(field_id + "-checkable-value-" + index)).value = check.Value;
          });

          break;
        default:
          break;
      }
    });
  }
  setFieldsAsDisabled(true)

  let group_structures = field_group.SubGroups
  console.log("GS" , group_structures)
  if(group_structures){
    group_structures.forEach(structure => {
      let button_ref = createNewSubgroup(base_button);
      (<HTMLInputElement>document.getElementById(button_ref + "-label")).value = structure.Label;
      (<HTMLInputElement>document.getElementById(button_ref + "-id")).value = structure.ID;
      (<HTMLInputElement>document.getElementById(button_ref + "-description")).value = structure.Description;

      let reference_button = document.getElementById(button_ref);
      recursiveFormFieldRebuild(structure , <HTMLButtonElement>reference_button)
    });
  }
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


export function createNewSubgroup( button: HTMLButtonElement): string{
  let parent_id = button.getAttribute("data-link-id")
  let parent_container =  document.getElementById(parent_id + "-group")
  console.log(parent_id , parent_container , button)

  let group_id = "group" + (Date.now() + Math.random())
  let container = document.createElement('DIV');
  container.className =  "subgroup form-group"
  container.id =  group_id + "-group"
  container.setAttribute("style" , "")
  container.innerHTML = `<LABEL>Group Label : <INPUT type="text" name="form-label" id="${group_id}-label"/> </LABEL> <br/>
  <LABEL>Group ID : <INPUT type="text" name="id" id="${group_id}-id"/></LABEL> <br/>
  <LABEL>Form Descriptor : <INPUT type="text" name="description" id="${group_id}-description"/></LABEL> <br/>
  <BUTTON id="${group_id}" onclick="FormLibrary.createNewResponseElement(this)" data-link-id="${group_id}" >Create New Respondable Below Last Respondable</BUTTON><br/>
  <BUTTON  onclick="FormLibrary.createNewSubgroup(this)" data-link-id="${group_id}">Create New Group Below Last Respondable</BUTTON><br/>
  <BUTTON  onclick="FormLibrary.deleteContainer('${parent_id + "-group"}' , '${ group_id + "-group"}')">Delete Subgroup</BUTTON><br/>
  <SPAN  class="respondable-container" id="${group_id}-respondables"></SPAN><BR/>
  `
  parent_container.appendChild(container)
  return group_id
}

export function createNewResponseElement(button: HTMLButtonElement ): ({base_id:string, respondable_container_id:string, container:HTMLDivElement }){
  let parent_id = button.getAttribute("data-link-id")
  let respondable_container =  document.getElementById(parent_id + "-respondables")
  console.log(parent_id , respondable_container)

  let res_id = "response" + (Date.now() + Math.random())
  let container = document.createElement('DIV');
  container.className =  "creation-prompt"
  container.id =  res_id + "-fields"
  container.setAttribute("style", "width:400px;min-height:200px")
  container.innerHTML = `Element Creation Info:<BR/>
    <UL>
      <LI>Item Type:
        <SELECT id="${res_id}-type">
          <OPTION value="textarea">TextArea</OPTION>
          <OPTION value="input">Input(unimplemented)</OPTION>
          <OPTION value="file">FileInput(unimplemented)</OPTION>
          <OPTION value="select">SelectGroup</OPTION>
          <OPTION value="option">OptionGroup(unimplemented)</OPTION>
        </SELECT>
      </LI>
      <LI>If we want any of the unimplemented features, then you'll have to ask me or wait until I personally require it</LI>
      <LI><BUTTON id="${res_id}" data-link-id="${res_id}" onclick="FormLibrary.responseTypeSelected('${parent_id + "-respondables"}' , this)">Next</BUTTON></LI>
      <LI><BUTTON  data-link-id="${res_id}" onclick="FormLibrary.deleteContainer('${parent_id + "-respondables"}' , '${res_id + "-fields"}')">Delete</BUTTON></LI>
    </UL>`

    // CREATE NEW AT END
    // INSERT IN POSITION
    respondable_container.appendChild(container)
    return { base_id: res_id, respondable_container_id: parent_id + "-respondables" , container: <HTMLDivElement>container  }
}

export function deleteContainer(base_container_id:string , sub_container_id:string){
  document.getElementById(base_container_id).removeChild(document.getElementById( sub_container_id));
}

export function createTextAreaInputs(base_id:string, respondable_container_id: string ,  container: HTMLDivElement): string{
  let field_id = "field" + base_id
  let ta_id = "text-area" + (Date.now() + Math.random())
  container.className =  "respondable-group"
  container.id = ta_id
  container.setAttribute("data-type", "textarea")
  container.setAttribute("style" ,"width:400px;min-height:200px")
  container.innerHTML = `TextArea Creation Info:<BR/>
    <UL>
      <LI>
        Name : <INPUT   data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT  data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT  data-field="1" name='Required' id="${field_id}-required" type="checkbox"/><BR/>
      </LI>
      <LI>
        Placeholder : <INPUT  data-field="1" name='Placeholder' id="${field_id}-placeholder"/><BR/>
      </LI>
      <LI><BUTTON id="${ta_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}' , '${ta_id}')">Delete</BUTTON></LI>
    </UL>`;
    return field_id
}

export function createSelectGroup(base_id:string, respondable_container_id:string ,  container: HTMLDivElement): string{
  let field_id = "field" + (Date.now() + Math.random())
  let select_id = "select" + (Date.now() + Math.random())
  container.className =  "field-container"
  container.id =  select_id
  container.setAttribute("data-type", "selectiongroup")
  container.setAttribute("style", "width:400px;min-height:200px")
  container.innerHTML = `TextArea Creation Info:<BR/>
    <UL>
      <LI>
        Name : <INPUT  data-field="1" name='Name' id="${field_id}-name"/><BR/>
      </LI>
      <LI>
        Label : <INPUT  data-field="1" name='Label' id="${field_id}-label"/><BR/>
      </LI>
      <LI>
        Required : <INPUT  data-field="1" name='Required' type="checkbox" id="${field_id}-required"/><BR/>
      </LI>
      <LI>
        Select type : <SELECT  data-field="1" name='SelectionCategory' type="checkbox" id="${field_id}-selectioncatergory">
          <OPTION value="checkbox">checkbox</OPTION>
          <OPTION value="radio">radio</OPTION>
        </SELECT><BR/>
      </LI>
      <LI>
        Group Items :
        <BUTTON id="${field_id}" onclick="FormLibrary.addCheckable(this)" data-link-id="${field_id}">+</BUTTON>
        <BUTTON onclick="FormLibrary.removeCheckable(this)" data-link-id="${field_id}">-</BUTTON><BR/>
        <OL data-field="1" data-select="1" id="${field_id}-checkable">
          <LI>
            <INPUT placeholder="Label"  data-checkable-no="0" data-field="1" name='Label' id="${field_id}-checkable-label-0"/>
            <INPUT placeholder="Value"  data-checkable-no="0" data-field="1" name='Value'  id="${field_id}-checkable-value-0"/>
          </LI>
        </OL>
      </LI>
      <LI><BUTTON id="${select_id}" onclick="FormLibrary.deleteContainer('${respondable_container_id}" , '${select_id}")">Delete</BUTTON></LI>
    </UL>`;
    return field_id
}

export function addCheckable(button : HTMLButtonElement) {
  let parent_id = button.getAttribute("data-link-id")
  let ol = document.getElementById(parent_id + "-checkable")
  let li = document.createElement("LI")
  let child_count = ol.children.length
  li.innerHTML = `
    <INPUT placeholder="Label" data-checkable-no="${child_count}" data-field="1" name='Label' id="${parent_id}-checkable-label-${child_count}"/>
    <INPUT placeholder="Value" data-checkable-no="${child_count}" data-field="1" name='Value' id="${parent_id}-checkable-value-${child_count}"/>
  `
  ol.appendChild(li)
}
export function removeCheckable(button : HTMLButtonElement) {
  let parent_id = button.getAttribute("data-link-id")
  let ol = document.getElementById(parent_id + "-checkable")
  if (ol.childNodes.length <= 1) {
    return
  }
  ol.removeChild(ol.lastChild)
}

export function responseTypeSelected(respondable_container_id:string , button: HTMLButtonElement ){
  let base_id = button.getAttribute("data-link-id")
  let select_obj = document.getElementById(base_id + "-type")
  let container = document.getElementById(base_id + "-fields") as HTMLDivElement
  console.log(base_id , select_obj , container)
  let selection = (select_obj as HTMLInputElement).value
  switch (selection) {
    case "textarea":
      createTextAreaInputs(base_id , respondable_container_id , container)
      break;
    case "select":
      createSelectGroup(base_id , respondable_container_id , container)

      break;

    default:
      break;
  }
}

/*
type FormConstruct struct {
  FormName string
  ID string
  Description string
  // With anon option set to true there is an ability for users to flag themselves as anonymous
  // Under other conditions there is no anonymity
  AnonOption bool
  FormFields []FormGroup
}

  type FormGroup struct {
    Label string
    ID string
    Description string
    SubGroup []FormGroup
    Respondables []UnmarshalerFormObject
  }

  //currently a placeholder to gain polymorphic properties
  type FormObject interface{
      ElementType() string
      GetName() string
      GetDescription() string
      GetRequired() bool
  }

  // implement Unmarshaler interface
  type UnmarshalerFormObject struct {
    Type FormObjectTag
    Object FormObject
  }
*/

export function submitForm(button: HTMLButtonElement , post_url){
  let response_object = {
    FormName: "",
    ID: "",
    Description: "",
    AnonOption: false,
    FormFields: [ ]
  }


  //  send on c.PostForm("form-construct-json")
  // marshals into type FormConstruct struct
  // get the first entry div and fill out FormConstruct's default params
  let form_stack = []
  let first_depth_no = 0;
  let node_list = document.getElementById("head-group").childNodes
  // console.log(node_list)
  for (let index = 0 ; index < node_list.length ; index++ ){
    if(node_list[index].nodeType != Node.ELEMENT_NODE){
      continue
    }
    // console.log(node_list[index].nodeType , Node.ELEMENT_NODE , node_list[index])
    if ((node_list[index] as HTMLElement).className.indexOf("form-group") != -1){
      form_stack.unshift( {node: (node_list[index] as HTMLElement) , paths: [first_depth_no]} )
      first_depth_no += 1;
      response_object["FormFields"].push({
        Label: "",
        ID: "",
        Description: "",
        SubGroup: [],
        Respondables: [],
      })
    } else if((node_list[index] as HTMLLabelElement).tagName.toUpperCase() == "LABEL"){
      let input_node = (node_list[index]  as HTMLLabelElement).getElementsByTagName("INPUT")[0] as HTMLInputElement
      // console.log(node_list[index] , input_node)
      switch (input_node.getAttribute("name")) {
        case "form-name":
        response_object["FormName"] = input_node.value

          break;
        case "id":
        response_object["ID"] = input_node.value

          break;
        case "description":
        response_object["Description"] = input_node.value

          break;
        case "anon-option":
        response_object["AnonOption"] = input_node.checked

          break;

        default:
          break;
      }
    }
  }
  // Step through the first div
    // collect the fields defining the subgroup by type FormGroup struct
    // get the respondables span container <INPUT type="hidden" name="type" value="textarea" />
    // defined by UnmarshalerFormObject where object will be a name as key ->value as value type object

  let depth = 0
  while( form_stack.length != 0 ) {
    let level_order = 0
    let subgroup = form_stack.shift()
    let routes = subgroup.paths
    let first_route = routes[0]
    let current_form_field = response_object["FormFields"][first_route]
    for (let route_position = 1 ; route_position < routes.length ; route_position++){
      current_form_field = current_form_field["SubGroup"][routes[route_position]];
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
        current_form_field["SubGroup"].push({
          Label: "",
          ID: "",
          Description: "",
          SubGroup: [],
          Respondables: [],
        })

        // a respondable group
      } else if(child.className.indexOf("respondable-container") != -1) {
        if(child.childNodes.length == 0 ){
          continue;
        }
        let response_list =  child.childNodes
        let len = response_list.length
        console.log(response_list , len)
        for (let res_index = 0 ; res_index < len ; res_index++){
          let response_child = response_list[res_index]
          let res = {Object: {Field: { Name: "", Label: "", Required: false} } , Type: response_child.getAttribute("data-type")}
          let inputs = response_child.getElementsByTagName("INPUT")
          let selects = response_child.getElementsByTagName("SELECT")
          if(selects.length){
            switch (response_child.getAttribute("data-type")) {
              case "selectiongroup":
                res.Object[selects[0].getAttribute("Name")] = selects[0].value
                break;

              default:
                break;
            }
          }
          for (let field_index = 0 ; field_index < inputs.length ; field_index++){
            if(response_child.getAttribute("data-type") == "selectiongroup"){
              // Selection groups
              let val:any = inputs[field_index].value
              if (inputs[field_index].getAttribute("type") == "checkbox"){
                 val = inputs[field_index].checked
              }

              if(inputs[field_index].getAttribute("data-checkable-no")) {
                if(!res.Object["CheckableItems"]){
                  res.Object["CheckableItems"] = []
                }
                let checkitem_index = parseInt(inputs[field_index].getAttribute("data-checkable-no"))
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
            } else if(response_child.getAttribute("data-type") == "fileinput"){
              // FileInputTag
            } else{
              // Textarea and etc
              let val:any = inputs[field_index].value
              if (inputs[field_index].getAttribute("type") == "checkbox"){
                 val = inputs[field_index].checked
              }
              console.log(inputs , inputs[field_index].getAttribute("name") ,  val)
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
          current_form_field["Respondables"].push(res)
        }
        // defined by UnmarshalerFormObject where object will be a name as key ->value as value type object
        //
      } else if((child as HTMLLabelElement).tagName.toUpperCase() == "LABEL"){
        let input_node = (child  as HTMLLabelElement).getElementsByTagName("INPUT")[0] as HTMLInputElement
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

  // write as json and send on "form-construct-json"

  console.log(response_object)
  let construct_string = JSON.stringify(response_object)
  sendCreateRequest(construct_string , post_url)
}

export function sendCreateRequest(form: string , url:string){
  var f = new FormData()
  f.append("json" , "1")
  f.append("form-construct-json" , form )

  var x = new XMLHttpRequest()
  x.open("POST" , url , true)
  x.onload = handleCreateComplete
  x.send(f)
}

export function handleCreateComplete(e) {
  let response_json:any = {}
  try {
    response_json = JSON.parse(e.originalTarget.responseText)
  } catch (error) {
    document.getElementById("response-container").innerHTML = `Hard server crash`
    return
  }
  console.log(response_json , e.originalTarget.status)
  if(response_json.error){
    if(!response_json['error-list']){
      document.getElementById("response-container").innerHTML = "Server error: " + response_json["error"]
    }
    let error_message =  "There are some issues preventing you from submitting...<br/><ul>"
    response_json['error-list'].forEach((err) => {
      error_message += "<li>" + err.FailPosition + " : " + err.FailType + "</li>"
    });
    error_message += "</ul>"
    document.getElementById("response-container").innerHTML = error_message

  } else{
    document.getElementById("response-container").innerHTML = `${response_json.message}...<br/>
      URL: <a href="${response_json.URL}">${response_json.URL}</a>
    `
  }
}

export function postDeleteForm(name: string , number:string){
  let conf = confirm('This will remove the form from the display and database.\n\
    Information files will retained on the server as a record until a duplicate named form is created')
  if (conf == true){
    let x = new XMLHttpRequest()
    x.open("POST" , `/mod/form/delete/${name}/${number}`)
    x.onload = (e:any) => {
      let response_json:any = {}
      try {
        response_json = JSON.parse(e.originalTarget.responseText)
      } catch (error) {
        alert("Error: Hard server crash")
        return
      }
      if(response_json.error){
        alert("Error: " + response_json.error)
      } else{
        alert("Success: " + response_json.message)
        let li = document.getElementById(`row-${name}-${number}`)
        //this is fine... there is nothing else that an LI can be contained by
        li.parentNode.removeChild(li)
      }
    }
    x.send()
    return false;
  }
}
export function postDeleteResponse(form_name: string , response_number:string){
  let conf = confirm('This will remove the response and all information associated with it')
  if (conf == true){
    let x = new XMLHttpRequest()
    x.open("POST" , `/mod/response/delete/${form_name}/${response_number}`)
    x.onload = (e:any) => {
      let response_json:any = {}
      try {
        response_json = JSON.parse(e.originalTarget.responseText)
      } catch (error) {
        alert("Error: Hard server crash")
        return
      }
      if(response_json.error){
        alert("Error: " + response_json.error)
      } else{
        alert("Success: " + response_json.message)
        let li = document.getElementById(`row-${form_name}-${response_number}`)
        //this is fine... there is nothing else that an LI can be contained by
        li.parentNode.removeChild(li)
      }
    }
    x.send()
    return false;
  }
}
