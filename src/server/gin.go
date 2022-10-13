package main

import (
  "github.com/gin-gonic/gin"
  "github.com/tyler-sommer/stick"
  "database/sql"
  "encoding/json"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/builder"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/destroyer"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/former/responder"

  "os"
  "fmt"
  "strconv"
  "time"
  "net/http"
)

func serveTwigTemplate(c *gin.Context , status int , template string){
    c.Header("Content-Type", "text/html")
    c.String(status , template)
}

// bundle args into a struct
func routeGin(cfg *types.ConfigurationSettings, db *sql.DB , stick *stick.Env ) *gin.Engine{

 // use args flags to set
 var gin_mode string
 if len(os.Args) > 1 && os.Args[1] == "--release" {
   gin_mode = "release"
 } else {
   gin_mode = "debug"
 }
 gin.SetMode( gin_mode )

 gin_engine := gin.Default()
 gin_engine.SetTrustedProxies([]string{"127.0.0.1"})

 {
   gin_engine.Static( "/assets", globals.RootDirectory + "public" ) //

   gin_engine.GET("/", generalGetHomepageHandler( stick )) //
   gin_engine.POST("/", modPostLoginForm( db , cfg , stick ))

    public_group := gin_engine.Group("/public")
    {
      public_group.GET("/", generalGetHomepageHandler( stick )) //
      // Handle form requests and build forms
      public_group.GET("/forms/:formname/:formnum", userServeForm(db , stick)) //
      public_group.POST("/forms/:formname/:formnum", userPostForm( db )) //

    }

   // Verify authentication down this route
   mod_group := gin_engine.Group("/mod")
   mod_group.Use(authenticationMiddleware( db  , stick ))
   {
     // list menu CREATE/VIEW
     mod_group.GET("/", modServeHomepageHandler( stick )) //
     // build a form
     mod_group.GET("/create", modServeCreateForm( stick )) //
     mod_group.POST("/create", modPostCreateForm( db , cfg )) //
     // edit a form
     mod_group.GET("/edit/:formnum", modServeEditForm( db  , stick ))
     mod_group.POST("/edit/:formnum", modPostEditForm( db ,  cfg ))
     // delete forms
     mod_group.POST("/form/delete/:formname/:formnum", modPostDeleteForm( db ))
     mod_group.POST("response/delete/:formname/:respnum", modPostDeleteResponse( db ))

     mod_group.GET("/form/delete/:formname/:formnum", modServeDelete( ))
     mod_group.GET("response/delete/:formname/:respnum", modServeDelete( ))
     // view all forms
     mod_group.GET("/view/", modServeViewAllForms(db  , stick))
     // view a form with responses
     mod_group.GET("/view/:formnum", modServeViewSingleForm( db  , stick )) // 6
     // view a response
     mod_group.GET("/view/:formnum/:respnum", modServeViewSingleResponse(db  , stick))
     // download everything of a form
     mod_group.GET("/download/:formname/:formnum", modServeDownloadForm(db ))
     mod_group.GET("/download/:formname/downloadable.tar.gz", modDownloadableForm())
     // Retrieve various forms
     api_group := mod_group.Group("api/")
     {
       // API calls to return details for given forms
       // Might not use these afterall and go for a classical HTML interface
       // Flesh out if mod pages are to be dynamic ReactJS
       api_group.GET("/all", modServeAPIGetAll())
       api_group.GET("/form/:formnum", modServeAPIGetForm())
       api_group.GET("/response/:respnum", modServeAPIGetResponse())

     }

   }
 }
 return gin_engine
}

func runGin(gin_engine *gin.Engine, port string ){
  gin_engine.Run(":" + port)
}

/* GET Handlers */
// Perform functions using the /former/ packages


func placeholderHandler() gin.HandlerFunc {
  return func (c *gin.Context) {
    fmt.Println("Placeholder ran")
    c.String(http.StatusOK , "Placeholder")
   }
}

func generalGetHomepageHandler(env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    template , err := templater.ReturnFilledTemplate(env , "user-views/user-home.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Home generation failed"  } )
      return
    }
    serveTwigTemplate(c , http.StatusOK , template)
  }
}

// Handle route to /mod
func modServeHomepageHandler(env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-home.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Homepage generation failed"} )
      return
    }
    serveTwigTemplate(c , http.StatusOK , template)
  }
  //Oct3
  // login
  // Page for Create and View
}

// Handle route /mod/create
func modServeCreateForm(env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-create.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Creator generation failed"} )
      return
    }
    serveTwigTemplate(c , http.StatusOK , template)
  }
  //Oct3
  // Display the form builder and JS to get it to work
}

// Handle /mod/edit/FORMNUMBER/
func modServeEditForm(db *sql.DB , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    formnum := c.Param("formnum")
    num , err := strconv.Atoi(formnum)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "URI not a formnumber"} )
      return
    }
    form_data , err := tools.GetFormOfID(db , int64(num))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Can't find form"} )
      return
    }
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form_data.FieldJSON) , &form_construct)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Invalid Unmarshaling of form"} )
      return
    }
    fmt.Println(form_data.FieldJSON)
    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-edit.twig" , map[string]stick.Value{
      "version" : globals.ProjectVersion , "id": form_data.ID , "form" : form_construct , "form_str" : form_data.FieldJSON,
    })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
      return
    }
    serveTwigTemplate(c , http.StatusOK , template)
  }
  // Display the form builder and JS to get it to work
}

// Handle route /mod/view and /mod/api/all
func modServeViewAllForms(db *sql.DB , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_data_list , err:= returner.GetAllForms(db)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Can't get forms"} )
      return
    }
    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-form-list.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "form_list" : form_data_list })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
      return
    }
    fmt.Println(form_data_list)
    serveTwigTemplate(c , http.StatusOK , template)
  }
  // view form list
}

// Handle route /mod/view/FORMNUMBER and /mod/api/#
func modServeViewSingleForm(db *sql.DB , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_num , err := strconv.Atoi(c.Param("formnum"))
    form_data , err := tools.GetFormOfID(db , int64(form_num))
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form_data.FieldJSON) , &form_construct)
    form_replies , err := returner.GetRepliesToForm(db , int64(form_num))
    var reply_list []map[string]string
    for _ , r := range form_replies{
      var r_map map[string]string
      err = json.Unmarshal([]byte(r.ResponseJSON) , &r_map)
      if err != nil{
        fmt.Println(err)
        c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a reply"} )
        return
      }
      fmt.Println(r_map)
      r_map["ID"] = strconv.Itoa(int(r.ID))
      r_map["FK_ID"] = strconv.Itoa(int(r.FK_ID))
      r_map["Identifier"] = r.Identifier
      r_map["SubmittedAt"] = strconv.Itoa(int(r.SubmittedAt))
      reply_list = append(reply_list , r_map)
    }
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Can't get form replies"} )
      return
    }
    fmt.Println(reply_list)
    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-reply-list.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "form" : form_construct , "replies": reply_list })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
      return
    }
    serveTwigTemplate(c , http.StatusInternalServerError , template)
  }
}

// Handle route /mod/view/FORMNUMBER/RESPONSENUMBER and /mod/api/#/#
func modServeViewSingleResponse(db *sql.DB , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_num , err := strconv.Atoi(c.Param("formnum"))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Form malformed"} )
      return
    }
    form_data , err := tools.GetFormOfID(db , int64(form_num))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Form malformed"} )
      return
    }
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form_data.FieldJSON) , &form_construct)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a form"} )
      return
    }

    reply_num , err := strconv.Atoi(c.Param("respnum"))
    reply_data , err := returner.GetResponseByID(db , int64(reply_num))
    var reply_construct map[string]string
    err = json.Unmarshal([]byte(reply_data.ResponseJSON) , &reply_construct)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a reply"} )
      return
    }
    reply_construct["ID"] = strconv.Itoa(int(reply_data.ID))
    reply_construct["FK_ID"] = strconv.Itoa(int(reply_data.FK_ID))
    reply_construct["Identifier"] = reply_data.Identifier
    reply_construct["SubmittedAt"] = strconv.Itoa(int(reply_data.SubmittedAt))

    template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-singular-reply.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "form" : form_construct , "reply": reply_construct })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
      return
    }
    serveTwigTemplate(c , http.StatusInternalServerError , template)
  }
  // view response
}

// Handle /mod/download/FORMNAME/FORMNUMBER
func modServeDownloadForm(db *sql.DB) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_num , err := strconv.Atoi(c.Param("formnum"))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a form"} )
      return
    }
    form_name := c.Param("formname")

    err = returner.CreateInstancedCSVForGivenForm(db , int64(form_num) , globals.RootDirectory)
    if err != nil{
      fmt.Println("CSV" , err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a form"} )
      return
    }
    err = returner.CreateReadmeForGivenForm(db , int64(form_num) , globals.RootDirectory)
    if err != nil{
      fmt.Println("README" , err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a form"} )
      return
    }
    // A tar.gz file containing the CSV, as it has zipped the entire form directory together
    err = tools.CreateDownloadableForGivenForm(globals.RootDirectory , form_name )
    if err != nil{
      fmt.Println("TAR" , err)
      c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "Issue parsing a form"} )
      return
    }
    fmt.Println(err , "redirect")
    c.Redirect(http.StatusFound, "/mod/download/" + form_name + "/downloadable.tar.gz")
  }
  // "file is being generated"
  // On click, generate file then redirect into  /mod/download/FORMNAME/downloadable.tar.gz which will serve the file
}
func modDownloadableForm() gin.HandlerFunc {
  return func (c *gin.Context) {
    form_name := c.Param("formname")
    fmt.Println("ASDF")
    now := strconv.Itoa(int(time.Now().Unix()))
    c.FileAttachment(globals.RootDirectory + "/data/" + form_name +"/downloadable.tar.gz" , form_name + "-" + now + "-archive.tar.gz")
  }

}

func modServeAPIGetAll() gin.HandlerFunc {
  return func (c *gin.Context) {
    c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "API is unimplemented"} )
  }

}
func modServeAPIGetForm() gin.HandlerFunc {
  return func (c *gin.Context) {
    c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "API is unimplemented"} )
  }

}
func modServeAPIGetResponse() gin.HandlerFunc {
  return func (c *gin.Context) {
    c.AbortWithStatusJSON( http.StatusInternalServerError ,  gin.H{"Error": "API is unimplemented"} )
  }
}

// Handle route to /forms/FORMNAME/NUMBER
func userServeForm(db *sql.DB , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    // /:formname/:formnum
    form_name := c.Param("formname")
    form_num , err := strconv.Atoi(c.Param("formnum"))
    if err != nil {
      c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "URL is malformed"})
      return
    }
    form_data , err := returner.GetFormByNameAndID(db , form_name , int64(form_num))
    var rebuild_group former.FormConstruct
    err = json.Unmarshal([]byte(form_data.FieldJSON), &rebuild_group)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve source file"})
      return
    }
    template , err := templater.ReturnFilledTemplate(env , "user-views/user-form.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "form" : rebuild_group })
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
      return
    }
    fmt.Println("SERVE")
    serveTwigTemplate(c , http.StatusOK , template)
  }
  //Oct3
}

func modServeDelete() gin.HandlerFunc {
  return func (c *gin.Context) {
    c.Header("Content-Type", "text/plain")
    c.String(http.StatusNotAcceptable , "For safety purposes, you can only submit deletes through the listing\nImagine if someone gave you a link shortener that redirects to this URL and caused you to delete someone's form...")
  }
}

/* POST Handlers */

func modPostLoginForm(db *sql.DB , cfg *types.ConfigurationSettings , env *stick.Env) gin.HandlerFunc {
  return func (c *gin.Context) {
    c.Header("Content-Type", "text/html")
    json := c.PostForm("json")
    stored_pass , err := getStoredPassword(db)
    if err != nil && json == "" {
      template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-login.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "error" : "DB Error" })
      if err != nil{
        fmt.Println(err)
        c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
        return
      }
      serveTwigTemplate(c , http.StatusInternalServerError , template)
    } else if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Internal error, Get stored"} )
    }
    ip := c.ClientIP()
    param_pass := c.PostForm("password")
    err = CheckPasswordValid(  param_pass , stored_pass.HashedPassword )
    if err != nil {
      fmt.Println(err)
      if json == "" {
        template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-login.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "error" : "Invalid Password" })
        if err != nil{
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
          return
        }
        serveTwigTemplate(c , http.StatusUnauthorized , template)
      } else if err != nil {
        fmt.Println(err)
        c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Internal error, Get stored"} )
      }
    } else{
      session_key_unencrypted := "ADMIN" + param_pass + strconv.Itoa(int(time.Now().Unix()))
      session_key_safe := CreateAuthenticationHash( session_key_unencrypted )
      // Store cookie
      var login_fields types.LoginDBFields
      login_fields = CreateLoginFields( session_key_safe , ip )
      err = StoreLogin(db , login_fields)
      if err != nil {
        if json == "" {
          template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-login.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion , "error" : "Login Storage Error" })
          fmt.Println(err)
          if err != nil{
            c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Template generation failed"} )
            return
          }
          serveTwigTemplate(c , http.StatusInternalServerError , template)
        } else if err != nil {
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"Error": "Login Storage Error"} )
        }
      }
      if json == "" {
        // (name, value string, maxAge int, path, domain string, secure, httpOnly bool)
        c.SetCookie("verified" , session_key_safe , int(30 * 24 * 60 * 60) , "/" , cfg.SiteName, true , true)
        c.Redirect(http.StatusMovedPermanently, "/mod")
      } else{
        c.AbortWithStatusJSON(http.StatusOK ,  gin.H{"message": "Success"} )
      }
    }
  }
}

func userPostForm(db *sql.DB) gin.HandlerFunc {
  return func (c *gin.Context) {
    var response_map map[string]string = make(map[string]string)
    var file_map map[string]former.MultipartFile = make(map[string]former.MultipartFile)

    form_name := c.Param("formname")
    form_num , err := strconv.Atoi(c.Param("formnum"))

    //Get a form
    form, err := tools.GetFormOfID(db , int64(form_num))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue getting form"} )
    }
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form.FieldJSON), &form_construct)
    if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue unmarshalling input"} )
    }

    responder.FillMapWithPostParams(c , response_map , form_construct)
    responder.FillMapWithPostFiles(c , file_map , form_construct)
    fmt.Println(form_construct , response_map)

    response_form := former.FormResponse{
      FormName: form_name,
      RelationalID: int64(form_num),
      ResponderID: c.ClientIP(),
      Responses: response_map,
      FileObjects:  file_map,
    }

    edit_mode , old_user_name, err := responder.CheckIfEdit(db  , response_form )
    if edit_mode{
      responder.DeleteResponderFolder( globals.RootDirectory , response_form , old_user_name )
      // Deleting is important because the responder ID could be set to scramble
      // Also easier and since nothing relies on the data it can be done
      responder.DeleteDatabaseResponse(db , int64(form_num) , old_user_name )
    }

    if _ , ok := c.GetPostForm("anon-option") ; ok {
      response_form.ScrambleResponderID()
    }

    // Check
    var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(response_form.Responses , form_construct)
    var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(response_form.FileObjects , form_construct)
    issue_array := append(text_issue_array, file_issue_array...)
    if len(issue_array) != 0{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "There are mistakes with the form" , "issue-list": issue_array } )
      return
    }
    err = responder.CreateResponderFolder( globals.RootDirectory , response_form )
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue creating responder data"} )
      os.RemoveAll(globals.RootDirectory + "/data/" + response_form.FormName + "/" + response_form.ResponderID + "/")
      return
    }
    error_list := tools.WriteFilesFromMultipart(globals.RootDirectory , response_form)
    if len(error_list) != 0 {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue creating responder data"} )
      os.RemoveAll(globals.RootDirectory + "/data/" + response_form.FormName + "/" + response_form.ResponderID + "/")
      return
    }
    err = tools.WriteResponsesToJSONFile(globals.RootDirectory , response_form)
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue creating responder data"} )
      os.RemoveAll(globals.RootDirectory + "/data/" + response_form.FormName + "/" + response_form.ResponderID + "/")
      return
    }
    response_fields , err := responder.FormResponseToDBFormat(response_form)
    if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue creating responder DB data"} )
      os.RemoveAll(globals.RootDirectory + "/data/" + response_form.FormName + "/" + response_form.ResponderID + "/")
      return
    }
    // A combination of Responses and File Locations listing a URL for file download where it will be served
    err = tools.StoreResponseToDB(db , response_fields)
    if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue creating responder DB data "} )
      os.RemoveAll(globals.RootDirectory + "/data/" + response_form.FormName + "/" + response_form.ResponderID + "/")
      return
    }

    if c.PostForm("json") != "" {
      c.String(http.StatusOK , "Submitted")
    } else{
      c.JSON(http.StatusOK , gin.H{"message": "Submitted"})
    }
  }
}

func modPostCreateForm(db *sql.DB , cfg *types.ConfigurationSettings) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_construct_raw := c.PostForm("form-construct-json")
    var form_construct former.FormConstruct
    err := json.Unmarshal([]byte(form_construct_raw), &form_construct)

    um , _ := json.Marshal(form_construct)
    fmt.Println( form_construct_raw )
    fmt.Println(string(um))

    if err != nil {
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue unmarshalling input"} )
    }
    issue_list := builder.ValidateForm(db, form_construct)
    if len(issue_list) > 0 {
      fmt.Println(issue_list)
      c.AbortWithStatusJSON(http.StatusNotAcceptable ,  gin.H{"error": "Invalid inputs" , "error-list" : issue_list } )
      return
    } else{
        insertable_form , err :=  builder.MakeFormWritable(form_construct)
        if err != nil {
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with Inputs"} )
          return
        }
        last_id , err := builder.StoreForm(db, insertable_form)
        if err != nil{
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with DB writing"} )
          return
        }
        ck_err := builder.CheckFormDirectoryExists(form_construct, globals.RootDirectory)
        if ck_err == nil{
          err = destroyer.UndoDirectory(form_construct , globals.RootDirectory)
          if err != nil {
            destroyer.UndoForm(db , insertable_form.Name)
            c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with Folder writing"} )
            return
          }
        }
        err = builder.CreateFormDirectory(form_construct , globals.RootDirectory)
        if err != nil {
          destroyer.UndoForm(db , insertable_form.Name)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with Folder writing"} )
          return
        }
        c.JSON(http.StatusOK, gin.H{ "message" : "Form written" , "URL": "https://" + cfg.SiteName + "/public/forms/" + form_construct.StorageName() + "/" +  strconv.Itoa(last_id) })
    }
  }
  //Oct3

}
func modPostEditForm(db *sql.DB , cfg *types.ConfigurationSettings) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_num , _ := strconv.Atoi(c.Param("formnum"))

    original_form , err := tools.GetFormOfID(db , int64(form_num))
    var original_form_construct former.FormConstruct
    if err != nil{
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue reading DB"} )
    }
    err = json.Unmarshal([]byte(original_form.FieldJSON), &original_form_construct)
    if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue unmarshalling input"} )
    }

    form_construct_raw := c.PostForm("form-construct-json")
    var form_construct former.FormConstruct
    err = json.Unmarshal([]byte(form_construct_raw), &form_construct)
    if err != nil {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Issue unmarshalling input"} )
    }
    issue_list := builder.ValidateFormEdit(original_form_construct, form_construct)
    if len(issue_list) > 0 {
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusNotAcceptable ,  gin.H{"error": "Invalid inputs" , "error-list" : issue_list } )
      return
    } else{
        insertable_form , err :=  builder.MakeFormWritable(form_construct)
        if err != nil {
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with Inputs"} )
          return
        }
        err = builder.UpdateForm(db, int64(form_num) , insertable_form)
        if err != nil{
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "Server Issue with DB writing"} )
          return
        }
        c.JSON(http.StatusOK, gin.H{ "message" : "Form was altered" , "URL": "https://" + cfg.SiteName + "/public/forms/" + form_construct.StorageName() + "/" +  strconv.Itoa(form_num) })
    }
  }

}
func modPostDeleteForm(db *sql.DB) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_name := c.Param("formname")
    form_num , err := strconv.Atoi(c.Param("formnum"))
    if err != nil{
      fmt.Println(err)
      c.JSON(http.StatusInternalServerError , gin.H{"error":"Malformed URL"})
    }
    err = destroyer.DeleteForm(db , form_name, int64(form_num))
    if err != nil{
      fmt.Println(err)
      c.JSON(http.StatusInternalServerError , gin.H{"error":"Issue on delete"})
    } else{
      c.JSON(http.StatusOK , gin.H{"message":"Deleted " + c.Param("formname") + " No. " + c.Param("formnum")})
    }
  }
}
func modPostDeleteResponse(db *sql.DB) gin.HandlerFunc {
  return func (c *gin.Context) {
    form_name := c.Param("formname")
    response_number , err := strconv.Atoi(c.Param("respnum"))
    if err != nil{
      fmt.Println(err)
      c.AbortWithStatusJSON(http.StatusInternalServerError ,  gin.H{"error": "URL malformed"} )
      return
    }
    response_fields , err := tools.GetResponseByID(db , int64(response_number))
    err = destroyer.DeleteResponse(db , globals.RootDirectory , int64(response_number) , form_name , response_fields.Identifier)
    if err != nil{
      fmt.Println(err)
      c.JSON(http.StatusInternalServerError , gin.H{"error":"Issue on delete"})
    } else{
      c.JSON(http.StatusOK , gin.H{"message":"Deleted " + c.Param("formname") + " No. " + c.Param("respnum")})
    }
  }
}


/* middleware */
// return function instead of handling directly to potentially pass in command line arguments on initialization

func authenticationMiddleware(db *sql.DB , env *stick.Env) gin.HandlerFunc{
  //Oct3
  	return func (c *gin.Context) {
      cookie_verification , err := c.Cookie("verified")
      json := c.PostForm("json")
      ip := c.ClientIP()

      err = CheckCookieValid(db , cookie_verification , ip)
      if err != nil && json != "" {
        fmt.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
      } else if err != nil{
        // ReturnFilledTemplate(env *stick.Env, template_path string, value_map map[string]stick.Value) (string , error)
        template , err := templater.ReturnFilledTemplate(env , "mod-views/mod-login.twig" , map[string]stick.Value{ "version" : globals.ProjectVersion })
        if err != nil {
          fmt.Println(err)
          c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized - Login Form Tenmplate Error"})
          return
        }
        c.Header("Content-Type", "text/html")
        serveTwigTemplate(c , http.StatusUnauthorized , template)
        c.Abort()
      }
      c.Next()
  	}
}
// func JWTDecodeMiddleware() gin.HandlerFunc {
// 	return func() gin.HandlerFunc {
//     // VALIDATE JWT
//     token_string, _ := c.Cookie("freeadstoken")
//     name, is_donor, is_mod, err := bannerjwt.IsAuth(token_string)
//     c.Set("name", name)
//     c.Set("is_donor", is_donor)
//     c.Set("is_mod", is_mod)
//     c.Set("valid_jwt", err == nil)
//     c.Next()
// 	}
// }
