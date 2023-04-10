package muxillary_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mux "github.com/Sibo9436/muxillary"
)

func TestMux(t *testing.T){
  req := httptest.NewRequest(http.MethodGet,"/test/test/testing/123",nil)
  res := httptest.NewRecorder()
  mux:= mux.NewMuxillaryHandler("/test")
  mux.Get("/test/testing/:id",func(w http.ResponseWriter, r* http.Request){
    w.Write([]byte("CAIAOAADFAOSD"))
  })
  mux.ServeHTTP(res,req)
  if res.Result().StatusCode != http.StatusOK{
    t.Errorf("Fail: Expected status code 200, got %d", res.Result().StatusCode)
  }
  body := make([]byte, 13)
  _, err := res.Result().Body.Read(body)
  if err != nil{
    t.Error("Could not read response body")
  }
  defer res.Result().Body.Close()
  if string(body) != "CAIAOAADFAOSD"{
    t.Error("Failure: Expected body to be CAIAOAADFAOSD, got ",string(body))
  }
}
func TestMuxPathParam(t *testing.T){
  req := httptest.NewRequest(http.MethodGet,"/test/123",nil)
  res := httptest.NewRecorder()
  mux := mux.NewMuxillaryHandler("/")
  mux.Get("/test/:id",func(w http.ResponseWriter, r *http.Request){
    fmt.Printf("%+v\n",r.Context())
    id  := r.Context().Value("mux_id").(string)
    w.Write([]byte(id))
  })
  mux.ServeHTTP(res,req)
  if res.Result().StatusCode != http.StatusOK{
    t.Error("Fail: expected status 200, got: ", res.Result().StatusCode)
  }

}

func TestMux_404(t *testing.T){
  req := httptest.NewRequest(http.MethodGet,"/test/test/testing/123",nil)
  res := httptest.NewRecorder()
  mux:= mux.NewMuxillaryHandler("/test")
  mux.Post("/test/testing/:id",func(w http.ResponseWriter, r* http.Request){
    w.Write([]byte("CAIAOAADFAOSD"))
  })
  mux.ServeHTTP(res,req)
  if res.Result().StatusCode != http.StatusNotFound{
    t.Errorf("Fail: Expected status code 404, got %d", res.Result().StatusCode)
  }
}
