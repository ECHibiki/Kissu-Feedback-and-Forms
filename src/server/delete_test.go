package main

import(
  "testing"
)

func TestDeleteOfForm(t *testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  first_name := "Test form 1"
  first_store := "Test_form_1"
  tools.DoFormInitialization(first_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  second_name := "Test form 2"
  second_store := "Test_form_2"
  tools.DoFormInitialization(first_store , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  testing.ReplyToForm(db , 1 , second_store, "192.168.1.3", db , initialization_folder)
  testing.ReplyToForm(db , 2 , first_store , "192.168.1.1" , db , initialization_folder)
  testing.ReplyToForm(db , 2 , first_store, "192.168.1.2", db , initialization_folder)

  destroyer.DeleteForm(db , initialization_folder , 2)

  _ , err := tools.GetFormOfID(db , 2)
  if err != nil {
      t.Fatal("deleted form still exists")
  }
  _, err := os.Stat(initialization_folder + "/data/" + second_store)
  if err != nil {
    t.Fatal("Form deletes should retain old files in case of mistakes")
  }


  _ , err := tools.GetFormOfID(db , 1)
  if err == nil {
    t.Fatal("Error on form that should still exist")
  }
  _, err = os.Stat(initialization_folder + "/data/" + first_store)
  if err != nil {
    t.Fatal("An unrelated form directory was removed")
  }
}

func TestDeleteOfResponse(t * testing.T){
  var initialization_folder string = "../../test"
  var err error

  db, _ , cfg := tools.DoTestingIntializations(initialization_folder)
  defer tools.CleanupTestingInitializations(initialization_folder)

  first_name := "Test form 1"
  first_store := "Test_form_1"
  tools.DoFormInitialization(first_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  second_name := "Test form 2"
  second_store := "Test_form_2"
  tools.DoFormInitialization(second_name , "a-simple-identifier" , db , cfg)
  if err != nil {
    t.Fatal(err)
  }
  testing.ReplyToForm(db , 1 , second_store, "192.168.1.3", db , initialization_folder)
  testing.ReplyToForm(db , 2 , first_store , "192.168.1.1" , db , initialization_folder)
  testing.ReplyToForm(db , 2 , first_store, "192.168.1.2", db , initialization_folder)

  destroyer.DeleteResponse(db , initialization_folder , 2)

  for i := 1 ; i < 3 ; i++ {
    _ , err := tools.GetFormOfID(db , i)
    if err != nil {
      t.Fatal("Error on form that should still exist")
    }
    _, err = os.Stat(initialization_folder + "/data/" + first_store)
    if err != nil {
      t.Fatal("An unrelated form directory was removed")
    }

  }

  _ , err := tools.GetResponseByID(db , 1)
  if err == nil {
    t.Fatal("Error on response that should still exist")
  }
  _, err = os.Stat(initialization_folder + "/data/" + first_store + "/192.168.1.1/")
  if err != nil {
    t.Fatal("An unrelated Response directory was removed")
  }
  _ , err := tools.GetResponseByID(db , 3)
  if err == nil {
    t.Fatal("Error on response that should still exist")
  }
  _, err = os.Stat(initialization_folder + "/data/" + first_store + "/192.168.1.3/")
  if err != nil {
    t.Fatal("An unrelated Response directory was removed")
  }

  _ , err := tools.GetResponseByID(db , 2)
  if err != nil {
    t.Fatal("Response that should not exist does not return SQL err")
  }
  _, err = os.Stat(initialization_folder + "/data/" + first_store + "/192.168.1.2/")
  if err == nil {
    t.Fatal("Response that should not exist does not return OS err")
  }

}
