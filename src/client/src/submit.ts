import {createFormStack , convertContainerToJSON , buildDropdownResponse} from "./builder"

export function submitUserPost(form){
  let fields = new FormData(form)
  fields.append("json" , "1")
  let x = new XMLHttpRequest()
  x.open("POST" , "" + window.location)
  x.onload = handleCreateComplete
  x.send(fields)

  return false;
}

export function postDeleteResponse(form_name: string , response_number:string){
  let conf = confirm('This will remove the response and all information associated with it')
  if (conf == true){
    let x = new XMLHttpRequest()
    x.open("POST" , `/mod/response/delete/${form_name}/${response_number}`)
    x.onload = (e:any) => {
      let response_json:any = {}
      try {
        response_json = JSON.parse(e.target.responseText)
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
  }
  return false;
}

function sendCreateRequest(form: string , url:string){
  var f = new FormData()
  f.append("json" , "1")
  f.append("form-construct-json" , form )

  var x = new XMLHttpRequest()
  x.open("POST" , url , true)
  x.onload = handleCreateComplete
  x.send(f)
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
        response_json = JSON.parse(e.target.responseText)
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
  }
  return false;
}

function handleCreateComplete(e) {
  let response_json:any = {}
  try {
    response_json = JSON.parse(e.target.responseText)
  } catch (error) {
    document.getElementById("response-container").innerHTML = `Hard server crash<BR/><TEXTAREA>${error.toString()}</TEXTAREA>`
    return
  }
  if(response_json.error){
    if(!response_json['error-list']){
      document.getElementById("response-container").innerHTML = "Server error: " + response_json["error"]
    } else{
      let error_message =  "There are some issues preventing you from submitting...<br/><ul>"
      response_json['error-list'].forEach((err) => {
        error_message += "<li>" + err.FailPosition + " : " + err.FailType + "</li>"
      });
      error_message += "</ul>"
      document.getElementById("response-container").innerHTML = error_message
    }

  } else{
    document.getElementById("response-container").innerHTML = `${response_json.message} - ${Date.now()}<br/>
      ${response_json.URL ? `URL: <a href="${response_json.URL}">${response_json.URL}</a>` : "" }
    `
  }
}

export function dropdownList(activator:HTMLAnchorElement, id:string) : boolean{
  let head = document.getElementById("container-" + id);
  let reply_container = head.getElementsByClassName("item-replies").item(0)

  if(activator.firstChild.textContent == "▶"){
    activator.firstChild.textContent = "( ...Loading... )";
    let x = new XMLHttpRequest()
    x.open("GET" , "/mod/api/form/" + id );
    x.onload = function(e:any){
      activator.firstChild.textContent = "▼";
      let response_json:any;
      try {
        response_json = JSON.parse(e.target.responseText)
      } catch (error) {
        reply_container.innerHTML = `Hard server crash<BR/><TEXTAREA>${error.toString()}</TEXTAREA>`
        return
      }
      if (response_json.error){
        reply_container.innerHTML = response_json.error
      } else{
        reply_container.innerHTML = buildDropdownResponse(reply_container as HTMLDivElement, response_json);
      }
    }
    x.onerror= function(){
      reply_container.innerHTML = `Server Issue`;
    }
    x.send();
  } else if(activator.firstChild.textContent == "▼"){
    (activator.firstChild).textContent = "▶";
    reply_container.innerHTML = "";
  }
  return false;
}

export function submitForm(button: HTMLButtonElement , post_url){
  let response_object = {
    FormName: "",
    ID: "",
    Description: "",
    AnonOption: false,
    FormFields: [ ]
  }

  let blank = document.querySelector("div[data-type=blank]")
  if( blank ){
    document.getElementById("response-container").innerHTML = "Formatting: You have an unfinished respondable. Delete or complete it.";
    return;
  }

  //  send on c.PostForm("form-construct-json")
  // marshals into type FormConstruct struct
  // get the first entry div and fill out FormConstruct's default params
  let node_list = document.getElementById("head-group").childNodes
  for (let index = 0 ; index < node_list.length ; index++ ){
    if(node_list[index].nodeType != Node.ELEMENT_NODE){
      continue
    }
    if((node_list[index] as HTMLLabelElement).tagName.toUpperCase() == "LABEL"){
      let input_node:any = (node_list[index] as HTMLLabelElement).getElementsByTagName("INPUT")[0] as HTMLInputElement

      if (!input_node){
        input_node = (node_list[index] as HTMLLabelElement).getElementsByTagName("SELECT")[0] as HTMLSelectElement
      }
      if (!input_node){
        input_node = (node_list[index] as HTMLLabelElement).getElementsByTagName("TEXTAREA")[0] as HTMLTextAreaElement
      }
      if (!input_node){
        continue;
      }
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
  let form_stack = createFormStack(node_list , response_object)
  convertContainerToJSON(response_object , form_stack)
  // write as json and send on "form-construct-json"

  let construct_string = JSON.stringify(response_object)
  sendCreateRequest(construct_string , post_url)
}
