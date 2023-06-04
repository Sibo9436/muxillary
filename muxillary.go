package muxillary

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST              = "POST"
	PUT               = "PUT"
	PATCH             = "PATCH"
	DELETE            = "DELETE"
)

type pathNode struct {
	value    string
	isAny    bool
	mappings map[HttpMethod]http.HandlerFunc
	children map[string]*pathNode
}

type MuxillaryHandler struct {
	root     *pathNode
	basepath string
	notFound func(http.ResponseWriter)
}

func NewMuxillaryHandler(basepath string) *MuxillaryHandler {
	return &MuxillaryHandler{
		root:     newPathNode("/"),
		basepath: basepath,
		notFound: notFound,
	}
}

func newPathNode(path string) *pathNode {
	isAny := strings.HasPrefix(path, ":")
	if isAny {
		path = strings.TrimLeft(path, ":")
	}
	return &pathNode{
		value:    path,
		isAny:    isAny,
		mappings: make(map[HttpMethod]http.HandlerFunc),
		children: make(map[string]*pathNode),
	}
}
func (p *pathNode) addPathVar(pathVar string) {

	pathVar = strings.TrimLeft(pathVar, ":")
	if len(p.value) > 0 {
		p.value = p.value + "/" + pathVar
	} else {
		p.value = pathVar
	}
}

func (p *pathNode) setMapping(mapping HttpMethod, f http.HandlerFunc) {
	p.mappings[mapping] = f

}
func (m *MuxillaryHandler) setMapping(path string) *pathNode {
	fullpath := m.basepath + path
	//TODO: Controllare che Split funzioni correttamente
	fullpath = strings.TrimLeft(fullpath, "/")
	paths := strings.Split(fullpath, "/")
	current := m.root
	for _, p := range paths {
		isAny := strings.HasPrefix(p, ":")
		path = p
		if isAny {
			p = "*"
		}
		if _, has := current.children[p]; !has {
			current.children[p] = newPathNode(path)
		} else {
			current.children[p].addPathVar(path)
		}
		current = current.children[p]
	}
	return current
}

func (m *MuxillaryHandler) Delete(path string, f http.HandlerFunc) {
	current := m.setMapping(path)
	current.setMapping(DELETE, f)
}
func (m *MuxillaryHandler) Patch(path string, f http.HandlerFunc) {
	current := m.setMapping(path)
	current.setMapping(PATCH, f)
}
func (m *MuxillaryHandler) Put(path string, f http.HandlerFunc) {
	current := m.setMapping(path)
	current.setMapping(PUT, f)
}
func (m *MuxillaryHandler) Post(path string, f http.HandlerFunc) {
	current := m.setMapping(path)
	current.setMapping(POST, f)
}
func (m *MuxillaryHandler) Get(path string, f http.HandlerFunc) {
	current := m.setMapping(path)
	current.setMapping(GET, f)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	m := make(map[string]any)
	m["error"] = "The page you are looking for could not be found"
	//Ignore the error
	msg, _ := json.Marshal(m)
	w.Write(msg)
}

func Value(key string, r *http.Request) string {
	//qui si potrebbe poi vedere anche se parsare altre cose in altri modi
	res := r.Context().Value("mux_" + key)
	if res != nil {
		return res.(string)
	}
	//TODO: decidere se implementare un sistema di errori un po' più smart
	return ""

}

// func print(node *pathNode) {
// fmt.Println(node.value, node.isAny)
// for _, child := range node.children {
// print(child)
// }
// }
func (m *MuxillaryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//Ricordarsi di estrarre questa funzione in modo da
	//poter creare una forma di componibilità
	//Voglio poter definire più Router e SubRouter
	////print(m.root)
	//fmt.Println("Received request at ", r.URL)
	path := strings.Split(r.URL.Path, "?")[0]
	paths := strings.Split(strings.TrimLeft(path, "/"), "/")
	current := m.root
	for _, path := range paths {
		//fmt.Println("Checking path", path)
		c, found := current.children[path]
		if !found {
			c, found = current.children["*"]
			//fmt.Println("Checking if has any", path)
			//fmt.Printf("%+v", current.children)
			if found {
				for _, p := range strings.Split(c.value, "/") {
					ctx := context.WithValue(r.Context(), "mux_"+p, path)
					r = r.WithContext(ctx)
				}
			}
		}
		if !found {
			m.notFound(rw)
			return
		}
		current = c
	}

	mapping, found := current.mappings[HttpMethod(r.Method)]
	//TODO: assolutamente da risolvere!
	//Meglio magari controllare se http2 ha una qualche specifica per questi pathParams

	if !found {
		m.notFound(rw)
		return
	}
	mapping(rw, r)
}
