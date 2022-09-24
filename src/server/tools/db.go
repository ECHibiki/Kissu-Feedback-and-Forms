package tools

import(
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
)

func BuildDBTables( db *sql.DB ){
  _, err := db.Exec("CREATE TABLE forms(id INT AUTO_INCREMENT PRIMARY KEY , field_json TEXT , updated_at INT)")
  if err != nil{
    panic(err)
  }
  _, err = db.Exec("CREATE TABLE responses(id INT AUTO_INCREMENT PRIMARY KEY , fk_id INT, identifier VARCHAR(60) , response_json TEXT , submitted_at INT  , FOREIGN KEY (fk_id) REFERENCES forms(id) )")
  if err != nil{
    panic(err)
  }
  _, err = db.Exec("CREATE TABLE passwords(hashed_password VARCHAR(255) , hash_system VARCHAR(255) , hash_scrambler VARCHAR(255))")
  if err != nil{
    panic(err)
  }
  _, err = db.Exec("CREATE TABLE logins(time_at INT , cookie VARCHAR(255) , ip VARCHAR(255))")
  if err != nil{
    panic(err)
  }
}

func QuickDBConnect(cfg types.ConfigurationSettings) *sql.DB{
  db , err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
    cfg.DBUserName,
    cfg.DBCredentials,
    cfg.DBAddr,
    cfg.DBName,
    ),
  )
  if err != nil {
    panic("DB connect error")
  }
  return db
}

func WritePassToDB(db *sql.DB , pass types.PasswordsDBFields) error{
  _ , err := db.Exec("INSERT INTO passwords VALUES( ? , ? , ? )" , pass.HashedPassword , pass.HashSystem , pass.HashScrambler)
  return err
}

func FetchPassword(db *sql.DB) (types.PasswordsDBFields , error){
  q := db.QueryRow("SELECT * FROM passwords");
  var pass types.PasswordsDBFields
  err := q.Scan(&pass.HashedPassword , &pass.HashSystem , &pass.HashScrambler)
  return pass , err
}

func StoreFormToDB(db *sql.DB , form types.FormDBFields ) error {
  _ , err := db.Exec("INSERT INTO forms VALUES( NULL , ? , ? )" , form.FieldJSON , form.UpdatedAt)
  return err
}
func GetFormOfID(db *sql.DB , id int64 ) (types.FormDBFields , error) {
  q := db.QueryRow("SELECT * FROM forms WHERE id=?" , id );
  var form types.FormDBFields
  err := q.Scan(&form.ID , &form.FieldJSON , &form.UpdatedAt)
  return form , err
}
