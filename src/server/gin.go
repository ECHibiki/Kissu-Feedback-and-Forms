package main

import (
	"database/sql"
	"encoding/json"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/builder"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/destroyer"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/responder"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/former/returner"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"github.com/gin-gonic/gin"
	"github.com/tyler-sommer/stick"

	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Sytax for project:
  // order of items is as follows: C , DB  , CFG , Stick
  // Form search is name then number

func serveTwigTemplate(c *gin.Context, status int, template string) {
	c.Header("Content-Type", "text/html")
	c.String(status, template)
}

// bundle args into a struct
func routeGin(cfg *types.ConfigurationSettings, db *sql.DB, stick *stick.Env) *gin.Engine {

	// use args flags to set
	var gin_mode string
	if len(os.Args) > 1 && os.Args[1] == "--release" {
		gin_mode = "release"
	} else {
		gin_mode = "debug"
	}
	gin.SetMode(gin_mode)

	gin_engine := gin.Default()
	gin_engine.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	{
		gin_engine.Static("/assets", globals.RootDirectory+"public") //

		gin_engine.GET("/", generalGetHomepageHandler(stick)) //
		gin_engine.POST("/", modPostLoginForm(db, cfg, stick))

		public_group := gin_engine.Group("/public")
		{
			public_group.GET("/", generalGetHomepageHandler(stick)) //
			// Handle form requests and build forms
			public_group.GET("/forms/:formname/:formnum", userServeForm(db, stick)) //
			public_group.POST("/forms/:formname/:formnum", userPostForm(db))        //

		}

		// Verify authentication down this route
		mod_group := gin_engine.Group("/mod")
		mod_group.Use(authenticationMiddleware(db, stick))
		{
			// list menu CREATE/VIEW
			mod_group.GET("/", modServeHomepageHandler(stick)) //
			// build a form
			mod_group.GET("/create", modServeCreateForm(stick))   //
			mod_group.POST("/create", modPostCreateForm(db, cfg)) //
			// edit a form
			mod_group.GET("/edit/:formnum", modServeEditForm(db, stick))
			mod_group.POST("/edit/:formnum", modPostEditForm(db, cfg))
			// delete forms
			mod_group.POST("/form/delete/:formname/:formnum", modPostDeleteForm(db))
			mod_group.POST("response/delete/:formname/:respnum", modPostDeleteResponse(db))

			mod_group.GET("/form/delete/:formname/:formnum", modServeDelete())
			mod_group.GET("response/delete/:formname/:respnum", modServeDelete())
			// view all forms
			mod_group.GET("/view/", modServeViewAllForms(db, stick))
			// view a form with responses
			mod_group.GET("/view/:formnum", modServeViewSingleForm(db, stick)) // 6
			// view a response
			mod_group.GET("/view/:formnum/:respnum", modServeViewSingleResponse(db, stick))
			// download everything of a form
			mod_group.GET("/download/:formname/:formnum", modServeDownloadForm(db))
			mod_group.GET("/download/:formname/downloadable.tar.gz", modDownloadableForm())

			mod_group.GET("/files/:formname/:id/:filename", modDownloadableFile())
			// Retrieve various forms
			api_group := mod_group.Group("api/")
			{
				api_group.GET("/form/:formnum", modServeAPIGetForm(db , stick))
			}

		}
	}
	return gin_engine
}

func runGin(gin_engine *gin.Engine, port string) {
	gin_engine.Run(":" + port)
}

/* GET Handlers */
// Perform functions using the /former/ packages

func placeholderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Placeholder ran\n")
		c.String(http.StatusOK, "Placeholder")
	}
}

func generalGetHomepageHandler(env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		template, err := templater.ReturnFilledTemplate(env, "user-views/user-home.twig", map[string]stick.Value{"version": globals.ProjectVersion})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Home generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
}

// Handle route to /mod
func modServeHomepageHandler(env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-home.twig", map[string]stick.Value{"version": globals.ProjectVersion})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Homepage generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	//Oct3
	// login
	// Page for Create and View
}

// Handle route /mod/create
func modServeCreateForm(env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-create.twig", map[string]stick.Value{"version": globals.ProjectVersion})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Creator generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	//Oct3
	// Display the form builder and JS to get it to work
}

// Handle /mod/edit/FORMNUMBER/
func modServeEditForm(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		formnum := c.Param("formnum")
		num, err := strconv.ParseInt(formnum , 10 , 64)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "URI not a formnumber"})
			return
		}
		form_data, err := returner.GetFormOfID(db, int64(num))
		if err != nil {
			fmt.Printf("%s" , err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Can't find form"})
			return
		}
		var form_construct former.FormConstruct
		err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Invalid Unmarshaling of form"})
			return
		}
		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-edit.twig", map[string]stick.Value{
			"version": globals.ProjectVersion, "id": form_data.ID, "form": form_construct, "form_str": form_data.FieldJSON,
		})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	// Display the form builder and JS to get it to work
}

// Handle route /mod/view and /mod/api/all
func modServeViewAllForms(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_data_list, err := returner.GetAllForms(db)
		if err != nil {
			fmt.Printf("%s" , err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Can't get forms"})
			return
		}
		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-form-list.twig", map[string]stick.Value{"version": globals.ProjectVersion, "form_list": form_data_list})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	// view form list
}

// Handle route /mod/view/FORMNUMBER and /mod/api/#
func modServeViewSingleForm(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
    if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Malformed request"})
      return
    }
		form_data, err := returner.GetFormOfID(db, form_num)
    if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Malformed request"})
      return
    }
		var form_construct former.FormConstruct
		err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
		form_replies, err := returner.GetRepliesToForm(db, form_num)
		var reply_list []map[string]string
		for _, r := range form_replies {
			var r_map map[string]string
			err = json.Unmarshal([]byte(r.ResponseJSON), &r_map)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a reply"})
				return
			}
			r_map["ID"] = strconv.Itoa(int(r.ID))
			r_map["FK_ID"] = strconv.Itoa(int(r.FK_ID))
			r_map["Identifier"] = r.Identifier
			r_map["SubmittedAt"] = strconv.Itoa(int(r.SubmittedAt))
			reply_list = append(reply_list, r_map)
		}
		if err != nil {
			fmt.Printf("%s" , err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Can't get form replies"})
			return
		}
		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-reply-list.twig", map[string]stick.Value{
			"version": globals.ProjectVersion, "form": form_construct, "formnum": form_data.ID, "storagename": form_data.Name, "replies": reply_list,
		})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
}

// Handle route /mod/view/FORMNUMBER/RESPONSENUMBER and /mod/api/#/#
func modServeViewSingleResponse(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Form malformed"})
			return
		}
		form_data, err := returner.GetFormOfID(db, form_num)
		if err != nil {
			fmt.Printf("%s" , err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Couldn't find form"})
			return
		}
		var form_construct former.FormConstruct
		err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a form"})
			return
		}

		reply_num, err := strconv.ParseInt(c.Param("respnum") , 10 , 64)
    if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Form malformed"})
      return
    }
		reply_data, err := returner.GetResponseByID(db, int64(reply_num))
		var reply_construct map[string]string
		err = json.Unmarshal([]byte(reply_data.ResponseJSON), &reply_construct)
    if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Could not find ID"})
      return
    }
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a reply"})
			return
		}
		reply_construct["ID"] = strconv.Itoa(int(reply_data.ID))
		reply_construct["FK_ID"] = strconv.Itoa(int(reply_data.FK_ID))
		reply_construct["Identifier"] = reply_data.Identifier
		reply_construct["SubmittedAt"] = strconv.Itoa(int(reply_data.SubmittedAt))

		template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-singular-reply.twig", map[string]stick.Value{"version": globals.ProjectVersion, "storagename": form_data.Name, "form": form_construct, "reply": reply_construct})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	// view response
}

// Handle /mod/download/FORMNAME/FORMNUMBER
func modServeDownloadForm(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a form"})
			return
		}
		form_name := c.Param("formname")

		err = returner.CreateInstancedCSVForGivenForm(db, form_num, globals.RootDirectory)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a form"})
			return
		}
		err = returner.CreateReadmeForGivenForm(db, form_num, globals.RootDirectory)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a form"})
			return
		}
		// A tar.gz file containing the CSV, as it has zipped the entire form directory together
		err = returner.CreateDownloadableForGivenForm( form_name , globals.RootDirectory)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a form"})
			return
		}
		c.Redirect(http.StatusFound, "/mod/download/"+form_name+"/downloadable.tar.gz")
	}
	// "file is being generated"
	// On click, generate file then redirect into  /mod/download/FORMNAME/downloadable.tar.gz which will serve the file
}
func modDownloadableForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		form_name := c.Param("formname")
		now := strconv.Itoa(int(time.Now().Unix()))
		c.FileAttachment(globals.RootDirectory+"/data/"+form_name+"/downloadable.tar.gz", form_name+"-"+now+"-archive.tar.gz")
	}
}

func modDownloadableFile() gin.HandlerFunc {
	return func(c *gin.Context){
		form_name := c.Param("formname")
		user_name := c.Param("id")
		file_name := c.Param("filename")

		c.FileAttachment(globals.RootDirectory+"/data/"+form_name+"/" + user_name + "/files/" + file_name , file_name)
	}
}

func modServeAPIGetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "API is unimplemented"})
	}

}
func modServeAPIGetForm(db *sql.DB,  env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
		if err != nil {
			fmt.Printf("%s" , err.Error())
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Failed to get data"})
				return
			}
		}

		form_replies, err := returner.GetRepliesToForm(db, form_num)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Malformed request"})
			return
		}
		form_data, err := returner.GetFormOfID(db, form_num)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Malformed request"})
			return
		}

		var form_construct former.FormConstruct
		var parse_list []struct{
			Body string
			Name string
			ID int
		}
		err = json.Unmarshal([]byte(form_data.FieldJSON), &form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing form"})
			return
		}

		for _, r := range form_replies {
			var r_map map[string]string
			err = json.Unmarshal([]byte(r.ResponseJSON), &r_map)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing a reply"})
				return
			}
			r_map["ID"] = strconv.Itoa(int(r.ID))
			r_map["FK_ID"] = strconv.Itoa(int(r.FK_ID))
			r_map["Identifier"] = r.Identifier
			r_map["SubmittedAt"] = strconv.Itoa(int(r.SubmittedAt))

			template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-reply-body.twig", map[string]stick.Value{
				"storagename": form_construct.StorageName(), "reply": r_map , "form": form_construct,
			})
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue parsing template"})
				return
			}
			parse_list = append(parse_list,  struct{
				Body string
				Name string
				ID int}{ Body: template , Name: form_construct.FormName , ID: int(r.ID) } )
		}
		c.JSON(http.StatusOK, gin.H{ "formatted_replies": parse_list})
	}

}
func modServeAPIGetResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "API is unimplemented"})
	}
}

// Handle route to /forms/FORMNAME/NUMBER
func userServeForm(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// /:formname/:formnum
		form_name := c.Param("formname")
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
		if err != nil {
      tools.AbortWithJSONError( c , http.StatusBadRequest , err.Error() , gin.H{"error": "URL is malformed"})
			return
		}
		form_data, err := returner.GetFormByNameAndID(db, form_name, form_num)
		var rebuild_group former.FormConstruct
		err = json.Unmarshal([]byte(form_data.FieldJSON), &rebuild_group)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Could not retrieve source file"})
			return
		}
		fmt.Println(rebuild_group)
		template, err := templater.ReturnFilledTemplate(env, "user-views/user-form.twig", map[string]stick.Value{"version": globals.ProjectVersion, "form": rebuild_group})
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
			return
		}
		serveTwigTemplate(c, http.StatusOK, template)
	}
	//Oct3
}

func modServeDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "You can only submit deletes through the listing."})
	}
}

/* POST Handlers */

func modPostLoginForm(db *sql.DB, cfg *types.ConfigurationSettings, env *stick.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		json := c.PostForm("json")
		stored_pass, err := getStoredPassword(db)
		if err != nil && json == "" {
			template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-login.twig", map[string]stick.Value{"version": globals.ProjectVersion, "error": "DB Error"})
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
				return
			}
			serveTwigTemplate(c, http.StatusInternalServerError, template)
		} else if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Internal error, Get stored"})
		}
		ip := c.ClientIP()
		param_pass := c.PostForm("password")
		err = CheckPasswordValid(param_pass, stored_pass.HashedPassword)
		if err != nil {
			fmt.Printf("%s" , err.Error())
			if json == "" {
				template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-login.twig", map[string]stick.Value{"version": globals.ProjectVersion, "error": "Invalid Password"})
				if err != nil {
          tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
					return
				}
				serveTwigTemplate(c, http.StatusUnauthorized, template)
			} else if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Internal error, Get stored"})
			}
		} else {
			session_key_unencrypted := "ADMIN" + param_pass + strconv.Itoa(int(time.Now().Unix()))
			session_key_safe := CreateAuthenticationHash(session_key_unencrypted)
			// Store cookie
			var login_fields types.LoginDBFields
			login_fields = CreateLoginFields(session_key_safe, ip)
			err = StoreLogin(db, login_fields)
			if err != nil {
				if json == "" {
					template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-login.twig", map[string]stick.Value{"version": globals.ProjectVersion, "error": "Login Storage Error"})
					if err != nil {
            tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Template generation failed"})
						return
					}
					serveTwigTemplate(c, http.StatusInternalServerError, template)
				} else if err != nil {
					tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Login Storage Error"})
				}
			}
			if json == "" {
				// (name, value string, maxAge int, path, domain string, secure, httpOnly bool)
				c.SetCookie("verified", session_key_safe, int(30*24*60*60), "/", cfg.SiteName, true, true)
				c.Redirect(http.StatusMovedPermanently, "/mod")
			} else {
        c.JSON(http.StatusOK , gin.H{"message": "Success"})
			}
		}
	}
}

func userPostForm(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response_map map[string]string = make(map[string]string)
		var file_map map[string]former.MultipartFile = make(map[string]former.MultipartFile)

		form_name := c.Param("formname")
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
    if err != nil{
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "URL malformed"})
      return
    }
		//Get a form
		form, err := returner.GetFormOfID(db, form_num)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue getting form"})
      return
		}
		var form_construct former.FormConstruct
		err = json.Unmarshal([]byte(form.FieldJSON), &form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue unmarshalling input"})
      return
		}

		responder.FillMapWithPostParams(c, response_map, form_construct)
		responder.FillMapWithPostFiles(c, file_map, form_construct)

		response_form := former.FormResponse{
			FormName:     form_name,
			RelationalID: form_num,
			ResponderID:  c.ClientIP(),
			Responses:    response_map,
			FileObjects:  file_map,
		}

		// Check
		var text_issue_array []former.FailureObject = responder.ValidateTextResponsesAgainstForm(response_form.Responses, form_construct)
		var file_issue_array []former.FailureObject = responder.ValidateFileObjectsAgainstForm(response_form.FileObjects, form_construct)
		issue_array := append(text_issue_array, file_issue_array...)
		if len(issue_array) != 0 {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , "Validation error" , gin.H{"error": "There are mistakes with the form", "error-list": issue_array})
			return
		}

		edit_mode, old_user_name, err := responder.CheckIfEdit(db, response_form)
		if edit_mode {
			// Deleting is important because the responder ID could be set to scramble
			// Also easier and since nothing relies on the data it can be done
			destroyer.UndoResponse(db, response_form, old_user_name, globals.RootDirectory)
		}

		if _, ok := c.GetPostForm("anon-option"); ok {
			response_form.ScrambleResponderID()
		}

		err = responder.CreateResponderFolder(globals.RootDirectory, response_form)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue creating responder data"})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}
		error_list := tools.WriteFilesFromMultipart(globals.RootDirectory, response_form)
		if len(error_list) != 0 {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue creating responder data"})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}
		err = responder.WriteResponsesToJSONFile(globals.RootDirectory, response_form)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue creating responder data"})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}
		response_fields, err := responder.FormResponseToDBFormat(response_form)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue creating responder DB data"})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}
		if len(response_fields.ResponseJSON) > 65000 {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , "len(response_fields.FieldJSON) > 65000 !" , gin.H{"error": "There is too much data in your form!"})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}
		// A combination of Responses and File Locations listing a URL for file download where it will be served
		err = responder.StoreResponseToDB(db, response_fields)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue creating responder DB data "})
			if !edit_mode {
				destroyer.UndoResponse(db, response_form, response_form.ResponderID, globals.RootDirectory)
			}
			return
		}

		if c.PostForm("json") == "" {
			c.String(http.StatusOK, "Submitted")
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Submitted"})
		}
	}
}

func modPostCreateForm(db *sql.DB, cfg *types.ConfigurationSettings) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_construct_raw := c.PostForm("form-construct-json")
		var form_construct former.FormConstruct
		err := json.Unmarshal([]byte(form_construct_raw), &form_construct)
		if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue unmarshalling input. Did you fill out everything?"})
      return
		}
		issue_list := builder.ValidateForm(db, form_construct)
		if len(issue_list) > 0 {
      tools.AbortWithJSONError( c , http.StatusNotAcceptable , fmt.Sprintf("%+v\n", issue_list) , gin.H{"error": "Invalid inputs", "error-list": issue_list})
      return
		} else {
			insertable_form, err := builder.MakeFormWritable(form_construct)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with Inputs"})
				return
			}
			last_id, err := builder.StoreForm(db, insertable_form)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with DB writing"})
				return
			}
			ck_err := builder.CheckFormDirectoryExists(form_construct, globals.RootDirectory)
			if ck_err == nil {
				err = destroyer.UndoFormDirectory(form_construct, globals.RootDirectory)
				if err != nil {
					destroyer.UndoForm(db, insertable_form.Name, globals.RootDirectory)
          tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with Folder writing"})
					return
				}
			}
			err = builder.CreateFormDirectory(form_construct, globals.RootDirectory)
			if err != nil {
				destroyer.UndoForm(db, insertable_form.Name, globals.RootDirectory)
        tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with Folder writing"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Form written", "URL": "https://" + cfg.SiteName + "/public/forms/" + form_construct.StorageName() + "/" + strconv.Itoa(last_id)})
		}
	}
	//Oct3

}
func modPostEditForm(db *sql.DB, cfg *types.ConfigurationSettings) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_num, _ := strconv.ParseInt(c.Param("formnum") , 10 , 64)

		original_form, err := returner.GetFormOfID(db, form_num)
		var original_form_construct former.FormConstruct
		if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue reading DB"})
      return
		}
		err = json.Unmarshal([]byte(original_form.FieldJSON), &original_form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue unmarshalling input"})
      return
		}

		form_construct_raw := c.PostForm("form-construct-json")
		var form_construct former.FormConstruct
		err = json.Unmarshal([]byte(form_construct_raw), &form_construct)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Issue unmarshalling input"})
      return
		}
		issue_list := builder.ValidateFormEdit(form_construct, original_form_construct)
		if len(issue_list) > 0 {
			tools.AbortWithJSONError( c , http.StatusNotAcceptable , fmt.Sprintf("%+v\n", issue_list) , gin.H{"error": "Invalid inputs", "error-list": issue_list})
			return
		} else {
			insertable_form, err := builder.MakeFormWritable(form_construct)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with Inputs"})
				return
			}
			err = builder.UpdateForm(db, form_num, insertable_form)
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Server Issue with DB writing"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Form was altered", "URL": "https://" + cfg.SiteName + "/public/forms/" + form_construct.StorageName() + "/" + strconv.Itoa(int(form_num))})
		}
	}

}
func modPostDeleteForm(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_name := c.Param("formname")
		form_num, err := strconv.ParseInt(c.Param("formnum") , 10 , 64)
		if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "Malformed URL"})
      return
		}
		err = destroyer.DeleteForm(db, form_name, form_num)
		if err != nil {
      tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() ,  gin.H{"error": "Issue on delete"})
      return
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Deleted " + c.Param("formname") + " No. " + c.Param("formnum")})
		}
	}
}
func modPostDeleteResponse(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		form_name := c.Param("formname")
		response_number, err := strconv.ParseInt(c.Param("respnum") , 10 , 64)
		if err != nil {
			tools.AbortWithJSONError( c , http.StatusInternalServerError , err.Error() , gin.H{"error": "URL malformed"})
			return
		}
		response_fields, err := returner.GetResponseByID(db, int64(response_number))
		err = destroyer.DeleteResponse(db, globals.RootDirectory, int64(response_number), form_name, response_fields.Identifier)
		if err != nil {
			fmt.Printf("%s" , err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Issue on delete"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Deleted " + c.Param("formname") + " No. " + c.Param("respnum")})
		}
	}
}

/* middleware */
// return function instead of handling directly to potentially pass in command line arguments on initialization

func authenticationMiddleware(db *sql.DB, env *stick.Env) gin.HandlerFunc {
	//Oct3
	return func(c *gin.Context) {
		cookie_verification, err := c.Cookie("verified")
		json := c.PostForm("json")
		ip := c.ClientIP()

		err = CheckCookieValid(db, cookie_verification, ip)
		if err != nil && json != "" {
			tools.AbortWithJSONError( c , http.StatusUnauthorized , err.Error() , gin.H{"error": "Unauthorized"})
		} else if err != nil {
			// ReturnFilledTemplate(env *stick.Env, template_path string, value_map map[string]stick.Value) (string , error)
			template, err := templater.ReturnFilledTemplate(env, "mod-views/mod-login.twig", map[string]stick.Value{"version": globals.ProjectVersion})
			if err != nil {
				tools.AbortWithJSONError( c , http.StatusUnauthorized , err.Error() , gin.H{"error": "Unauthorized - Login Form Tenmplate Error"})
				return
			}
			c.Header("Content-Type", "text/html")
			serveTwigTemplate(c, http.StatusUnauthorized, template)
			c.Abort()
		}
		c.Next()
	}
}
