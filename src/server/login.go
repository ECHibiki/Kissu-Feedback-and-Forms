package main

import (
  "golang.org/x/crypto/bcrypt"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "database/sql"
  "strconv"
  "errors"
)

// var stored_pass types.PasswordsDBFields = hashPassword( init_fields.ApplicationPassword , "bcrypt" , 10 )
func hashPassword(raw_password string , hash_method string , hash_scrambler string) types.PasswordsDBFields{
  if hash_method == "bcrypt"{
    rounds , err := strconv.Atoi(hash_scrambler)
    if err != nil{
      panic(err)
    }
    hashed_bytes , err := bcrypt.GenerateFromPassword([]byte(raw_password), rounds)
    if err != nil{
      panic(err)
    }
    return types.PasswordsDBFields {
      HashedPassword: string(hashed_bytes),
      HashSystem: hash_method,
      HashScrambler: hash_scrambler,
    }
  } else{
    panic("Unkown password setting used")
  }
}

func checkPasswordValid(submitted_password string, stored_pass string) bool{
  err := bcrypt.CompareHashAndPassword([]byte(stored_pass) , []byte(submitted_password))
  return err == nil
}

func storePassword(db *sql.DB, password_data types.PasswordsDBFields) error{
  stored_pass , err := tools.FetchPassword(db)

  if stored_pass.HashedPassword != "" {
    if err != nil{
      return err
    }
    return errors.New("Setting two values into the password DB results in an error")
  }
  err = tools.WritePassToDB(db, password_data)
  return err
}

func getStoredPassword(db *sql.DB) (types.PasswordsDBFields , error){
  return tools.FetchPassword(db)
}
