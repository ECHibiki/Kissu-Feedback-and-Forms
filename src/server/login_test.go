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

  db, init_fields , _ := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  var stored_pass types.PasswordsDBFields  = hashPassword( init_fields.ApplicationPassword , "bcrypt" , "10" )
  if stored_pass.HashedPassword == init_fields.ApplicationPassword {
    t.Fatal("HashedPassword is same as init_fields.ApplicationPassword")
  }
  if stored_pass.HashSystem != "bcrypt"{
    t.Fatal("HashSystem is not bcrpyt")
  }
  if stored_pass.HashScrambler != "10"{
    t.Fatal("HashScrambler is not using a tested value")
  }
  valid = checkPasswordValid(init_fields.ApplicationPassword , stored_pass.HashedPassword)
  if !valid{
    t.Fatal("Assigned password does not register as correct before storage")
  }

  err = storePassword( db , stored_pass )
  if err != nil{
    t.Fatal("Error on password storage" , err)
  }
  var a_second_stored_pass types.PasswordsDBFields = hashPassword( "second-" + init_fields.ApplicationPassword , "bcrypt" , "10" )
  err = storePassword( db , a_second_stored_pass )
  if err == nil{
    t.Fatal( err )
  }
  // write to it anyways, immitating the effect of a manual DB insertion. It shouldn't effect the outcome
  tools.WritePassToDB(db , a_second_stored_pass)
  potentially_the_second_stored_pass , err := getStoredPassword( db )
  if err != nil{
    t.Fatal( err )
  }
  if potentially_the_second_stored_pass.HashedPassword == a_second_stored_pass.HashedPassword{
    t.Fatal("Passwords being read incorrectly from DB, password table with two collumns should always be reading the top-most, first inserted value")
  }

  retrieved_pass , err := getStoredPassword( db )
  if err != nil{
    t.Fatal( err )
  }
  if stored_pass.HashedPassword != retrieved_pass.HashedPassword {
    t.Fatal("HashedPassword was not stored correctly" , retrieved_pass, stored_pass)
  }
  if stored_pass.HashSystem != retrieved_pass.HashSystem{
    t.Fatal("HashSystem was not stored correctly" , retrieved_pass, stored_pass)
  }
  if stored_pass.HashScrambler != retrieved_pass.HashScrambler{
    t.Fatal("HashScrambler was not stored correctly", retrieved_pass, stored_pass)
  }
  valid = checkPasswordValid(init_fields.ApplicationPassword , retrieved_pass.HashedPassword)
  if !valid{
    t.Fatal("Assigned password does not register as correct after storage" , init_fields.ApplicationPassword , retrieved_pass)
  }

  valid = checkPasswordValid("Not-" + init_fields.ApplicationPassword , retrieved_pass.HashedPassword)
  if valid{
    t.Fatal("The incorrect password registers as valid" , init_fields.ApplicationPassword , retrieved_pass)
  }
}
