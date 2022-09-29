package tools

import (
  "os"
)

func LogError(storage_dir string, message string){
  err_handler , err := os.OpenFile(storage_dir +  "errors.log" , os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    panic(err)
  }
  defer err_handler.Close()
  _ , err = err_handler.WriteString("File write fail: " + message + "\n")
  if err != nil {
    panic(err)
  }
}

func WriteResponsesToJSONFile(root_dir string , resp former.FormResponse) error {
  storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"

  json_resp := ConvertFormResponseToJSONFormResponse(root_dir , resp)

  json_bytes , err := json.MarshalIndent(json_resp , "" , " ")
  if err != nil {
    return err
  }
  err = ioutil.WriteFile(storage_dir + "responses.json" , json_bytes , 0644)
  return err
}

func ConvertFormResponseToJSONFormResponse(root_dir string, resp former.FormResponse) former.JSONFormResponse {
  json_resp := former.JSONFormResponse{}
  json_resp.FormName= resp.FormName
  json_resp.RelationalID = resp.RelationalID
  json_resp.ResponderID = resp.ResponderID
  json_resp.Responses = resp.Responses
  json_resp.FilePaths = make(map[string]string)
  storage_dir := root_dir + "/data/" + resp.FormName + "/" + resp.ResponderID + "/"
  for k, v := range(resp.FileObjects) {
    json_resp.FilePaths[k] = storage_dir + "files/" + v.Header.Filename
  }

  return json_resp
}

func WriteFilesFromMultipart(root_dir string , response_struct former.FormResponse) []error{
  storage_dir := root_dir + "/data/" + response_struct.FormName + "/" + response_struct.ResponderID + "/files/"
  var err_list []error = []error{}
  for field_name , file_object := range response_struct.FileObjects {
    fname := field_name + "-" + file_object.Header.Filename
    if strings.Contains(fname , "/"){
      tools.LogError( storage_dir , storage_dir +  fname )
      err_list = append(err_list , errors.New("File " + fname +  " contained illegal characters"))
      continue
    }
    handler, err := os.OpenFile(storage_dir +  fname , os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
      tools.LogError( storage_dir , storage_dir +  fname )
      err_list = append(err_list , err)
      continue
    }
    defer handler.Close()
    _ , err = io.Copy(handler , file_object.File )
     if err != nil {
       tools.LogError( storage_dir , storage_dir +  fname )
       err_list = append(err_list , err)
       continue
     }
  }
  return err_list
}
