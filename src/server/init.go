package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"io/ioutil"
	"os"
)

func checkIntialized(root_dir string) bool {
	// manually test the existence of folders
	var err error
	_, err = os.Stat(root_dir + "/templates/user-views/")
	if err != nil {
		panic("/templates/user-views must exist from download")
	}
	_, err = os.Stat(root_dir + "/templates/mod-views/")
	if err != nil {
		panic("/templates/mod-views must exist from download")
	}
	_, err = os.Stat(root_dir + "/public/css/")
	if err != nil {
		panic("/public/css/ must exist from download")
	}
	_, err = os.Stat(root_dir + "/public/js/")
	if err != nil {
		panic("/public/js/ must exist from download")
	}
	_, err = os.Stat(root_dir + "/public/img/")
	if err != nil {
		panic("/public/img/ must exist from download")
	}

	_, err = os.Stat(root_dir + "/settings/")
	if err != nil {
		return false
	}

	_, err = os.Stat(root_dir + "/data/")
	if err != nil {
		return false
	}

	return true
}
func checkConfigured(root_dir string) bool {
	var cfg types.ConfigurationSettings
	cfg_bytes, err := ioutil.ReadFile(root_dir + "/settings/config.json")
	if err != nil {
		return false
	}

	err = json.Unmarshal(cfg_bytes, &cfg)
	if err != nil {
		panic("Attempting to read config.json resulted in an error")
	}
	if cfg.DBName == "" {
		fmt.Printf("%+v\n", cfg)
		panic("Error reading Config.json. DBName must be set")
	}
	if cfg.DBUserName == "" {
		fmt.Printf("%+v\n", cfg)
		panic("Error reading Config.json. DBUserName must be set")
	}
	if cfg.DBCredentials == "" {
		// No problem
	}
	if cfg.DBAddr == "" {
		fmt.Printf("%+v\n", cfg)
		panic("Error reading Config.json. DBAddr should have a value")
	}
	if cfg.StartupPort == "" {
		fmt.Printf("%+v\n", cfg)
		panic("Error reading Config.json. StartupPort must be set")
	}
	if cfg.SiteName == "" {
		fmt.Printf("%+v\n", cfg)
		panic("Error reading Config.json. SiteName must be set")
	}
	return true
}

func createConfigurationFile(root_dir string, init_fields types.ConfigurationInitializerFields) {
	if init_fields.DBName == "" {
		panic("DBName must be set")
	}
	if init_fields.DBUserName == "" {
		panic("DBUserName must be set")
	}
	if init_fields.DBCredentials == "" {
		// No problem
	}
	if init_fields.DBAddr == "" {
		init_fields.DBAddr = "127.0.0.1"
	}
	if init_fields.ApplicationPassword == "" {
		panic("ApplicationPassword must be set")
	}
	if init_fields.StartupPort == "" {
		panic("StartupPort must be set")
	}
	if init_fields.SiteName == "" {
		panic("SiteName must be set")
	}

	cfg := types.ConfigurationSettings{
		DBName:        init_fields.DBName,
		DBUserName:    init_fields.DBUserName,
		DBCredentials: init_fields.DBCredentials,
		DBAddr:        init_fields.DBAddr,
		StartupPort:   init_fields.StartupPort,
		SiteName:      init_fields.SiteName,
	}

	byte_json, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(root_dir+"/settings/config.json", byte_json, 0655)
	if err != nil {
		panic(err)
	}
}

func createDB(init_fields types.ConfigurationInitializerFields) {
	cfg := types.ConfigurationSettings{
		DBName:        init_fields.DBName,
		DBUserName:    init_fields.DBUserName,
		DBCredentials: init_fields.DBCredentials,
		DBAddr:        init_fields.DBAddr,
		StartupPort:   init_fields.StartupPort,
		SiteName:      init_fields.SiteName,
	}
	db := tools.QuickDBConnect(cfg)
	tools.BuildDBTables(db)
}

//   db, cfg := initializeFromConfigFile(root_dir)
func initializeFromConfigFile(root_dir string) (*sql.DB, types.ConfigurationSettings) {
	var cfg types.ConfigurationSettings
	cfg_bytes, err := ioutil.ReadFile(root_dir + "/settings/config.json")
	if err != nil {
		panic("Config.json missing")
	}
	err = json.Unmarshal(cfg_bytes, &cfg)
	if err != nil {
		panic("Attempting to read config.json resulted in an error")
	}

	db := tools.QuickDBConnect(cfg)
	return db, cfg
}

func createProjectToFolder(root_dir string) {
	// We don't want this to result in a crash, but if it can be created then it should
	_ = os.Mkdir(root_dir+"/settings/", 0755)
	_ = os.Mkdir(root_dir+"/data/", 0755)

}
