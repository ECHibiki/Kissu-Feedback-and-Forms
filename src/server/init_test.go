package main

import (
	"testing"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "os"
)

var initialization_folder string = "../../test"



func TestInitialization(t *testing.T) {

  defer tools.CleanupTestingInitializations()

  is_init := checkIntialized(initialization_folder)
  if is_init {
     t.Fatal("The test folder has been preinitialized")
  }
  intializeProjectToFolder(initialization_folder, init_fields)

  // manually test the existence of /settings/
  _, err := os.Stat(initialization_folder + "/settings/")
  if err != nil {
    t.Fatal("/settings/ still does not exist")
  }
  // manually test the existence of /data/
  _, err := os.Stat(initialization_folder + "/data/")
  if err != nil {
    t.Fatal("/data/ still does not exist")
  }

  is_init = checkIntialized(initialization_folder)
  if !is_init {
    t.Fatal("Despite all this, the folder check fails")
  }

  is_configured := checkConfigured(initialization_folder)
  if is_configured {
     t.Fatal("The config file already exists")
  }
  init_fields := ConfigurationFields{
    DBName: "test-db",
    DBUserName: "test-user",
    DBCredentials: "test-user",
    ApplicationPassword: "test-password",
    StartupPort: "4960",
    SiteName: "example.com",
  }
  initializeConfigurationFile(initialization_folder , init_fields)

  // Manually test the config file exists
  _, err := os.Stat(initialization_folder + "/settings/config.json")
  if err != nil {
    t.Fatal("/settings/ still does not exist")
  }
  is_configured = checkConfigured(initialization_folder)
  if !is_configured {
     t.Fatal("Despite this the config file can not be found")
  }

  initializeFromConfigFile(initialization_folder)
  assumed_config_file := tools.ConfigurationSettings{
    DBName: "test-db",
    DBUserName: "test-user",
    DBCredentials: "test-user",
    StartupPort: "4960",
    SiteName: "example.com",
  }
  if tools.StoredConfig.DBName != assumed_config_file.DBName {
    t.Fatal("DBName mismatch", tools.StoredConfig.DBName , assumed_config_file.DBName )
  }
  if tools.StoredConfig.DBUserName != assumed_config_file.DBUserName {
    t.Fatal("DBUserName mismatch", tools.StoredConfig.DBUserName , assumed_config_file.DBUserName )
  }
  if tools.StoredConfig.DBCredentials != assumed_config_file.DBCredentials {
    t.Fatal("DBCredentials mismatch", tools.StoredConfig.DBCredentials , assumed_config_file.DBCredentials )
  }
  if tools.StoredConfig.StartupPort != assumed_config_file.StartupPort {
    t.Fatal("StartupPort mismatch", tools.StoredConfig.StartupPort , assumed_config_file.StartupPort )
  }
  if tools.StoredConfig.SiteName != assumed_config_file.SiteName {
    t.Fatal("SiteName mismatch", tools.StoredConfig.SiteName , assumed_config_file.SiteName )
  }

  db_active := checkDBConnection()
  if !db_active {
     t.Fatal("The DB should be active")
  }
  // Use the should-be set DB connection to:
    // Manually check the DB was created with propper tables
  _, err := tools.Query("SELECT id, field_json, updated_at FROM forms")
  if err != nil {
    t.Fatal("DB forms error")
  }
  _, err = tools.Query("SELECT fk_id, identifier, response_json, submitted_at FROM responses")
  if err != nil {
    t.Fatal("DB forms error")
  }
  _, err = tools.Query("SELECT hashed_password, hash_system, hash_scrambler FROM passwords")
  if err != nil {
    t.Fatal("DB passwords error")
  }
  _, err = tools.Query("SELECT login, cookie, ip FROM passwords")
  if err != nil {
    t.Fatal("DB passwords error")
  }

  // everything set up
  // remove all connections and handlers on the DB
  closeDBConnection()
  db_active = checkDBConnection()
  if db_active {
     t.Fatal("The DB should be inactive")
  }
  // This is now testing for a second time startup
  is_init = checkIntialized(initialization_folder)
  if !is_init {
     t.Fatal("Project folder should be initialized")
  }
  is_configured = checkConfigured(initialization_folder)
  if !is_configured {
     t.Fatal("The config file should exist")
  }
  connectToDB(initialization_folder)
  db_active = checkDBConnection(initialization_folder)
  if !db_active {
     t.Fatal("The DB should be active")
  }

  // DB is running and files are confirmed to exist
  // Nothing left to do
  // defer: remove data for the next test
}
