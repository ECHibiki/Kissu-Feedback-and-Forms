package tools

import (
  "strings"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "errors"
  "encoding/json"
  "io/ioutil"
  "os"
)

func testTesting() bool{
  // https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
  if strings.HasSuffix(os.Args[0], ".test") {
    return true
  }
  fmt.Println("normal run")
  return false
}

func dropDBOnlyForTesting(db *sql.DB , db_name string) {
  if !testTesting() {
    return
  }
  var err error
  _, err = db.Exec("DROP TABLE responses")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE forms")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE passwords")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE logins")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
}

func connectToDBForTesting(dir string) (types.ConfigurationSettings , *sql.DB , error){
  if !testTesting() {
    return types.ConfigurationSettings{}, nil, errors.New("Not in testing")
  }
  var cfg types.ConfigurationSettings
  cfg_bytes, err := ioutil.ReadFile(dir + "/settings/config.json")
  if err != nil{
    return types.ConfigurationSettings{}, nil, err
  }
  err = json.Unmarshal(cfg_bytes, &cfg)
  if err != nil{
    return types.ConfigurationSettings{}, nil, err
  }

  db , err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
    cfg.DBUserName,
    cfg.DBCredentials,
    cfg.DBAddr,
    cfg.DBName,
    ),
  )
  return cfg , db, err

}

func CleanupTestingInitializations(initialization_folder string){
  if !testTesting() {
    return
  }

  var err error

  cfg , db , err := connectToDBForTesting(initialization_folder)
  if err != nil {
    fmt.Println("err: connectToDBForTesting" , err)
  }
  dropDBOnlyForTesting(db , cfg.DBName)


  err = os.RemoveAll("../../test/settings/")
  if err != nil{
    fmt.Println(err)
  }
  err = os.RemoveAll("../../test/data/")
  if err != nil{
    fmt.Println(err)
  }

}
