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

func DoTestingIntializations(initialization_folder string) (*sql.DB , types.ConfigurationInitializerFields , types.ConfigurationSettings){
  err := os.Mkdir(initialization_folder + "/settings/", 0755)
  if err != nil {
    panic("Initialization of project settings folder failed")
  }
  err = os.Mkdir(initialization_folder + "/data/", 0755)
  if err != nil {
    panic("Initialization of project data folder failed")
  }
  init_fields := types.ConfigurationInitializerFields{
    DBName: "feedback_tests",
    DBUserName: "testuser",
    DBCredentials: "",
    DBAddr: "127.0.0.1",
    ApplicationPassword: "test-password",
    StartupPort: ":4960",
    SiteName: "example.com",
    ResoruceDirectory: initialization_folder,
  }
  cfg := types.ConfigurationSettings{
    DBName: init_fields.DBName,
    DBUserName: init_fields.DBUserName,
    DBCredentials: init_fields.DBCredentials,
    DBAddr: init_fields.DBAddr,
    StartupPort: init_fields.StartupPort,
    SiteName: init_fields.SiteName,
    ResoruceDirectory: init_fields.ResoruceDirectory,
  }

  byte_json , err := json.Marshal(cfg)
  if err != nil {
    panic(err)
  }
  err = ioutil.WriteFile(initialization_folder + "/settings/config.json", byte_json, 0655)
  if err != nil {
    panic(err)
  }

  db := QuickDBConnect(cfg)
  BuildDBTables( db )
  return db, init_fields , cfg
}

func DoFormInitialization(){

}
