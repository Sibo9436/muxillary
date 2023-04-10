package muxillary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)
type HttpMethod string
const (
  GET HttpMethod= "GET"
  POST = "POST"
  PUT = "PUT"
  PATCH = "PATCH"
  DELETE = "DELETE"
)
type pathNode struct{
  value string
  isAny bool
  mappings map[HttpMethod]http.HandlerFunc
  children map[string]*pathNode
}

type MuxillaryHandler struct{
    root *pathNode
    basepath string
}

func NewMuxillaryHandler(basepath string) *MuxillaryHandler{
  return &MuxillaryHandler{
    root: newPathNode("/"),
    basepath: basepath,
  }
}

func newPathNode(path string) *pathNode{
  isAny := strings.HasPrefix(path, ":")
  if isAny {
    path = strings.TrimLeft(path,":")
  }
  return &pathNode{
    value : path, 
    isAny: isAny,
    mappings: make(map[HttpMethod]http.HandlerFunc),
    children: make(map[string]*pathNode),
  }
}

func (p*pathNode) setMapping(mapping HttpMethod,f http.HandlerFunc){
  p.mappings[mapping] = f

}
func (m*MuxillaryHandler) setMapping(path string) *pathNode{
  fullpath := m.basepath + path
  //TODO: Controllare che Split funzioni correttamente 
  fullpath = strings.TrimLeft(fullpath, "/")
  paths := strings.Split(fullpath, "/")
  current := m.root
  fmt.Println(fullpath)
  fmt.Printf("%+v", paths)
  for _,p := range paths{
    if _, has := current.children[p]; !has{
      fmt.Println("Inserting ", p , " into ", current.value)
      current.children[p] = newPathNode(p)
    }
    current = current.children[p]
  }
  return current
}
func (m*MuxillaryHandler) Post(path string , f http.HandlerFunc){
  current := m.setMapping(path)
  current.setMapping(POST, f)
}
func (m*MuxillaryHandler) Get(path string , f http.HandlerFunc){
  current := m.setMapping(path)
  current.setMapping(GET, f)
}

func notFound(w http.ResponseWriter){
  w.WriteHeader(http.StatusNotFound)
  m:= make(map[string]any)
  m["error"] = "The page you are looking for could not be found"
  //Ignore the error
  msg,_:= json.Marshal(m)
  w.Write(msg)
}
func (m* MuxillaryHandler) ServeHTTP(rw http.ResponseWriter,r* http.Request){
  fmt.Println("Received request at ", r.URL)
  paths := strings.Split(strings.TrimLeft(r.URL.Path,"/"), "/")
  current := m.root
  for _, path := range paths{
    c, has := current.children[path]
    //SCHIFO , RISOLVERE ASSOLUTAMENTE
    for k, v := range current.children{
      fmt.Printf("%s: %+v\n", k, v)
      if v.isAny{
        c = v
        has = true
        ctx := context.WithValue(r.Context(),"mux_"+c.value,path)
        fmt.Printf("Adding %s to context, with value %s\n",c.value, path)
        r = r.WithContext(ctx)
      }
    }
    if  !has{
      notFound(rw)
      return
    }
    current = c
  }

  mapping, found := current.mappings[HttpMethod(r.Method)]
  //TODO: assolutamente da risolvere!
  //Meglio magari controllare se http2 ha una qualche specifica per questi pathParams
  
  if !found {
    notFound(rw)
    return
  }
  mapping(rw, r)
}









