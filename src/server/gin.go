package main

import (
  "github.com/gin-gonic/gin"
  "github.com/tyler-sommer/stick"
  "database/sql"
  "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"

  "os"
)

func PlaceholderHandler(c *gin.Context){ }

// bundle args into a struct
func routeGin(cfg *types.ConfigurationSettings, db *sql.DB , stick *stick.Env ) *gin.Engine{

 // use args flags to set
 var gin_mode string
 if os.Args[1] == "--release" {
   gin_mode = "release"
 } else {
   gin_mode = "debug"
 }
 gin.SetMode( gin_mode )

 gin_engine := gin.Default()
 gin_engine.SetTrustedProxies([]string{"127.0.0.1"})

 {
   gin_engine.GET("/",  PlaceholderHandler)

   // Verify authentication down this route
   mod_group := gin_engine.Group("mod/")
   mod_group.Use(authenticationMiddleware())
   {
     mod_group.GET("/", PlaceholderHandler)
     // Retrieve various forms
     api_group := mod_group.Group("api/")
     {
       // API calls to return details for given forms
       api_group.GET("/", PlaceholderHandler)

     }

   }

   display_group := gin_engine.Group("display/")
   {
     display_group.GET("/", PlaceholderHandler)
     // Handle form requests and build forms

   }
 }
 return gin_engine
}

func runGin(gin_engine *gin.Engine, cfg *types.ConfigurationSettings ){
  gin_engine.Run(cfg.StartupPort)
}

/* middleware */
// return function instead of handling directly to potentially pass in command line arguments on initialization

func authenticationMiddleware() gin.HandlerFunc{
  	return func(c *gin.Context) {
      c.Next()
  	}
}
// func JWTDecodeMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
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
