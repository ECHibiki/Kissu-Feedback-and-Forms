package tools

import (
  "strings"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/builder"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "errors"
  "encoding/json"
  "io/ioutil"
  "os"
  "mime/multipart"
  // "net/textproto"
  // "path/filepath"
)

func testTesting() bool{
  // https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
  if strings.HasSuffix(os.Args[0], ".test") {
    return true
  }
  fmt.Println("normal run")
  return false
}

func dropDBOnlyForTesting(db *sql.DB , db_name string) {
  if !testTesting() {
    return
  }
  var err error
  _, err = db.Exec("DROP TABLE responses")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE forms")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE passwords")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
  _, err = db.Exec("DROP TABLE logins")
  if err != nil{
    fmt.Println("err: dropDBOnlyForTesting " , err)
  }
}

func connectToDBForTesting(dir string) (types.ConfigurationSettings , *sql.DB , error){
  if !testTesting() {
    return types.ConfigurationSettings{}, nil, errors.New("Not in testing")
  }
  var cfg types.ConfigurationSettings
  cfg_bytes, err := ioutil.ReadFile(dir + "/settings/config.json")
  if err != nil{
    return types.ConfigurationSettings{}, nil, err
  }
  err = json.Unmarshal(cfg_bytes, &cfg)
  if err != nil{
    return types.ConfigurationSettings{}, nil, err
  }

  db , err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
    cfg.DBUserName,
    cfg.DBCredentials,
    cfg.DBAddr,
    cfg.DBName,
    ),
  )
  return cfg , db, err

}

func CleanupTestingInitializations(initialization_folder string){
  if !testTesting() {
    return
  }

  var err error

  cfg , db , err := connectToDBForTesting(initialization_folder)
  if err != nil {
    fmt.Println("err: connectToDBForTesting" , err)
  }
  dropDBOnlyForTesting(db , cfg.DBName)


  err = os.RemoveAll("../../test/settings/")
  if err != nil{
    fmt.Println(err)
  }
  err = os.RemoveAll("../../test/data/")
  if err != nil{
    fmt.Println(err)
  }

}

func DoTestingIntializations(initialization_folder string) (*sql.DB , types.ConfigurationInitializerFields , types.ConfigurationSettings){
  err := os.Mkdir(initialization_folder + "/settings/", 0755)
  if err != nil {
    panic("Initialization of project settings folder failed")
  }
  err = os.Mkdir(initialization_folder + "/data/", 0755)
  if err != nil {
    panic("Initialization of project data folder failed")
  }
  init_fields := types.ConfigurationInitializerFields{
    DBName: "feedback_tests",
    DBUserName: "testuser",
    DBCredentials: "",
    DBAddr: "127.0.0.1",
    ApplicationPassword: "test-password",
    StartupPort: ":4960",
    SiteName: "example.com",
    ResoruceDirectory: initialization_folder,
  }
  cfg := types.ConfigurationSettings{
    DBName: init_fields.DBName,
    DBUserName: init_fields.DBUserName,
    DBCredentials: init_fields.DBCredentials,
    DBAddr: init_fields.DBAddr,
    StartupPort: init_fields.StartupPort,
    SiteName: init_fields.SiteName,
    ResoruceDirectory: init_fields.ResoruceDirectory,
  }

  byte_json , err := json.Marshal(cfg)
  if err != nil {
    panic(err)
  }
  err = ioutil.WriteFile(initialization_folder + "/settings/config.json", byte_json, 0655)
  if err != nil {
    panic(err)
  }

  db := QuickDBConnect(cfg)
  BuildDBTables( db )
  return db, init_fields , cfg
}

func DoFormInitialization(form_name string, form_id string, db *sql.DB, ){
  var base_demo_form former.FormConstruct = former.FormConstruct{
      FormName: form_name ,
      ID: form_id,
      Description: "First test form",
      AnonOption: true,
      FormFields:[]former.FormGroup{
        {
          Label:"test-group1",
          ID: "test-group1",
          Description: "Groups and subgroups may have a description, when set it does not need respondables",
          // SubGroup: []former.FormGroup{},
          Respondables:[]former.UnmarshalerFormObject{
              {
                Type: former.TextAreaTag ,
                Object: former.TextArea{
                  Field: former.Field{
                    Label:"Test-Text-Area",
                    Name:"Test-TA",
                    Required:true,
                  },
                  Placeholder:"This is a test TA",
                },
              } ,
              {
                Type: former.GenericInputTag ,
                Object: former.GenericInput{
                  Field: former.Field{
                    Label:"Test-GenericInput",
                    Name:"Test-GI",
                    Required:true,
                  },
                  Placeholder:"This is a test GI",
                  Type:former.Text, // former.InputType
                },
              } ,
              {
                Type: former.FileInputTag ,
                Object: former.FileInput{
                  Field: former.Field{
                    Label:"Test-FileInput",
                    Name:"Test-FI",
                    Required:false,
                  },
                    AllowedExtRegex:"jpg",
                    MaxSize:200000, // ~200kb
                  },
              },
              {
                Type: former.FileInputTag ,
                Object: former.FileInput{
                  Field: former.Field{
                    Label:"Test-FileInput",
                    Name:"Test-FI-2",
                    Required:false,
                  },
                    AllowedExtRegex:"jpg",
                    MaxSize:10000000, // ~10mb
                },
              } ,
              {
                Type: former.SelectionGroupTag ,
                Object: former.SelectionGroup{
                  Field: former.Field{
                    Label:"Test-Chk-SelectGroup",
                    Name:"Test-Chk-SG",
                    Required:true,
                  },
                  SelectionCategory: former.Checkbox,
                  CheckableItems:[]former.Checkable{
                    {Label:"A check Item", Value:"ck1"},
                    {Label:"Another check Item", Value:"ck2"},
                    {Label:"final check Item", Value:"ck3"},
                  },
                },
              },
              {
                Type: former.SelectionGroupTag,
                Object: former.SelectionGroup{
                  Field: former.Field{
                    Label:"Test-rdo-SelectGroup",
                    Name:"Test-rdo-SG",
                    Required:true,
                  },
                  SelectionCategory: former.Radio,
                  CheckableItems:[]former.Checkable{
                    {Label:"A radio Item", Value:"rd1"},
                    {Label:"Another radio Item", Value:"rd2"},
                  },
                },
              },
              {
                Type: former.OptionGroupTag,
                Object: former.OptionGroup{
                  Field: former.Field{
                    Label:"Test-optGrp",
                    Name:"Test-optGrp",
                    Required:true,
                  },
                  Options:[]former.OptionItem{
                    {
                      Label:"Item 1",
                      Value: "item-1",
                    } ,
                    {
                      Label:"Item 2",
                      Value: "item-2",
                    } ,
                  },
                },
              },
          },
        },
      },
  }
  issue_array := builder.ValidateForm(base_demo_form)
  if len(issue_array) != 0 {
    fmt.Println(issue_array)
    panic("Issue array, issues detected")
  }
  var insertable_form types.FormDBFields
  insertable_form , err :=  builder.MakeFormWritable(base_demo_form)
  if err != nil{
    panic(err)
  }
  err = StoreFormToDB(db, insertable_form)
  if err != nil{
    panic(err)
  }
}

// Place files into the tmp folder from the testing-data folder
func CopyTestFilesToMemory(root_dir string , image_names map[string]string ) ( map[string]former.MultipartFile ) {
  var processed_images map[string]former.MultipartFile = make(map[string]former.MultipartFile)

  for key , fname := range image_names {
    // type FileHeader struct {
    // 	Filename string
    // 	Header   textproto.MIMEHeader
    // 	Size     int64
    // 	// contains filtered or unexported fields
    // }

    // type File interface {
      // 	io.Reader
      // 	io.ReaderAt
      // 	io.Seeker
      // 	io.Closer
    // }
    // File is an interface to access the file part of a multipart message. Its contents may be either stored in memory or on disk. If stored on disk, the File's underlying concrete type will be an *os.File.
    // Read Only .Open
    file_handle , err := os.Open(root_dir + "/testing-data/images/" + fname)
    if err != nil{
      panic(err)
    }
    finfo , err := file_handle.Stat()
    if err != nil{
      panic(err)
    }
    size := finfo.Size()
    processed_images[key] = former.MultipartFile {
      File: file_handle,
      Header: &multipart.FileHeader{
        Filename: fname,
        Header: nil , // we won't use this unless we must
        Size: size,
      },
    }
  }
  return  processed_images
}
