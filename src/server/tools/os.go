package tools

import (
  "os"
)

func LogError(storage_dir string, message string){
  err_handler , err := os.OpenFile(storage_dir +  "errors.log" , os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    panic(err)
  }
  defer err_handler.Close()
  _ , err = err_handler.WriteString("File write fail: " + message + "\n")
  if err != nil {
    panic(err)
  }
}
