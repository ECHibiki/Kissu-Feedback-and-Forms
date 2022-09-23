package main

import (
	"testing"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"

)


func TestPassCreation(t *testing.T)  {
  var initialization_folder string = "../../test"
  var err error
  var valid bool

  db, cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  var stored_pass types.PasswordsDBFields = hashPassword( cfg.ApplicationPassword , "bcrypt" , 10 )
  if stored_pass.HashedPassword == cfg.ApplicationPassword {
    t.Fatal("HashedPassword is same as cfg.ApplicationPassword")
  }
  if stored_pass.HashSystem != "bcrypt"{
    t.Fatal("HashSystem is not bcrpyt")
  }
  if stored_pass.HashScrambler != "10"{
    t.Fatal("HashScrambler is not using a tested value")
  }
  valid = checkPasswordValid(cfg.ApplicationPassword , stored_pass)
  if !valid{
    t.Fatal("Assigned password does not register as correct before storage")
  }

  err = storePassword( db , stored_pass )
  if err != nil{
    t.Fatal("Error on password storage" , err)
  }
  var a_second_stored_pass types.PasswordsDBFields = hashPassword( "second-" + cfg.ApplicationPassword , "bcrypt" , 10 )
  err = storePassword( db , a_second_stored_pass )
  if err == nil{
    t.Fatal("Setting two values into the password DB should result in an error" )
  }
  // write to it anyways, immitating the effect of a manual DB insertion. It shouldn't effect the outcome
  tools.WritePassToDB(db , a_second_stored_pass)
  var potentially_the_second_stored_pass types.PasswordsDBFields = getStoredPassword( db )
  if potentially_the_second_stored_pass.HashedPassword == a_second_stored_pass.HashedPassword{
    t.Fatal("Passwords being read unexpectedly from DB, password table with two collumns should always be reading the top-most, first inserted value")
  }

  var retrieved_pass types.PasswordsDBFields = getStoredPassword( db )
  if stored_pass.HashedPassword != retrieved_pass.HashedPassword {
    t.Fatal("HashedPassword was not stored correctly" , retrieved_pass, stored_pass)
  }
  if stored_pass.HashSystem != retrieved_pass.HashSystem{
    t.Fatal("HashSystem was not stored correctly" , retrieved_pass, stored_pass)
  }
  if stored_pass.HashScrambler != retrieved_pass.HashScrambler{
    t.Fatal("HashScrambler was not stored correctly", retrieved_pass, stored_pass)
  }
  valid = checkPasswordValid(cfg.ApplicationPassword , retrieved_pass)
  if !valid{
    t.Fatal("Assigned password does not register as correct after storage" , cfg.ApplicationPassword , retrieved_pass)
  }

  valid = checkPasswordValid("Not-" + cfg.ApplicationPassword , retrieved_pass)
  if valid{
    t.Fatal("The incorrect password registers as valid" , cfg.ApplicationPassword , retrieved_pass)
  }
}
