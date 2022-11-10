!function webpackUniversalModuleDefinition(e,t){"object"==typeof exports&&"object"==typeof module?module.exports=t():"function"==typeof define&&define.amd?define("FormLibrary",[],t):"object"==typeof exports?exports.FormLibrary=t():e.FormLibrary=t()}(self,(()=>(()=>{"use strict";var e={614:(e,t,n)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.createFormStack=t.convertRespondable=t.convertContainerToJSON=t.rebuildRespondable=t.recursiveFormFieldRebuild=t.buildDropdownResponse=t.createNewResponseElement=t.setFieldsAsDisabled=t.removeListItem=t.addListField=t.createNewSubgroup=t.responseTypeSelected=t.deleteContainer=void 0;const a=n(525);function createNewSubgroup(e,t){let n=e.getAttribute("data-link-id"),r=document.getElementById(n+"-group"),o="group"+1e4*(Date.now()+Math.random()),l=document.createElement("DIV");return l.setAttribute("style",""),l.className="sub-group form-group feedback-group",l.id=o+"-group",l.innerHTML=`<LABEL>Group Label : <INPUT ondragstart="return false" draggable="false" type="text" name="form-label" id="${o}-label"/> </LABEL> <br/>\n  <LABEL>Group ID : <INPUT ondragstart="return false" draggable="false" type="text" name="id" id="${o}-id"/></LABEL> <br/>\n  <LABEL>Group Descriptor : <br/> <TEXTAREA name="description" id="${o}-description"></TEXTAREA></LABEL> <br/>\n  <BUTTON id="${o}" onclick="FormLibrary.createNewResponseElement(this)" data-link-id="${o}" >Create New Respondable Below Last Respondable</BUTTON><br/>\n  <BUTTON  onclick="FormLibrary.createNewSubgroup(this)" data-link-id="${o}">Create New Group Below Last Respondable</BUTTON><br/>\n  <BUTTON  onclick="FormLibrary.deleteContainer('${n+"-group"}' , '${o+"-group"}')">Delete Subgroup</BUTTON><br/>\n  <SPAN  class="respondable-container" id="${o}-respondables"></SPAN><BR/>\n  `,l.ondrop=a.handleContainerDropWithinParent,l.ondragenter=a.handleContainerDragEnterWithinParent,l.ondragleave=a.handleContainerDragLeaveWithinParent,l.ondragover=function(e){e.preventDefault()},appendDropIcon(l),null==t&&(t=r.lastChild),r.insertBefore(l,t),o}function addListField(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-"+t),r=document.createElement("LI"),o=a.children.length;r.innerHTML=`\n    <INPUT ondragstart="return false" draggable="false" placeholder="Label" data-list-item-no="${o}" data-field="1" name='Label' id="${n}-${t}-label-${o}"/>\n    <INPUT ondragstart="return false" draggable="false" placeholder="Value" data-list-item-no="${o}" data-field="1" name='Value' id="${n}-${t}-value-${o}"/>\n  `,a.appendChild(r)}function setFieldsAsDisabled(e){let t=Array.from(document.getElementsByTagName("BUTTON")),n=Array.from(document.getElementsByTagName("INPUT")),a=Array.from(document.getElementsByTagName("SELECT")),r=[].concat(t,n,a);for(let t=0;t<r.length;t++)e?r[t].setAttribute("disabled",""+e):r[t].removeAttribute("disabled")}function createNewResponseElement(e,t){let n=e.getAttribute("data-link-id"),r=document.getElementById(n+"-respondables"),o="response"+1e4*(Date.now()+Math.random()),l=document.createElement("DIV");return l.setAttribute("style","width:400px;min-height:200px"),l.setAttribute("data-type","blank"),l.className="creation-prompt respondable-group feedback-group",l.id=o+"-fields",l.innerHTML=`<SPAN>Element Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        <LABEL>Item Type:</LABEL>\n        <SELECT id="${o}-type" ondragstart="return false" draggable="false">\n          <OPTION value="textarea">TextArea</OPTION>\n          <OPTION value="input">Input</OPTION>\n          <OPTION value="file">FileInput</OPTION>\n          <OPTION value="select">SelectGroup</OPTION>\n          <OPTION value="option">OptionGroup</OPTION>\n        </SELECT>\n      </LI>\n      <LI>If we want any of the unimplemented features, then you'll have to ask me or wait until I personally require it</LI>\n      <LI><BUTTON id="${o}" data-link-id="${o}" onclick="FormLibrary.responseTypeSelected('${n+"-respondables"}' , this)">Next</BUTTON></LI>\n      <LI><BUTTON  data-link-id="${o}" onclick="FormLibrary.deleteContainer('${n+"-respondables"}' , '${o+"-fields"}')">Delete</BUTTON></LI>\n    </UL>`,l.ondrop=a.handleContainerDropWithinParent,l.ondragenter=a.handleContainerDragEnterWithinParent,l.ondragleave=a.handleContainerDragLeaveWithinParent,l.ondragover=function(e){e.preventDefault()},appendDropIcon(l),null==t?r.appendChild(l):r.insertBefore(l,t),{base_id:o,respondable_container_id:n+"-respondables",container:l}}function rebuildRespondable(e,t,n){let a=createNewResponseElement(t,n),r="";switch(e.Type){case"textarea":r=createTextAreaInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-placeholder").value=e.Object.Placeholder;break;case"genericinput":r=createInputInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-placeholder").value=e.Object.Placeholder,document.getElementById(r+"-type").value=e.Object.Type;break;case"fileinput":r=createFileInputs(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required,document.getElementById(r+"-allowed-ext").value=e.Object.AllowedExtRegex,document.getElementById(r+"-max-size").value=e.Object.MaxSize;break;case"selectiongroup":r=createSelectGroup(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required;for(let t=1;t<e.Object.CheckableItems.length;t++)addListField(document.getElementById(r),"checkable");e.Object.CheckableItems.forEach(((e,t)=>{document.getElementById(r+"-checkable-label-"+t).value=e.Label,document.getElementById(r+"-checkable-value-"+t).value=e.Value}));break;case"optiongroup":r=createOptionsGroup(a.base_id,a.respondable_container_id,a.container),document.getElementById(r+"-name").value=e.Object.Field.Name,document.getElementById(r+"-label").value=e.Object.Field.Label,document.getElementById(r+"-required").checked=e.Object.Field.Required;for(let t=1;t<e.Object.Options.length;t++)addListField(document.getElementById(r),"option");e.Object.Options.forEach(((e,t)=>{document.getElementById(r+"-option-label-"+t).value=e.Label,document.getElementById(r+"-option-value-"+t).value=e.Value}))}}function appendDropIcon(e){let t=document.createElement("DIV");t.setAttribute("draggable","true"),t.className="drop-icon",t.textContent=" . . . ",t.title="Drag this to relocate item",t.ondragstart=a.handleContainerDragStartWithinParent,t.ondragend=a.handleContainerDragEndWithinParent,t.onmouseover=function(e){e.target.parentNode.style.borderColor="red"},t.onmouseleave=function(e){e.target.parentNode.style.borderColor=""},e.appendChild(t)}function convertRespondable(e){let t={Object:{Field:{Name:"",Label:"",Required:!1}},Type:e.getAttribute("data-type")},n=e.getElementsByTagName("INPUT"),a=e.getElementsByTagName("SELECT");if(a.length)switch(e.getAttribute("data-type")){case"selectiongroup":case"genericinput":t.Object[a[0].getAttribute("Name")]=a[0].value}for(let a=0;a<n.length;a++){let r=n[a].value;if("checkbox"==n[a].getAttribute("type")?r=n[a].checked:"number"==n[a].getAttribute("type")&&(r=parseInt(n[a].value)),"selectiongroup"==e.getAttribute("data-type"))if(n[a].getAttribute("data-list-item-no")){t.Object.CheckableItems||(t.Object.CheckableItems=[]);let e=parseInt(n[a].getAttribute("data-list-item-no"));t.Object.CheckableItems[e]||(t.Object.CheckableItems[e]={}),t.Object.CheckableItems[e][n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}else if("optiongroup"==e.getAttribute("data-type"))if(n[a].getAttribute("data-list-item-no")){t.Object.Options||(t.Object.Options=[]);let e=parseInt(n[a].getAttribute("data-list-item-no"));t.Object.Options[e]||(t.Object.Options[e]={}),t.Object.Options[e][n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}else switch(n[a].getAttribute("name")){case"Name":case"Label":case"Required":t.Object.Field[n[a].getAttribute("name")]=r;break;default:t.Object[n[a].getAttribute("name")]=r}}return t}function createTextAreaInputs(e,t,n){let a="field"+e,r="text-area"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","textarea"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>TextArea Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${a}-placeholder"/><BR/>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createFileInputs(e,t,n){let a="field"+e,r="file"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","fileinput"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>FileInput Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Allowed Extention Pattern : <INPUT ondragstart="return false" draggable="false" data-field="1" name='AllowedExtRegex' id="${a}-allowed-ext"/><BR/>\n      </LI>\n      <LI>\n        Max Filesize(Bytes) : <INPUT ondragstart="return false" draggable="false" data-field="1" type="number" name='MaxSize' id="${a}-max-size"/><BR/>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createInputInputs(e,t,n){let a="field"+e,r="input"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","genericinput"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>Input Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false"  data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' id="${a}-required" type="checkbox"/><BR/>\n      </LI>\n      <LI>\n        Placeholder : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Placeholder' id="${a}-placeholder"/><BR/>\n      </LI>\n      <LI>\n        Type : <SELECT data-field="1" name='Type' id="${a}-type" ondragstart="return false" draggable="false">\n          <OPTGROUP label="Text Types">\n            <OPTION value="text">Text</OPTION>\n            <OPTION value="email">Email</OPTION>\n            <OPTION value="number">Number</OPTION>\n            <OPTION value="password">Password</OPTION>\n            <OPTION value="url">URL</OPTION>\n          </OPTGROUP>\n          <OPTGROUP label="Time Types">\n            <OPTION value="time">Time</OPTION>\n            <OPTION value="date">Date</OPTION>\n          </OPTGROUP>\n          <OPTGROUP label="Special Types">\n            <OPTION value="color">Color Picker</OPTION>\n            <OPTION value="range">Number Range</OPTION>\n          </OPTGROUP>\n        </SELECT>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createSelectGroup(e,t,n){let a="field"+1e4*(Date.now()+Math.random()),r="select"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","selectiongroup"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>SelectGroup Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${a}-required"/><BR/>\n      </LI>\n      <LI>\n        Select type : <SELECT  data-field="1" name='SelectionCategory' type="checkbox" id="${a}-selectioncatergory" ondragstart="return false" draggable="false">\n          <OPTION value="checkbox">checkbox</OPTION>\n          <OPTION value="radio">radio</OPTION>\n        </SELECT><BR/>\n      </LI>\n      <LI>\n        Group Items :\n        <BUTTON id="${a}" onclick="FormLibrary.addListField(this , 'checkable')" data-link-id="${a}">+</BUTTON>\n        <BUTTON onclick="FormLibrary.removeListItem(this, 'checkable')" data-link-id="${a}">-</BUTTON><BR/>\n        <OL data-field="1" data-select="1" id="${a}-checkable">\n          <LI>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${a}-checkable-label-0"/>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${a}-checkable-value-0"/>\n          </LI>\n        </OL>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}function createOptionsGroup(e,t,n){let a="field"+1e4*(Date.now()+Math.random()),r="option"+1e4*(Date.now()+Math.random());return n.className="respondable-group feedback-group",n.id=r,n.setAttribute("data-type","optiongroup"),n.setAttribute("style","width:400px;min-height:200px"),n.innerHTML=`<SPAN>OptionGroup Creation Info:</SPAN><BR/>\n    <UL>\n      <LI>\n        Name : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Name' id="${a}-name"/><BR/>\n      </LI>\n      <LI>\n        Label : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Label' id="${a}-label"/><BR/>\n      </LI>\n      <LI>\n        Required : <INPUT ondragstart="return false" draggable="false" data-field="1" name='Required' type="checkbox" id="${a}-required"/><BR/>\n      </LI>\n      <LI>\n      Option Fields :\n        <BUTTON id="${a}" onclick="FormLibrary.addListField(this , 'option')" data-link-id="${a}">+</BUTTON>\n        <BUTTON onclick="FormLibrary.removeListItem(this, 'option')" data-link-id="${a}">-</BUTTON><BR/>\n        <OL data-field="1" data-select="1" id="${a}-option">\n          <LI>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Label"  data-list-item-no="0" data-field="1" name='Label' id="${a}-option-label-0"/>\n            <INPUT ondragstart="return false" draggable="false" placeholder="Value"  data-list-item-no="0" data-field="1" name='Value'  id="${a}-option-value-0"/>\n          </LI>\n        </OL>\n      </LI>\n      <LI><BUTTON id="${r}" onclick="FormLibrary.deleteContainer('${t}' , '${r}')">Delete</BUTTON></LI>\n    </UL>`,appendDropIcon(n),a}t.deleteContainer=function deleteContainer(e,t){document.getElementById(e).removeChild(document.getElementById(t))},t.responseTypeSelected=function responseTypeSelected(e,t){let n=t.getAttribute("data-link-id"),a=document.getElementById(n+"-type"),r=document.getElementById(n+"-fields");switch(a.value){case"textarea":createTextAreaInputs(n,e,r);break;case"input":createInputInputs(n,e,r);break;case"file":createFileInputs(n,e,r);break;case"select":createSelectGroup(n,e,r);break;case"option":createOptionsGroup(n,e,r)}},t.createNewSubgroup=createNewSubgroup,t.addListField=addListField,t.removeListItem=function removeListItem(e,t){let n=e.getAttribute("data-link-id"),a=document.getElementById(n+"-"+t);a.childNodes.length<=1||a.removeChild(a.lastChild)},t.setFieldsAsDisabled=setFieldsAsDisabled,t.createNewResponseElement=createNewResponseElement,t.buildDropdownResponse=function buildDropdownResponse(e,t){return`<UL class="item-ul">${(()=>{if(!t.formatted_replies)return"<LI class='item-li'>Empty Set</LI>";let e="";return t.formatted_replies.forEach((t=>{e+=`<LI class='item-li' id='row-${t.Name}-${t.ID}'>${t.Body}</LI>`})),e})()}</UL>`},t.recursiveFormFieldRebuild=function recursiveFormFieldRebuild(e,t,n){if(!e)return;let a=e.Respondables;a&&a.forEach((e=>{rebuildRespondable(e,t)})),setFieldsAsDisabled(!0);let r=e.SubGroups;r&&r.forEach((e=>{let a=createNewSubgroup(t,n);document.getElementById(a+"-label").value=e.Label,document.getElementById(a+"-id").value=e.ID,document.getElementById(a+"-description").value=e.Description,recursiveFormFieldRebuild(e,document.getElementById(a))}))},t.rebuildRespondable=rebuildRespondable,t.convertContainerToJSON=function convertContainerToJSON(e,t){for(;0!=t.length;){let n=0,a=t.shift(),r=a.paths,o=r[0],l=e.FormFields[o],i=0;for(let e=1;e<r.length;e++)i++,l=l.SubGroups[r[e].pop()];for(let e=0;e<a.node.childNodes.length;e++){let o=a.node.childNodes[e];if(o.nodeType==Node.ELEMENT_NODE)if(-1!=o.className.indexOf("form-group")){let e=r;!e[i+1]&&e.push([]),e[i+1].push(n),n+=1,t.unshift({node:o,paths:e}),l.SubGroups.push({Label:"",ID:"",Description:"",SubGroups:[],Respondables:[]})}else if(-1!=o.className.indexOf("respondable-container")){if(0==o.childNodes.length)continue;let e=o.childNodes,t=e.length;for(let n=0;n<t;n++){let t=e[n];l.Respondables.push(convertRespondable(t))}}else if("LABEL"==o.tagName.toUpperCase()){let e=o.getElementsByTagName("INPUT")[0];if(e||(e=o.getElementsByTagName("SELECT")[0]),e||(e=o.getElementsByTagName("TEXTAREA")[0]),!e)continue;switch(e.getAttribute("name")){case"form-label":l.Label=e.value;break;case"id":l.ID=e.value;break;case"description":l.Description=e.value}}}}},t.convertRespondable=convertRespondable,t.createFormStack=function createFormStack(e,t){let n=[],a=0;for(let r=0;r<e.length;r++)e[r].nodeType==Node.ELEMENT_NODE&&-1!=e[r].className.indexOf("form-group")&&(n.unshift({node:e[r],paths:[a]}),a+=1,t.FormFields.push({Label:"",ID:"",Description:"",SubGroups:[],Respondables:[]}));return n}},525:(e,t,n)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.handleContainerDropWithinParent=t.handleContainerDragLeaveWithinParent=t.handleContainerDragEnterWithinParent=t.handleContainerDragEndWithinParent=t.handleContainerDragStartWithinParent=void 0;const a=n(614);var r;t.handleContainerDragStartWithinParent=function handleContainerDragStartWithinParent(e){e.stopPropagation(),this==e.target?(r=this.parentNode,this.parentNode.style.opacity="0.4",this.style.cursor="grabbing",function createZIndexModifier(e){let t=document.createElement("STYLE");t.innerHTML=`#${e} > div * {\n    z-index: -1;\n    position: relative;\n  }`,t.id="zindexstyle",document.body.appendChild(t)}(r.parentNode.id)):e.preventDefault()},t.handleContainerDragEndWithinParent=function handleContainerDragEndWithinParent(e){e.stopPropagation(),e.preventDefault(),r=void 0,this.parentNode.style.opacity="1.0",this.style.cursor="",function removeZIndexModifier(){let e=document.getElementById("zindexstyle");document.body.removeChild(e)}()},t.handleContainerDragEnterWithinParent=function handleContainerDragEnterWithinParent(e){e.currentTarget.parentNode==r.parentNode&&e.currentTarget!=r?e.currentTarget.style.border="1px dashed red":e.preventDefault()},t.handleContainerDragLeaveWithinParent=function handleContainerDragLeaveWithinParent(e){e.currentTarget.parentNode==r.parentNode&&e.currentTarget!=r?e.currentTarget.style.border="":e.preventDefault()},t.handleContainerDropWithinParent=function handleContainerDropWithinParent(e){if(e.currentTarget.parentNode!=r.parentNode||e.currentTarget==r)return void e.preventDefault();let t,n=this.parentNode;if(-1!=this.className.indexOf("sub-group")){let e=[this,...Array.from(this.childNodes)],o={FormFields:[]},l=(0,a.createFormStack)(e,o);t=(0,a.convertContainerToJSON)(o,l);let i={SubGroups:o.FormFields};(0,a.recursiveFormFieldRebuild)(i,n.getElementsByTagName("BUTTON").item(0),r),(0,a.setFieldsAsDisabled)(!1)}else t=(0,a.convertRespondable)(this),(0,a.rebuildRespondable)(t,n.parentNode.getElementsByTagName("BUTTON").item(0),r);n.insertBefore(r,this),n.removeChild(this)}},802:function(e,t,n){var a=this&&this.__createBinding||(Object.create?function(e,t,n,a){void 0===a&&(a=n);var r=Object.getOwnPropertyDescriptor(t,n);r&&!("get"in r?!t.__esModule:r.writable||r.configurable)||(r={enumerable:!0,get:function(){return t[n]}}),Object.defineProperty(e,a,r)}:function(e,t,n,a){void 0===a&&(a=n),e[a]=t[n]}),r=this&&this.__setModuleDefault||(Object.create?function(e,t){Object.defineProperty(e,"default",{enumerable:!0,value:t})}:function(e,t){e.default=t}),o=this&&this.__importStar||function(e){if(e&&e.__esModule)return e;var t={};if(null!=e)for(var n in e)"default"!==n&&Object.prototype.hasOwnProperty.call(e,n)&&a(t,e,n);return r(t,e),t};Object.defineProperty(t,"__esModule",{value:!0}),t.submitUserPost=t.postDeleteResponse=t.postDeleteForm=t.setFieldsAsDisabled=t.rebuildFromRaw=t.removeListItem=t.addListField=t.responseTypeSelected=t.deleteContainer=t.createNewSubgroup=t.createNewResponseElement=t.attatchCreators=t.dropdownList=void 0,console.log("FormLibrary initialized.\nFeedback&Forms product of Kissu.moe");const l=o(n(290)),i=o(n(970)),d=o(n(614));t.dropdownList=function dropdownList(e,t){return l.dropdownList(e,t)},t.attatchCreators=function attatchCreators(e){i.attatchCreators(e)},t.createNewResponseElement=function createNewResponseElement(e,t){return d.createNewResponseElement(e,t)},t.createNewSubgroup=function createNewSubgroup(e,t){return d.createNewSubgroup(e,t)},t.deleteContainer=function deleteContainer(e,t){d.deleteContainer(e,t)},t.responseTypeSelected=function responseTypeSelected(e,t){d.responseTypeSelected(e,t)},t.addListField=function addListField(e,t){d.addListField(e,t)},t.removeListItem=function removeListItem(e,t){d.removeListItem(e,t)},t.rebuildFromRaw=function rebuildFromRaw(e){i.rebuildFromRaw(e)},t.setFieldsAsDisabled=function setFieldsAsDisabled(e){d.setFieldsAsDisabled(e)},t.postDeleteForm=function postDeleteForm(e,t){return l.postDeleteForm(e,t)},t.postDeleteResponse=function postDeleteResponse(e,t){return l.postDeleteResponse(e,t)},t.submitUserPost=function submitUserPost(e){return l.submitUserPost(e)}},970:(e,t,n)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.rebuildFromRaw=t.attatchCreators=t.helloWorld=void 0;const a=n(614),r=n(290);t.helloWorld=function helloWorld(){console.log("hello client")},t.attatchCreators=function attatchCreators(e){let t=document.getElementById("sub-create");t.onclick=()=>(0,a.createNewSubgroup)(t);let n=document.getElementById("form-submit-button");n&&(n.onclick=()=>(0,r.submitForm)(n,"/mod/create"));let o=document.getElementById("form-edit-button");o&&(o.onclick=()=>(0,r.submitForm)(o,"/mod/edit/"+e.form_number))},t.rebuildFromRaw=function rebuildFromRaw(e){try{e=(e=(e=e.replace(/\n/g,"\\n")).replace(/\t/g,"\\t")).replace(/[\r\n\t\f\v]/g,"");let t={SubGroups:JSON.parse(e).FormFields},n=document.getElementById("sub-create");(0,a.recursiveFormFieldRebuild)(t,n)}catch(e){console.error(e),alert("Issue with rebuilding")}}},290:(e,t,n)=>{Object.defineProperty(t,"__esModule",{value:!0}),t.submitForm=t.dropdownList=t.postDeleteForm=t.postDeleteResponse=t.submitUserPost=void 0;const a=n(614);function handleCreateComplete(e){let t={};try{t=JSON.parse(e.target.responseText)}catch(e){return void(document.getElementById("response-container").innerHTML=`Hard server crash<BR/><TEXTAREA>${e.toString()}</TEXTAREA>`)}if(t.error)if(t["error-list"]){let e="There are some issues preventing you from submitting...<br/><ul>";t["error-list"].forEach((t=>{e+="<li>"+t.FailPosition+" : "+t.FailType+"</li>"})),e+="</ul>",document.getElementById("response-container").innerHTML=e}else document.getElementById("response-container").innerHTML="Server error: "+t.error;else document.getElementById("response-container").innerHTML=`${t.message} - ${Date.now()}<br/>\n      ${t.URL?`URL: <a href="${t.URL}">${t.URL}</a>`:""}\n    `}t.submitUserPost=function submitUserPost(e){let t=new FormData(e);t.append("json","1");let n=new XMLHttpRequest;return n.open("POST",""+window.location),n.onload=handleCreateComplete,n.send(t),!1},t.postDeleteResponse=function postDeleteResponse(e,t){if(1==confirm("This will remove the response and all information associated with it")){let n=new XMLHttpRequest;n.open("POST",`/mod/response/delete/${e}/${t}`),n.onload=n=>{let a={};try{a=JSON.parse(n.target.responseText)}catch(e){return void alert("Error: Hard server crash")}if(a.error)alert("Error: "+a.error);else{alert("Success: "+a.message);let n=document.getElementById(`row-${e}-${t}`);n.parentNode.removeChild(n)}},n.send()}return!1},t.postDeleteForm=function postDeleteForm(e,t){if(1==confirm("This will remove the form from the display and database.\n    Information files will retained on the server as a record until a duplicate named form is created")){let n=new XMLHttpRequest;n.open("POST",`/mod/form/delete/${e}/${t}`),n.onload=n=>{let a={};try{a=JSON.parse(n.target.responseText)}catch(e){return void alert("Error: Hard server crash")}if(a.error)alert("Error: "+a.error);else{alert("Success: "+a.message);let n=document.getElementById(`row-${e}-${t}`);n.parentNode.removeChild(n)}},n.send()}return!1},t.dropdownList=function dropdownList(e,t){let n=document.getElementById("container-"+t).getElementsByClassName("item-replies").item(0);if("▶"==e.firstChild.textContent){e.firstChild.textContent="( ...Loading... )";let r=new XMLHttpRequest;r.open("GET","/mod/api/form/"+t),r.onload=function(t){let r;e.firstChild.textContent="▼";try{r=JSON.parse(t.target.responseText)}catch(e){return void(n.innerHTML=`Hard server crash<BR/><TEXTAREA>${e.toString()}</TEXTAREA>`)}r.error?n.innerHTML=r.error:n.innerHTML=(0,a.buildDropdownResponse)(n,r)},r.onerror=function(){n.innerHTML="Server Issue"},r.send()}else"▼"==e.firstChild.textContent&&(e.firstChild.textContent="▶",n.innerHTML="");return!1},t.submitForm=function submitForm(e,t){let n={FormName:"",ID:"",Description:"",AnonOption:!1,FormFields:[]};if(document.querySelector("div[data-type=blank]"))return void(document.getElementById("response-container").innerHTML="Formatting: You have an unfinished respondable. Delete or complete it.");let r=document.getElementById("head-group").childNodes;for(let e=0;e<r.length;e++)if(r[e].nodeType==Node.ELEMENT_NODE&&"LABEL"==r[e].tagName.toUpperCase()){let t=r[e].getElementsByTagName("INPUT")[0];if(t||(t=r[e].getElementsByTagName("SELECT")[0]),t||(t=r[e].getElementsByTagName("TEXTAREA")[0]),!t)continue;switch(t.getAttribute("name")){case"form-name":n.FormName=t.value;break;case"id":n.ID=t.value;break;case"description":n.Description=t.value;break;case"anon-option":n.AnonOption=t.checked}}let o=(0,a.createFormStack)(r,n);(0,a.convertContainerToJSON)(n,o),function sendCreateRequest(e,t){var n=new FormData;n.append("json","1"),n.append("form-construct-json",e);var a=new XMLHttpRequest;a.open("POST",t,!0),a.onload=handleCreateComplete,a.send(n)}(JSON.stringify(n),t)}}},t={};var n=function __webpack_require__(n){var a=t[n];if(void 0!==a)return a.exports;var r=t[n]={exports:{}};return e[n].call(r.exports,r,r.exports,__webpack_require__),r.exports}(802);return n})()));