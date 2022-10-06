package main

import (
  "golang.org/x/crypto/bcrypt"
  "crypto/sha256"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "database/sql"
  "strconv"
  "errors"
  "time"
  "fmt"
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

func CheckPasswordValid(submitted_password string, stored_pass string) error{
  err := bcrypt.CompareHashAndPassword([]byte(stored_pass) , []byte(submitted_password))
  return err
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
func CheckCookieValid(db  *sql.DB , cookie string , ip string) ( error ){
  login  := db.QueryRow("SELECT * FROM logins WHERE cookie = ? AND ip = ?" , cookie , ip)
  var login_fields types.LoginDBFields
  err := login.Scan( &login_fields.TimeAt , &login_fields.Cookie , &login_fields.IP )
  return err
}
func CreateLoginFields( cookie string , ip string ) ( login types.LoginDBFields ) {
  login.TimeAt = time.Now().Unix()
  login.Cookie = cookie
  login.IP = ip
  return
}
func CreateAuthenticationHash( key string ) ( string ) {
  sha_key := sha256.Sum256([]byte(key))
  return fmt.Sprintf("%x", sha_key)
}
func StoreLogin(db *sql.DB , fields types.LoginDBFields ) error {
  _ , err := db.Exec("INSERT INTO logins VALUES( ? , ? , ? )" , fields.TimeAt , fields.Cookie, fields.IP)
  return err
}
