package main

import (
  prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	// "github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/templater"
	// "github.com/ECHibiki/Kissu-Feedback-and-Forms/globals"
	// "os"
	"fmt"
	"testing"

  "strings"
	"net/http/httptest"
	// "time"
	"net/http"
	"net/url"
)

func FuzzHelloWorld(f *testing.F){
  f.Add("")
  f.Fuzz(func(t *testing.T , str string){
    fmt.Println(str)
  })
}

func FuzzLoginSubmit(f *testing.F){
  var initialization_folder string = "./../../test"
  db, init , cfg := prebuilder.DoTestingIntializations(initialization_folder)
  prebuilder.WritePassword(db , init.ApplicationPassword, "bcrypt", "10")
  defer prebuilder.CleanupTestingInitializations(initialization_folder)
  templater.SetRootDir(initialization_folder + "/templates/")

    f.Add("base")
    f.Fuzz(func(t *testing.T , pass string){
    	gin_engine := routeGin(&cfg, db)

      fmt.Println("Started" )
      w := httptest.NewRecorder()
      data := url.Values{"password": {"ghghg"}}
      fmt.Println("before request" , data)
      req, e := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
      if e != nil {
        t.Fatal(e)
      }
      gin_engine.ServeHTTP(w, req)
      fmt.Println("Got here", w.Code )
      if w.Code != 401{
        t.Fatal("http code" , w.Code)
      }
    })
}
