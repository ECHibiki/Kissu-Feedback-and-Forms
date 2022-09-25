package builder

import (
  "encoding/json"
  "time"
  "os"
  "fmt"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
)


func ConstructFormObject(inputs map[string]string) former.FormConstruct {
  return former.FormConstruct{}
}

func MakeFormWritable(form former.FormConstruct) (types.FormDBFields  , error){
  marshaled_form , err := json.Marshal(form)
  if err != nil{
    return types.FormDBFields{}, err
  }
  return types.FormDBFields{
    ID: 0 ,
    Name: form.StorageName() ,
    FieldJSON: string(marshaled_form),
    UpdatedAt: time.Now().Unix(),
  } , nil
}

func CreateFormDirectory(form former.FormConstruct , cfg types.ConfigurationSettings) error{
  safe_name := form.StorageName()
  err := os.Mkdir(cfg.ResoruceDirectory + "/data/" + safe_name + "/" , 0775 )
  fmt.Println("Writting new form to " , cfg.ResoruceDirectory + "/data/" + safe_name + "/")
  if err != nil {
    return err
  }
  return nil
}
