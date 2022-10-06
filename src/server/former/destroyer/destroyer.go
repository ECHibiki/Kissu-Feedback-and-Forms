package destroyer

import(
  "database/sql"
  "os"
)

func DeleteForm(db *sql.DB , form_name string) error{
  _ , err := db.Exec("DELETE FROM forms WHERE name = ?", form_name )
  // The form folder will linger
  if err != nil{
    return err
  }
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
