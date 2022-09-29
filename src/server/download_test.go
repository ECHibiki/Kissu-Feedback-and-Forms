package main
import(
  "testing"
)

func TestCompressionOfGivenForm(t *testing.T){
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

  // a CSV of all replies for given form data, allowing for Excel viewing
  // Placed into the form's directory
    // Steps: Get form for all Columns from FormConstruct and initialize a map[string][]string, get all rows and fill into the map
    //        In order defined by FormConstruct, createa a [][]sring for CSV creation
  tools.CreateInstancedCSVForGivenForm(2)
  // A tar.gz file containing the CSV, as it has zipped the entire form directory together
  tools.CreateDownloadableForGivenForm(2)

  return
  // serving of files by http done without tests..
}
