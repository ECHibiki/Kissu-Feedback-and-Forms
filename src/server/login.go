package main

import (
  "golang.org/x/crypto/bcrypt"
)

// var login_obfuscator int = 10
//
// func createPassword(cfg types.ConfigurationSettings , db *sql.DB){
//   hashed_bytes , err := bcrypt.GenerateFromPassword([]byte(cfg.ApplicationPassword), login_obfuscator)
//   if err != nil{
//     panic("issue doing bcrypt on password")
//   }
//
//   err = db.Exec("INSERT INTO passwords VALUES (? , ? , ?)" , string(hashed_bytes) , "bcrypt" , "" + login_obfuscator)
//   if err != nil{
//     panic("issue inserting password to DB")
//   }
// }
