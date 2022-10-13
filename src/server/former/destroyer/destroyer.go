package destroyer

import(
  "database/sql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "os"
)

func DeleteForm(db *sql.DB , form_name string, form_num int64) error{
  _ , err := db.Exec("DELETE FROM forms WHERE name = ? and id = ?", form_name , form_num )
  // The form folder will linger
  if err != nil{
    return err
  }
  return err
}
func UndoForm(db *sql.DB , form_name string, root_dir string) error{
  _ , err := db.Exec("DELETE FROM forms WHERE name = ?", form_name )
  if err != nil{
    return err
  }
  err = os.RemoveAll(root_dir + "/data/" + form_name + "/")
  return err
}

func DeleteResponse(db *sql.DB , initialization_folder string , reply_id int64 , form_name string , identifier string ) error{
  _ , err := db.Exec("DELETE FROM responses WHERE id=? AND identifier=?", reply_id , identifier)
  if err != nil{
    return err
  }
  err = os.RemoveAll(initialization_folder + "/data/" + form_name + "/" + identifier)
  return err
}

func UndoDirectory(form former.FormConstruct , root_dir string) error{
  safe_name := form.StorageName()
  err := os.RemoveAll(root_dir + "/data/" + safe_name + "/"  )
  return err
}
