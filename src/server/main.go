package main
import (
  "fmt"
  "database/sql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
)

func main(){
  fmt.Println(`Starting Kissu Feedback & Forms
    Project of ECVerniy for Kissu.moe
    All files herein are open-source under the MPL-2.0 license
    Feedback & Forms hopes to aid communities by promoting communication between staff and users
--------------`)
  // initialize settings
    // intialize database
    // initialize folders
  is_init := checkIntialized(globals.RootDirectory)
  if !is_init {
    createProjectToFolder(globals.RootDirectory)
  }

  var db *sql.DB
  var cfg types.ConfigurationSettings
  is_configured := checkConfigured(globals.RootDirectory)
  if !is_configured {
    var init_fields types.ConfigurationInitializerFields
    // console inputs or pass
    fmt.Println("Configuration file not found. Creating one through console inputs...")
    init_fields.DBName = tools.ReadLine("DBName")
    init_fields.DBUserName = tools.ReadLine("DBUserName")
    init_fields.DBCredentials = tools.ReadLine("DBCredentials")
    init_fields.DBAddr = tools.ReadLine("DBAddr")
    init_fields.ApplicationPassword = tools.ReadLine("ApplicationPassword")
    init_fields.StartupPort = tools.ReadLine("StartupPort")
    init_fields.SiteName = tools.ReadLine("SiteName")

    createConfigurationFile(globals.RootDirectory, init_fields)
    createDB(init_fields)
    var password_data types.PasswordsDBFields = hashPassword( init_fields.ApplicationPassword , "bcrypt" , "10" )

    db , cfg = initializeFromConfigFile(globals.RootDirectory)
    storePassword(db , password_data)
  } else{
    db , cfg = initializeFromConfigFile(globals.RootDirectory)
  }
  stick := templater.ReturnFileSystemTemplateHandler(globals.RootDirectory)

  gin_engine := routeGin(&cfg, db , stick )
  runGin(gin_engine , cfg.StartupPort)
}
