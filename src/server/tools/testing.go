package tools

import (
  "strings"
  "fmt"
  "os"
)

func CleanupTestingInitializations(initialization_folder string){
  if strings.HasSuffix(os.Args[0], ".test") {
    fmt.Println("normal run")
    return
  }

  connectToDB(initialization_folder)
  // query to remove the given database
  dropDBOnlyForTesting(initialization_folder)

  os.Remove(initialization_folder)
  os.Mkdir(initialization_folder, 0755)
}

// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
func BasicIntialization(initialization_folder string){
  if strings.HasSuffix(os.Args[0], ".test") {
    fmt.Println("normal run")
    return
  }

  intializeProjectToFolder(initialization_folder, init_fields)
  init_fields := ConfigurationFields{
    DBName: "test-db",
    DBUserName: "test-user",
    DBCredentials: "test-user",
    ApplicationPassword: "test-password",
    StartupPort: "4960",
    SiteName: "example.com",
  }
  initializeConfigurationFile(initialization_folder , init_fields)
}
