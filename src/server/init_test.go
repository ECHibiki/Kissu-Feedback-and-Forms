package main

import (
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
	"os"
	"fmt"
	"testing"
	"net/http/httptest"
	"net/http"

)

func TestInitialization(t *testing.T) {
	var initialization_folder string = "../../test"

	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	// Check for basic folders
	is_init := checkIntialized(initialization_folder)
	if is_init {
		t.Fatal("The test folder has been preinitialized")
	}
	createProjectToFolder(initialization_folder)

	// manually test the existence of /settings/
	_, err := os.Stat(initialization_folder + "/settings/")
	if err != nil {
		t.Fatal("/settings/ still does not exist")
	}
	// manually test the existence of /data/
	_, err = os.Stat(initialization_folder + "/data/")
	if err != nil {
		t.Fatal("/data/ still does not exist")
	}

	is_init = checkIntialized(initialization_folder)
	if !is_init {
		t.Fatal("Despite all this, the folder check fails")
	}

	// test config.json
	is_configured := checkConfigured(initialization_folder)
	if is_configured {
		t.Fatal("The config file already exists")
	}
	init_fields := types.ConfigurationInitializerFields{
		DBName:              "feedback_tests",
		DBUserName:          "testuser",
		DBCredentials:       "",
		DBAddr:              "127.0.0.1",
		ApplicationPassword: "test-password",
		StartupPort:         ":4960",
		SiteName:            "example.com",
	}
	//to create files and data
	createConfigurationFile(initialization_folder, init_fields)
	createDB(init_fields)

	// Manually test the config file exists
	_, err = os.Stat(initialization_folder + "/settings/config.json")
	if err != nil {
		t.Fatal("/settings/ still does not exist")
	}
	is_configured = checkConfigured(initialization_folder)
	if !is_configured {
		t.Fatal("Despite this the config file can not be found")
	}

	// to load existing data
	db, cfg := initializeFromConfigFile(initialization_folder)
	if db == nil {
		t.Fatal("The DB should be active")
	}
	assumed_config_file := types.ConfigurationSettings{
		DBName:        "feedback_tests",
		DBUserName:    "testuser",
		DBCredentials: "",
		DBAddr:        "127.0.0.1",
		StartupPort:   ":4960",
		SiteName:      "example.com",
	}
	if cfg.DBName != assumed_config_file.DBName {
		t.Fatal("DBName mismatch", cfg.DBName, assumed_config_file.DBName)
	}
	if cfg.DBUserName != assumed_config_file.DBUserName {
		t.Fatal("DBUserName mismatch", cfg.DBUserName, assumed_config_file.DBUserName)
	}
	if cfg.DBCredentials != assumed_config_file.DBCredentials {
		t.Fatal("DBCredentials mismatch", cfg.DBCredentials, assumed_config_file.DBCredentials)
	}
	if cfg.StartupPort != assumed_config_file.StartupPort {
		t.Fatal("StartupPort mismatch", cfg.StartupPort, assumed_config_file.StartupPort)
	}
	if cfg.SiteName != assumed_config_file.SiteName {
		t.Fatal("SiteName mismatch", cfg.SiteName, assumed_config_file.SiteName)
	}

	// Use the should-be set DB connection to:
	// Manually check the DB was created with propper tables
	_, err = db.Query("SELECT id, name, field_json, updated_at FROM forms")
	if err != nil {
		t.Fatal("DB forms error", "DB Name "+cfg.DBName+" err: ", err)
	}
	_, err = db.Query("SELECT id, fk_id, identifier, response_json, submitted_at FROM responses")
	if err != nil {
		t.Fatal("DB responses error", "DB Name "+cfg.DBName+" err: ", err)
	}
	_, err = db.Query("SELECT hashed_password, hash_system, hash_scrambler FROM passwords")
	if err != nil {
		t.Fatal("DB passwords error", "DB Name "+cfg.DBName+" err: ", err)
	}
	_, err = db.Query("SELECT time_at, cookie, ip FROM logins")
	if err != nil {
		t.Fatal("DB logins error", "DB Name "+cfg.DBName+" err: ", err)
	}

	// everything set up
	// remove all connections and handlers on the DB
	db = nil
	// This is now testing for a second time startup
	is_init = checkIntialized(initialization_folder)
	if !is_init {
		t.Fatal("Project folder should be initialized")
	}
	is_configured = checkConfigured(initialization_folder)
	if !is_configured {
		t.Fatal("The config file should exist")
	}
	db, cfg = initializeFromConfigFile(initialization_folder)
	if db == nil {
		t.Fatal("The DB should be active")
	}

	// DB is running and files are confirmed to exist
	// Nothing left to do
	// defer: remove data for the next test
}

func TestRuning(t *testing.T){
	var initialization_folder string = "./../../test"
	db, _ , cfg := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)
 	templater.SetRootDir(initialization_folder + "/templates/")
	gin_engine := routeGin(&cfg, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	gin_engine.ServeHTTP(w, req)

	if 200 != w.Code {
		t.Fatal("Bad code" , w.Code)
	}

	expected := "<!DOCTYPE html>\n<HTML>\n  <HEAD>\n      <META http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n  <META name=\"Language\" content=\"en\">\n  <META name=\"robots\" content=\"index, follow, all\">\n  <SCRIPT src=\"/assets/js/main(" + globals.ProjectVersion + ").js\"></SCRIPT>\n  <LINK href=\"/assets/css/style.css\" rel=\"stylesheet\">\n  <LINK rel=\"shortcut icon\" href=\"/assets/img/favicon.png\" type=\"image/png\">\n\n    <TITLE>Kissu Feedback & Forms - Home</TITLE>\n  </HEAD>\n  <BODY>\n    <P>Please access forms directly through their URL</P>\n    <HR/>\n    Feedback & Forms V" + globals.ProjectVersion + " - Product of Kissu.moe under MPL2.0<br/>\nRun it yourself from <a href=\"https://github.com/ECHibiki/Kissu-Feedback-and-Forms\">the source code</a>\n\n  </BODY>\n</HTML>\n"
	if expected != w.Body.String() {
		t.Fatal(fmt.Sprintf("%#v" , expected) , "\n" , fmt.Sprintf("%#v" ,w.Body.String()))
	}
}
