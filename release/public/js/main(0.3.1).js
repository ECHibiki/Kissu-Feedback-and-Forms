!function webpackUniversalModuleDefinition(e,t){"object"==typeof exports&&"object"==typeof module?module.exports=t():"function"==typeof define&&define.amd?define("FormLibrary",[],t):"object"==typeof exports?exports.FormLibrary=t():e.FormLibrary=t()}(self,(()=>(()=>{"use strict";var e={};return(()=>{var t,n=e;function recursiveFormFieldRebuild(e,t,n){if(!e)return;let a=e.Respondables;a&&a.forEach((e=>{rebuildRespondable(e,t)})),setFieldsAsDisabled(!0);let r=e.SubGroups;r&&r.forEach((e=>{let a=createNewSubgroup(t,n);document.getElementById(a+"-label").value=e.Label,document.getElementById(a+"-id").value=e.ID,document.getElementById(a+"-description").value=e.Description,recursiveFormFieldRebuild(e,document.getElementById(a))}))}function rebuildRespondable(e,t,n){let a=createNewResponseElement(t,n),r="";switch(e.Type){case"textarea":r=createTextAreaInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-placeholder").value=e.Object.Placeholder;break;case"genericinput":r=createInputInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-placeholder").value=e.Object.Placeholder,document.getElementById(r+"-type").value=e.Object.Type;break;case"fileinput":r=createFileInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-allowed-ext").value=e.Object.AllowedExtRegex,document.getElementById(r+"-max-size").value=e.Object.MaxSize;break;case"selectiongroup":r=createSelectGroup(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required;for(let t=1;t<e.Object.CheckableItems.length;t++)addListField(document.getElementById(r),"checkable");e.Object.CheckableItems.forEach(((e,t)=>{document.getElementById(r+"-checkable-label-"+t).value=e.Label,document.getElementById(r+"-checkable-value-"+t).value=e.Value}));break;case"optiongroup":r=createOptionsGroup(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required;for(let t=1;t<e.Object.Options.length;t++)addListField(document.getElementById(r),"option");e.Object.Options.forEach(((e,t)=>{document.getElementById(r+"-option-label-"+t).value=e.Label,document.getElementById(r+"-option-value-"+t).value=e.Value}))}}function setFieldsAsDisabled(e){let t=Array.from(document.getElementsByTagName("BUTTON")),n=Array.from(document.getElementsByTagName("INPUT")),a=Array.from(document.getElementsByTagName("SELECT")),r=[].concat(t,n,a);for(let t=0;t<r.length;t++)e?r[t].setAttribute("disabled",""+e):r[t].removeAttribute("disabled")}function createNewSubgroup(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-group"),r="group"+1e4*(Date.now()+Math.random()),l=document.createElement("DIV");return l.setAttribute("style",""),l.className="sub-group form-group feedback-group",l.id=r+"-group",l.innerHTML=`<LABEL>Group Label : <INPUT ondragstart="return false" draggable="false" type="text" name="form-label" id="${r}-label"/> </LABEL> <br/>\n  <LABEL>Group ID : <INPUT ondragstart="return false" draggable="false" type="text" name="id" id="${r}-id"/></LABEL> <br/>\n  <LABEL>Group Descriptor : <br/> <TEXTAREA name="description" id="${r}-description"></TEXTAREA></LABEL> <br/>\n  <BUTTON id="${r}" onclick="FormLibrary.createNewResponseElement(this)" data-link-id="${r}" >Create New Respondable Below Last Respondable</BUTTON><br/>\n  <BUTTON  onclick="FormLibrary.createNewSubgroup(this)" data-link-id="${r}">Create New Group Below Last Respondable</BUTTON><br/>\n  <BUTTON  onclick="FormLibrary.deleteContainer('${n+"-group"}' , '${r+"-group"}')">Delete Subgroup</BUTTON><br/>\n  <SPAN  class="respondable-container" id="${r}-respondables"></SPAN><BR/>\n  `,l.ondrop=handleContainerDropWithinParent,l.ondragenter=handleContainerDragEnterWithinParent,l.ondragleave=handleContainerDragLeaveWithinParent,l.ondragover=function(e){e.preventDefault()},appendDropIcon(l),null==t&&(t=a.lastChild),a.insertBefore(l,t),r}function createNewResponseElement(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-respondables"),r="response"+1e4*(Date.now()+Math.random()),l=document.createElement("DIV");return l.setAttribute("style","width:400px;min-height:200px"),l.setAttribute("data-type","blank"),l.className="creation-prompt respondable-group feedback-group",l.id=r+"-fields",l.innerHTML=`<SPAN>Element Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        <LABEL>Item Type:</LABEL>\n        <SELECT id="${r}-type" ondragstart="return false" draggable="false">\n          <OPTION value="textarea">TextArea</OPTION>\n          <OPTION value="input">Input</OPTION>\n          <OPTION value="file">FileInput</OPTION>\n          <OPTION value="select">SelectGroup</OPTION>\n          <OPTION value="option">OptionGroup</OPTION>\n        </SELECT>\n      </LI>\n      <LI>If we want any of the unimplemented features, then you'll have to ask me or wait until I personally require it</LI>\n      <LI><BUTTON id="${r}" data-link-id="${r}" onclick="FormLibrary.responseTypeSelected('${n+"-respondables"}' , this)">Next</BUTTON></LI>\n      <LI><BUTTON  data-link-id="${r}" onclick="FormLibrary.deleteContainer('${n+"-respondables"}' , '${r+"-fields"}')">Delete</BUTTON></LI>\n    </UL>`,l.ondrop=handleContainerDropWithinParent,l.ondragenter=handleContainerDragEnterWithinParent,l.ondragleave=handleContainerDragLeaveWithinParent,l.ondragover=function(e){e.preventDefault()},appendDropIcon(l),null==t?a.appendChild(l):a.insertBefore(l,t),{base_id:r,respondable_container_id:n+"-respondables",container:l}}function appendDropIcon(e){let t=document.createElement("DIV");t.setAttribute("draggable","true"),t.className="drop-icon",t.textContent=" . . . ",t.title="Drag this to relocate item",t.ondragstart=handleContainerDragStartWithinParent,t.ondragend=handleContainerDragEndWithinParent,t.onmouseover=function(e){e.target.parentNode.style.borderColor="red"},t.onmouseleave=function(e){e.target.parentNode.style.borderColor=""},e.appendChild(t)}function handleContainerDragStartWithinParent(e){e.stopPropagation(),this==e.target?(t=this.parentNode,this.parentNode.style.opacity="0.4",this.style.cursor="grabbing",function createZIndexModifier(e){let t=document.createElement("STYLE");t.innerHTML=`#${e} > div * {\n    z-index: -1;\n    position: relative;\n  }`,t.id="zindexstyle",document.body.appendChild(t)}(t.parentNode.id)):e.preventDefault()}function handleContainerDragEndWithinParent(e){e.stopPropagation(),this==e.target?(t=void 0,this.parentNode.style.opacity="1.0",this.style.cursor=void 0,function removeZIndexModifier(){let e=document.getElementById("zindexstyle");document.body.removeChild(e)}()):e.preventDefault()}function handleContainerDragEnterWithinParent(e){e.currentTarget.parentNode==t.parentNode&&e.currentTarget!=t?e.currentTarget.style.border="1px dashed red":e.preventDefault()}function handleContainerDragLeaveWithinParent(e){e.currentTarget.parentNode==t.parentNode&&e.currentTarget!=t?e.currentTarget.style.border="":e.preventDefault()}function handleContainerDropWithinParent(e){if(e.currentTarget.parentNode!=t.parentNode||e.currentTarget==t)return void e.preventDefault();let n,a=this.parentNode;if(-1!=this.className.indexOf("sub-group")){let e={FormFields:[]};n=convertContainerToJSON(e,createFormStack([this,...Array.from(this.childNodes)],e)),recursiveFormFieldRebuild({SubGroups:e.FormFields},a.getElementsByTagName("BUTTON").item(0),t),setFieldsAsDisabled(!1)}else n=convertRespondable(this),rebuildRespondable(n,a.parentNode.getElementsByTagName("BUTTON").item(0),t);a.insertBefore(t,this),a.removeChild(this)}function convertContainerToJSON(e,t){for(;0!=t.length;){let n=0,a=t.shift(),r=a.paths,l=r[0],o=e.FormFields[l];for(let e=1;e<r.length;e++)o=o.SubGroups[r[e]];for(let e=0;e<a.node.childNodes.length;e++){let l=a.node.childNodes[e];if(l.nodeType==Node.ELEMENT_NODE)if(-1!=l.className.indexOf("form-group")){let e=r;e.push(n),t.unshift({node:l,paths:e}),n+=1,o.SubGroups.push({Label:"",ID:"",Description:"",SubGroups:[],Respondables:[]})}else if(-1!=l.className.indexOf("respondable-container")){if(0==l.childNodes.length)continue;let e=l.childNodes,t=e.length;for(let n=0;n<t;n++){let t=e[n];o.Respondables.push(convertRespondable(t))}}else if("LABEL"==l.tagName.toUpperCase()){let e=l.getElementsByTagName("INPUT")[0];if(e||(e=l.getElementsByTagName("SELECT")[0]),e||(e=l.getElementsByTagName("TEXTAREA")[0]),!e)continue;switch(e.getAttribute("name")){case"form-label":o.Label=e.value;break;case"id":o.ID=e.value;break;case"description":o.Description=e.value}}}}}function convertRespondable(e){let t={Object:{Field:{Name:"",Label:"",Required:!1}},Type:e.getAttribute("data-type")},n=e.getElementsByTagName("INPUT"),a=e.getElementsByTagName("SELECT");if(a.length)switch(e.getAttribute("data-type")){case"selectiongroup":case"genericinput":t.Object[a[0].getAttribute("Name")]=a[0].value}for(let a=0;a<n.length;a++){let r=n[a].value;if("checkbox"==n[a].getAttribute("type")?r=n[a].checked:"number"==n[a].getAttribute("type")&&(r=parseInt(n[a].value)),"selectiongroup"==e.getAttribute("data-type"))if(n[a].getAttribute("data-list-item-no")){t.Object.CheckableItems||(t.Object.CheckableItems=[]);let e=parseInt(n[a].getAttribute("data-list-item-no"));t.Object.CheckableItems[e]||(t.Object.CheckableItems[e]={}),t.Object.CheckableItems[e][n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}else if("optiongroup"==e.getAttribute("data-type"))if(n[a].getAttribute("data-list-item-no")){t.Object.Options||(t.Object.Options=[]);let e=parseInt(n[a].getAttribute("data-list-item-no"));t.Object.Options[e]||(t.Object.Options[e]={}),t.Object.Options[e][n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}}return t}function createFormStack(e,t){let n=[],a=0;for(let r=0;r<e.length;r++)e[r].nodeType==Node.ELEMENT_NODE&&-1!=e[r].className.indexOf("form-group")&&(n.unshift({node:e[r],paths:[a]}),a+=1,t.FormFields.push({Label:"",ID:"",Description:"",SubGroups:[],Respondables:[]}));return n}function createTextAreaInputs(e,t,n){let a="field"+e,r="text-area"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","textarea"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>TextArea Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${a}-placeholder"/><BR/>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createFileInputs(e,t,n){let a="field"+e,r="file"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","fileinput"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>FileInput Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Allowed Extention Pattern : <INPUT ondragstart="return false" draggable="false" data-field="1" name='AllowedExtRegex' id="${a}-allowed-ext"/><BR/>\n      </LI>\n      <LI>\n        Max Filesize(Bytes) : <INPUT ondragstart="return false" draggable="false" data-field="1" type="number" name='MaxSize' id="${a}-max-size"/><BR/>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createInputInputs(e,t,n){let a="field"+e,r="input"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","genericinput"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>Input Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${a}-placeholder"/><BR/>\n      </LI>\n      <LI>\n        Type : <SELECT data-field="1" name='Type' id="${a}-type" ondragstart="return false" draggable="false">\n          <OPTGROUP label="Text Types">\n            <OPTION value="text">Text</OPTION>\n            <OPTION value="email">Email</OPTION>\n            <OPTION value="number">Number</OPTION>\n            <OPTION value="password">Password</OPTION>\n            <OPTION value="url">URL</OPTION>\n          </OPTGROUP>\n          <OPTGROUP label="Time Types">\n            <OPTION value="time">Time</OPTION>\n            <OPTION value="date">Date</OPTION>\n          </OPTGROUP>\n          <OPTGROUP label="Special Types">\n            <OPTION value="color">Color Picker</OPTION>\n            <OPTION value="range">Number Range</OPTION>\n          </OPTGROUP>\n        </SELECT>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createSelectGroup(e,t,n){let a="field"+1e4*(Date.now()+Math.random()),r="select"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","selectiongroup"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>SelectGroup Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${a}-required"/><BR/>\n      </LI>\n      <LI>\n        Select type : <SELECT  data-field="1" name='SelectionCategory' type="checkbox" id="${a}-selectioncatergory" ondragstart="return false" draggable="false">\n          <OPTION value="checkbox">checkbox</OPTION>\n          <OPTION value="radio">radio</OPTION>\n        </SELECT><BR/>\n      </LI>\n      <LI>\n        Group Items :\n        <BUTTON id="${a}" onclick="FormLibrary.addListField(this , 'checkable')" data-link-id="${a}">+</BUTTON>\n        <BUTTON onclick="FormLibrary.removeListItem(this, 'checkable')" data-link-id="${a}">-</BUTTON><BR/>\n        <OL data-field="1" data-select="1" id="${a}-checkable">\n          <LI>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${a}-checkable-label-0"/>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${a}-checkable-value-0"/>\n          </LI>\n        </OL>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function addListField(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-"+t),r=document.createElement("LI"),l=a.children.length;r.innerHTML=`\n    <INPUT ondragstart="return false" draggable="false" placeholder="Label" data-list-item-no="${l}" data-field="1" name='Label' id="${n}-${t}-label-${l}"/>\n    <INPUT ondragstart="return false" draggable="false" placeholder="Value" data-list-item-no="${l}" data-field="1" name='Value' id="${n}-${t}-value-${l}"/>\n  `,a.appendChild(r)}function createOptionsGroup(e,t,n){let a="field"+1e4*(Date.now()+Math.random()),r="option"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","optiongroup"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>OptionGroup Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${a}-required"/><BR/>\n      </LI>\n      <LI>\n      Option Fields :\n        <BUTTON id="${a}" onclick="FormLibrary.addListField(this , 'option')" data-link-id="${a}">+</BUTTON>\n        <BUTTON onclick="FormLibrary.removeListItem(this, 'option')" data-link-id="${a}">-</BUTTON><BR/>\n        <OL data-field="1" data-select="1" id="${a}-option">\n          <LI>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${a}-option-label-0"/>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${a}-option-value-0"/>\n          </LI>\n        </OL>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function submitForm(e,t){let n={FormName:"",ID:"",Description:"",AnonOption:!1,FormFields:[]};if(document.querySelector("div[data-type=blank]"))return void(document.getElementById("response-container").innerHTML="Formatting: You have an unfinished respondable. Delete or complete it.");let a=document.getElementById("head-group").childNodes;for(let e=0;e<a.length;e++)if(a[e].nodeType==Node.ELEMENT_NODE&&"LABEL"==a[e].tagName.toUpperCase()){let t=a[e].getElementsByTagName("INPUT")[0];if(t||(t=a[e].getElementsByTagName("SELECT")[0]),t||(t=a[e].getElementsByTagName("TEXTAREA")[0]),!t)continue;switch(t.getAttribute("name")){case"form-name":n.FormName=t.value;break;case"id":n.ID=t.value;break;case"description":n.Description=t.value;break;case"anon-option":n.AnonOption=t.checked}}convertContainerToJSON(n,createFormStack(a,n)),sendCreateRequest(JSON.stringify(n),t)}function sendCreateRequest(e,t){var n=new FormData;n.append("json","1"),n.append("form-construct-json",e);var a=new XMLHttpRequest;a.open("POST",t,!0),a.onload=handleCreateComplete,a.send(n)}function handleCreateComplete(e){let t={};try{t=JSON.parse(e.target.responseText)}catch(e){return void(document.getElementById("response-container").innerHTML=`Hard server crash<BR/><TEXTAREA>${e.toString()}</TEXTAREA>`)}if(t.error)if(t["error-list"]){let e="There are some issues preventing you from submitting...<br/><ul>";t["error-list"].forEach((t=>{e+="<li>"+t.FailPosition+" : "+t.FailType+"</li>"})),e+="</ul>",document.getElementById("response-container").innerHTML=e}else document.getElementById("response-container").innerHTML="Server error: "+t.error;else document.getElementById("response-container").innerHTML=`${t.message} - ${Date.now()}<br/>\n      ${t.URL?`URL: <a href="${t.URL}">${t.URL}</a>`:""}\n    `}Object.defineProperty(n,"__esModule",{value:!0}),n.postDeleteResponse=n.postDeleteForm=n.handleCreateComplete=n.sendCreateRequest=n.submitForm=n.responseTypeSelected=n.createOptionsGroup=n.removeListItem=n.addListField=n.createSelectGroup=n.createInputInputs=n.createFileInputs=n.createTextAreaInputs=n.deleteContainer=n.createNewResponseElement=n.createNewSubgroup=n.setFieldsAsDisabled=n.rebuildFromRaw=n.attatchCreators=n.submitUserPost=n.dropdownList=n.helloWorld=void 0,n.helloWorld=function helloWorld(){console.log("hello client")},n.dropdownList=function dropdownList(e,t){let n=document.getElementById("container-"+t).getElementsByClassName("item-replies").item(0);if("▶"==e.firstChild.textContent){e.firstChild.textContent="( ...Loading... )";let a=new XMLHttpRequest;a.open("GET","/mod/api/form/"+t),a.onload=function(t){let a;e.firstChild.textContent="▼";try{a=JSON.parse(t.target.responseText)}catch(e){return void(n.innerHTML=`Hard server crash<BR/><TEXTAREA>${e.toString()}</TEXTAREA>`)}a.error?n.innerHTML=a.error:n.innerHTML=function buildDropdownResponse(e,t){return`<UL class="item-ul">${(()=>{if(!t.formatted_replies)return"<LI class='item-li'>Empty Set</LI>";let e="";return t.formatted_replies.forEach((t=>{e+=`<LI class='item-li'>${t}</LI>`})),e})()}</UL>`}(0,a)},a.onerror=function(){n.innerHTML="Server Issue"},a.send()}else"▼"==e.firstChild.textContent&&(e.firstChild.textContent="▶",n.innerHTML="");return!1},n.submitUserPost=function submitUserPost(e){let t=new FormData(e);t.append("json","1");let n=new XMLHttpRequest;return n.open("POST",""+window.location),n.onload=handleCreateComplete,n.send(t),!1},n.attatchCreators=function attatchCreators(e){let t=document.getElementById("sub-create");t.onclick=()=>createNewSubgroup(t);let n=document.getElementById("form-submit-button");n&&(n.onclick=()=>submitForm(n,"/mod/create"));let a=document.getElementById("form-edit-button");a&&(a.onclick=()=>submitForm(a,"/mod/edit/"+e.form_number))},n.rebuildFromRaw=function rebuildFromRaw(e){try{e=(e=(e=e.replace(/\n/g,"\\n")).replace(/\t/g,"\\t")).replace(/[\r\n\t\f\v]/g,""),recursiveFormFieldRebuild({SubGroups:JSON.parse(e).FormFields},document.getElementById("sub-create"))}catch(e){console.error(e),alert("Issue with rebuilding")}},n.setFieldsAsDisabled=setFieldsAsDisabled,n.createNewSubgroup=createNewSubgroup,n.createNewResponseElement=createNewResponseElement,n.deleteContainer=function deleteContainer(e,t){document.getElementById(e).removeChild(document.getElementById(t))},n.createTextAreaInputs=createTextAreaInputs,n.createFileInputs=createFileInputs,n.createInputInputs=createInputInputs,n.createSelectGroup=createSelectGroup,n.addListField=addListField,n.removeListItem=function removeListItem(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-"+t);a.childNodes.length<=1||a.removeChild(a.lastChild)},n.createOptionsGroup=createOptionsGroup,n.responseTypeSelected=function responseTypeSelected(e,t){let n=t.getAttribute("data-link-id"),a=document.getElementById(n+"-type"),r=document.getElementById(n+"-fields");switch(a.value){case"textarea":createTextAreaInputs(n,e,r);break;case"input":createInputInputs(n,e,r);break;case"file":createFileInputs(n,e,r);break;case"select":createSelectGroup(0,e,r);break;case"option":createOptionsGroup(0,e,r)}},n.submitForm=submitForm,n.sendCreateRequest=sendCreateRequest,n.handleCreateComplete=handleCreateComplete,n.postDeleteForm=function postDeleteForm(e,t){if(1==confirm("This will remove the form from the display and database.\n    Information files will retained on the server as a record until a duplicate named form is created")){let n=new XMLHttpRequest;n.open("POST",`/mod/form/delete/${e}/${t}`),n.onload=n=>{let a={};try{a=JSON.parse(n.target.responseText)}catch(e){return void alert("Error: Hard server crash")}if(a.error)alert("Error: "+a.error);else{alert("Success: "+a.message);let n=document.getElementById(`row-${e}-${t}`);n.parentNode.removeChild(n)}},n.send()}return!1},n.postDeleteResponse=function postDeleteResponse(e,t){if(1==confirm("This will remove the response and all information associated with it")){let n=new XMLHttpRequest;n.open("POST",`/mod/response/delete/${e}/${t}`),n.onload=n=>{let a={};try{a=JSON.parse(n.target.responseText)}catch(e){return void alert("Error: Hard server crash")}if(a.error)alert("Error: "+a.error);else{alert("Success: "+a.message);let n=document.getElementById(`row-${e}-${t}`);n.parentNode.removeChild(n)}},n.send()}return!1}})(),e})()));