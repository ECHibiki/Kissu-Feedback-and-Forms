package destroyer

import(
  "database/sql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "os"
  "fmt"
)

func DeleteForm(db *sql.DB , form_name string, form_num int64) error{
  _ , err := db.Exec("DELETE FROM forms WHERE name = ? and id = ?", form_name , form_num )
  // The form folder will linger
  if err != nil{
    return err
  }
  return err
}
func UndoFormDirectory(form former.FormConstruct , root_dir string) error{
  safe_name := form.StorageName()
  safe := tools.CheckSafeDirectoryName( safe_name )
  if !safe {
    panic("WTF THIS IS DANGEROUS PATH TO DELETE, " + safe_name)
  }
  err := os.RemoveAll(root_dir + "/data/" + safe_name + "/"  )
  return err
}
func UndoForm(db *sql.DB , form_name string, root_dir string) error{
  _ , err := db.Exec("DELETE FROM forms WHERE name = ?", form_name )
  if err != nil{
    return err
  }
  safe := tools.CheckSafeDirectoryName( form_name )
  if !safe {
    panic("WTF THIS IS DANGEROUS PATH TO DELETE, " + form_name)
  }
  err = os.RemoveAll(root_dir + "/data/" + form_name + "/")
  return err
}

// Delete a specific ID response with identifier
func DeleteResponse(db *sql.DB , root_dir string , reply_id int64 , form_name string , identifier string ) error{
  _ , err := db.Exec("DELETE FROM responses WHERE id = ? AND identifier = ? LIMIT 1", reply_id , identifier)
  if err != nil{
    return err
  }
  safe := tools.CheckSafeDirectoryName( identifier )
  if !safe {
    panic("WTF THIS IS DANGEROUS PATH TO DELETE, " + identifier)
  }
  safe = tools.CheckSafeDirectoryName( form_name )
  if !safe {
    panic("WTF THIS IS DANGEROUS PATH TO DELETE, " + form_name)
  }

  fmt.Printf("Deleted: %sdata/%s/%s\n" , root_dir , form_name , identifier)
  err = os.RemoveAll(root_dir + "/data/" + form_name + "/" + identifier)
  return err
}

// Delete a form's response with given identifier
func UndoResponse(db *sql.DB , form_response former.FormResponse, responder_name string , root_dir string ){
  deleteDatabaseResponse( db , form_response.RelationalID , responder_name )
  deleteResponderFolder( root_dir , form_response.FormName , responder_name )
}

func deleteResponderFolder(root_dir string,  form_name string , old_user_name string) error{
  return os.RemoveAll(root_dir + "/data/" + form_name + "/" + old_user_name)
}

func deleteDatabaseResponse(db *sql.DB , relational_id int64 , old_user_name string) error{
    _ , err := db.Exec("DELETE FROM responses WHERE fk_id= ? AND identifier = ? LIMIT 1" , relational_id , old_user_name)
    return err
}
