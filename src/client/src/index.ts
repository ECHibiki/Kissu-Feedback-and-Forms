What's left to do...:

correct this file
Submit a form
hndle errors
Render the form for a client
Allow client to respond
Allow mod to view responses(forms, form, replies)
Upload beta version to site
Finish implementing other fields, create and recieve
Allow form edits
Allow downloads
Allow deletes, form and reply

API for displaying data
Remake this file with react components and a better mod view mode

// the API file for displaying various views
// I'm gonna have to not hardcode this later... consider this file essentially a draft

export function helloWorld() {
  console.log("hello client");
}

export function attatchCreators(){
  let root_group_button = document.getElementById("head-root-group-create")
  root_group_button.onclick = createNewSubgroup(root_group_button , 0)

  let submit_button = document.getElementById("form-submit-button")
  submit_button.onclick = submitForm(submit_button )
}

function createNewSubgroup( button: HTMLButtonElement, base_no: number  ){
  let parent_container =  button.parentNode
  base_no += 1
  let container = document.createElement('DIV');
  container.className =  "subgroup form-group"
  container.style= "border:1px solid black;"
  container.innerHTML = `
  <P>Required Fields:</P>
  <LABEL>Form Name : <INPUT type="text" name="form-name"/> </LABEL> <br/>
  <LABEL>Form ID : <INPUT type="text" name="id"/> <br/>
  <LABEL>Form Descriptor : <INPUT type="text" name="description"/> <br/>
  <LABEL>Allow Anonymous Submissions :<INPUT type="checkbox" name="anon-option"/>
  <SPAN data-depth="${base_no}" class="respondable-container"> <SPAN><BR/>
  <BUTTON  onclick="FormLibrary.createNewResponseElement(this , ${base_no})">Create New Respondable Below Last Respondable</BUTTON><br/>
  <BUTTON  onclick="FormLibrary.createNewSubgroup(this , ${base_no})">Create New Group Below Last Respondable</BUTTON><br/>
  `
  parent_container.appendChild(container)
}

function createNewResponseElement(button: HTMLButtonElement , base_no: number ){
  respondable_container = button.previousSibling.previousSibling

  let container = document.createElement('DIV');
  container.className =  "creation-prompt"
  container.style= "border:1px solid black;width:400px;height:200px"
  container.innerHTML = `Element Creation Info:<BR/>
    <UL>
      <LI>Item Type:
        <SELECT>
          <OPTION value="textarea">TextArea</OPTION>
          <OPTION value="input">Input(unimplemented)</OPTION>
          <OPTION value="file">FileInput(unimplemented)</OPTION>
          <OPTION value="select">SelectGroup</OPTION>
          <OPTION value="option">OptionGroup(unimplemented)</OPTION>
        </SELECT>
      </LI>
      <LI><BUTTON onclick="FormLibrary.responseTypeSelected(this, ${base_no})">Next</BUTTON></LI>
    </UL>`

    // CREATE NEW AT END
    // INSERT IN POSITION
    respondable_container.appendChild(container)
}

function createTextAreaInputs(container: HTMLDivElement, base_no: number){
  container.className =  "field-container"
  container.style= "border:1px solid black;width:400px;height:200px"
  container.innerHTML = `TextArea Creation Info:<BR/>
  <INPUT data-depth="${base_no}" type="hidden" name="type" value="textarea" />
    <UL>
      <LI>
        Name : <INPUT  data-depth="${base_no}" data-field="1" name='name'/><BR/>
      </LI>
      <LI>
        Label : <INPUT data-depth="${base_no}" data-field="1" name='Label'/><BR/>
      </LI>
      <LI>
        Required : <INPUT data-depth="${base_no}" data-field="1" name='Required' type="checkbox"/><BR/>
      </LI>
      <LI>
        Placeholder : <INPUT data-depth="${base_no}" data-field="1" name='placeholder'/><BR/>
      </LI>
    </UL>`;
}

function createSelectGroup(container: HTMLDivElement , base_no: number){
  container.className =  "field-container"
  container.style= "border:1px solid black;width:400px;height:200px"
  container.innerHTML = `TextArea Creation Info:<BR/>
    <INPUT type="hidden" name="type" value="selectiongroup" />
    <UL>
      <LI>
        Name : <INPUT data-depth="${base_no}" data-field="1" name='name'/><BR/>
      </LI>
      <LI>
        Label : <INPUT data-depth="${base_no}" data-field="1" name='Label'/><BR/>
      </LI>
      <LI>
        Required : <INPUT data-depth="${base_no}" data-field="1" name='Required' type="checkbox"/><BR/>
      </LI>
      <LI>
        Select type : <SELECT data-depth="${base_no}" data-field="1" name='SelectionCategory' type="checkbox">
          <OPTION value="checkbox">checkbox</OPTION>
          <OPTION value="radio">radio</OPTION>
        </SELECT><BR/>
      </LI>
      <LI>
        Group Items :
        <BUTTON onclick="FormLibrary.addCheckable(this)">+</BUTTON>
        <BUTTON onclick="FormLibrary.removeCheckable(this)">-</BUTTON><BR/>
        <OL data-field="1" data-select="1">
          <LI>
            <INPUT placeholder="Label" data-depth="${base_no}" data-checkable-no="1" data-field="1" name='Label'/>
            <INPUT placeholder="Value" data-depth="${base_no}" data-checkable-no="1" data-field="1" name='Value'/>
          </LI>
        </OL>
      </LI>
    </UL>`;
}

function addCheckable(button : HTMLButtonElement) {
  let ol = button.nextSibling.nextSibling.nextSibling
  let li = document.createElement("LI")
  let child_count = li.children.length
  li.InnterHTML = `
    <INPUT placeholder="Label" data-depth="${base_no}" data-checkable-no="${child_count.length}" data-field="1" name='Label'/>
    <INPUT placeholder="Value" data-depth="${base_no}" data-checkable-no="${child_count.length}" data-field="1" name='Value'/>
  `
  ol.appendChild(li)
}
function removeCheckable(button : HTMLButtonElement) {
  let ol = button.nextSibling.nextSibling
  ol.removeChild(ol.lastChild)
}

function responseTypeSelected(button: HTMLButtonElement , base_no: number){
  let select_obj = button.parentNode.previousSibling.firstChild
  let container = button.parentNode.parentNode.parentNode
  let selection = select_obj.value
  switch (selection) {
    case "textarea":
      createTextAreaInputs(container , base_no)
      break;
    case "select":
      createSelectGroup(container , base_no)

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

function submitForm(button: HTMLButtonElement){
  let form_container = button.parentNode
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
  for (let child in form_container.firstChild.childNodes){
    if (child.className.indexOf("form-group") != -1){
      form_stack.unshift( {node: child , paths: [0]} )
      response_object["FormFields"].push({
        FormName: "",
        ID: "",
        Description: "",
        SubGroup: [],
        Respondables: []],
      })
    } else{
      switch (child.getAttribute("name")) {
        case "form-name":
        response_object["FormName"] = child.value

          break;
        case "id":
        response_object["ID"] = child.value

          break;
        case "description":
        response_object["Description"] = child.value

          break;
        case "anon-option":
        response_object["AnonOption"] = child.value

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
    let current_form_field = response_object["FormFields"][0]
    routes.forEach((route:number) => {
      current_form_field = current_form_field["SubGroup"][route]
    });
    for (let c = 0 ; c < subgroup.node.childNodes.length ; c++){
      let child = subgroup.node.childNodes[c];
      if (child.className.indexOf("form-group") != -1){
        let r = routes
        r.push(level_order)
        form_stack.unshift( {node: child , paths: r} )
        level_order += 1;
        // If response_object Fields are not passed by reference then we'll have to do a reconstruction of the response object using current_form_field
        // Probably would be slow,  but is the first thing that comes to mind
        current_form_field["SubGroup"].push({
          FormName: "",
          ID: "",
          Description: "",
          SubGroup: [],
          Respondables: [],
        })

        // a respondable group
      } else if(child.className.indexOf("respondable-container") != -1) {
        for (let c = 0 ; c < child.childNodes.length ; c++){
          if (!child.childNodes[c].getAttribute("data-field") ){
            continue
          }
          let res = {}
          for (let res_index = 0 ; res_index < child.childNodes[c].childNodes.length ; res_index++){
            let r = child.childNodes[c].childNodes[res_index]
            if ( r.getAttribute("data-select") ){
              res["CheckableItems"] = {}
              //label
              res["CheckableItems"]["Label"] = r.value
              //value
              res["CheckableItems"]["Value"] = r.value
            } else{
              res[r.getAttribute("data-name")] = r.value
            }
          }
          current_form_field["Respondables"].push(res)
        }
        // defined by UnmarshalerFormObject where object will be a name as key ->value as value type object
        //
      } else{
        switch (child.getAttribute("name")) {
          case "form-name":
          current_form_field["FormName"] = child.value

            break;
          case "id":
          current_form_field["ID"] = child.value

            break;
          case "description":
          current_form_field["Description"] = child.value

            break;

          default:
            break;
        }
      }
    }
  }

  // write as json and send on "form-construct-json"
  let construct_string = JSON.stringify(response_object)
  sendCreateRequest(construct_string)
}

function sendCreateRequest(form: string){
  var f = new formData()
  f.append("json" , "1")
  f.append("form-construct-json" , form )

  var x = new XMLHttpRequest()
  x.open("POST" , "/mod/create" , true)
  x.load = handleCreateComplete
  x.send(f)
}

function handleCreateComplete(e) {
  console.log(e)
}
