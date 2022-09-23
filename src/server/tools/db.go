package tools

import(
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
)

type FormDBFields struct{
  id int64
  field_json string
  updated_at int64
}
type ResponseDBFields struct{
  id int64
  fk_id int64
  identifier string
  response_json string
  submitted_at int64
}
type PasswordsDBFields struct{
  hashed_password string
  hash_system string
  hash_scrambler string
}
type LoginsDBFields struct{
  time_at int64
  cookie string
  ip string
}

func BuildDBTables(init_fields types.ConfigurationSettings  , db *sql.DB ){
  _, err := db.Exec("CREATE TABLE forms(id INT PRIMARY KEY , field_json TEXT , updated_at INT)")
  if err != nil{
    panic(err)
  }
  _, err = db.Exec("CREATE TABLE responses(id INT PRIMARY KEY , fk_id INT, identifier VARCHAR(60) , response_json TEXT , submitted_at INT  , FOREIGN KEY (fk_id) REFERENCES forms(id) )")
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
