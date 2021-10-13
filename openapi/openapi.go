package openapi

import (
	"embed"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultPath is root path docs will be served up unless configured differently
	DefaultPath = "/docs"

	staticPath = "static/"
)

var (
	// ErrServeMuxMustNotBeNil error
	ErrServeMuxMustNotBeNil = errors.New("openapi: serve mux can not be nil")
	// ErrJSONDataMustNotBeNil error
	ErrJSONDataMustNotBeNil = errors.New("openapi: JSON data can not be nil")

	//go:embed static
	staticFiles embed.FS

	//go:embed static/index.html
	indexHTML []byte
	//go:embed static/oauth2-redirect.html
	redirectHTML []byte
)

// OpenAPI holds the basic config and mux
type OpenAPI struct {
	Spec    spec.Swagger
	Schemes []string
	Path    string

	mux             *http.ServeMux
	securitySchemes map[string]*spec.SecurityScheme
	security        []map[string][]string
	index           []byte
}

// New initiates openapi server
func New(mux *http.ServeMux, jsonData []byte) (*OpenAPI, error) {

	if mux == nil {
		return nil, ErrServeMuxMustNotBeNil
	}

	if jsonData == nil {
		return nil, ErrJSONDataMustNotBeNil
	}

	s := &OpenAPI{
		mux:     mux,
		Schemes: []string{"http", "https"},
		Path:    DefaultPath,
	}

	err := s.Spec.UnmarshalJSON(jsonData)
	if err != nil {
		return nil, err
	}

	s.makeIndexHTML(s.Path, indexHTML)

	return s, nil
}

// AddSecurityScheme is a helper to add security schemes such as OAuth2, etc.
func (s *OpenAPI) AddSecurityScheme(name string, scheme spec.SecurityScheme) {

	if s.securitySchemes == nil {
		s.securitySchemes = map[string]*spec.SecurityScheme{}
	}

	if s.security == nil {
		s.security = []map[string][]string{}
	}

	s.security = append(s.security, map[string][]string{
		name: {},
	})

	s.securitySchemes[name] = &scheme

}

// Serve handles the docs, openapi.json, and server.js
func (s *OpenAPI) Serve() {

	s.cleanSpec()

	s.mux.HandleFunc(fmt.Sprintf("%s/openapi.json", s.Path), func(w http.ResponseWriter, req *http.Request) {
		data, err := s.Spec.MarshalJSON()
		if err != nil {
			w.Write(errorJSON(err))
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	})

	s.mux.HandleFunc(fmt.Sprintf("%s/openapi.yaml", s.Path), func(w http.ResponseWriter, req *http.Request) {
		data, err := s.Spec.MarshalJSON()
		if err != nil {
			w.Write(errorJSON(err))
		}

		var tmp interface{}
		err = yaml.Unmarshal(data, &tmp)
		if err != nil {
			w.Write(errorJSON(err))
		}

		data, err = yaml.Marshal(tmp)
		if err != nil {
			w.Write(errorJSON(err))
		}

		w.Header().Add("Content-Type", "text/yaml")
		w.Write(data)
	})

	docTitle := s.Spec.Info.Title
	if len(s.Spec.Info.Title) > 1 {
		docTitle = fmt.Sprintf("%s%s", strings.ToUpper(s.Spec.Info.Title[:1]), s.Spec.Info.Title[1:])
	}

	serviceJS := `
function service() {
	document.title = "` + docTitle + ` API Documentation";
};
	`

	if s.Path == "/" || s.Path[len(s.Path)-1] == '/' {
		s.Path = strings.TrimRight(s.Path, "/")
		s.mux.HandleFunc(s.Path, http.RedirectHandler(s.Path, 301).ServeHTTP)
	}

	fs := http.FileServer(http.FS(staticFiles))
	fs = http.StripPrefix(s.Path+"/", fs)

	s.mux.Handle(s.Path+"/"+staticPath, fs)

	s.mux.HandleFunc(fmt.Sprintf("%s/service.js", s.Path), func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		w.WriteHeader(200)
		w.Write([]byte(serviceJS))
	})

	mime.AddExtensionType(".svg", "image/svg+xml")

	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(s.index)
	})
	redirect := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(200)
		w.Write(redirectHTML)
	})

	s.mux.Handle(s.Path, index)
	s.mux.Handle(s.Path+"/", index)
	s.mux.Handle(s.Path+"/index.html", index)
	s.mux.Handle(s.Path+"/oauth2-redirect.html", redirect)

}

func (s *OpenAPI) cleanSpec() {

	if s.securitySchemes != nil {
		s.Spec.SecurityDefinitions = s.securitySchemes
	}

	if len(s.security) > 0 {

		if s.Spec.Paths != nil {

			for k := range s.Spec.Paths.Paths {

				if s.Spec.Paths.Paths[k].Get != nil {
					s.Spec.Paths.Paths[k].Get.Security = s.security
				}
				if s.Spec.Paths.Paths[k].Post != nil {
					s.Spec.Paths.Paths[k].Post.Security = s.security
				}
				if s.Spec.Paths.Paths[k].Put != nil {
					s.Spec.Paths.Paths[k].Put.Security = s.security
				}
				if s.Spec.Paths.Paths[k].Delete != nil {
					s.Spec.Paths.Paths[k].Delete.Security = s.security
				}
				if s.Spec.Paths.Paths[k].Patch != nil {
					s.Spec.Paths.Paths[k].Patch.Security = s.security
				}
			}
		}
	}
}

func (s *OpenAPI) makeIndexHTML(path string, data []byte) {
	if path == "" {
		path = DefaultPath
	}
	if path[:1] == "/" {
		path = path[1:]
	}
	idx := string(indexHTML)
	idx = strings.ReplaceAll(idx, `href="./`, `href="./`+path+"/"+staticPath)
	idx = strings.ReplaceAll(idx, `src="./`, `src="./`+path+"/"+staticPath)

	injectJS := `
	<script src="./` + path + `/service.js"> </script>
	<script>
		window.onload = function () {
			service();

			const ui = SwaggerUIBundle({
				url: window.location.href.replace(location.hash, "") + "/openapi.json",
				oauth2RedirectUrl: window.location.href.replace(location.hash, "") + '/oauth2-redirect.html',
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				plugins: [
					SwaggerUIBundle.plugins.DownloadUrl
				],
				layout: "StandaloneLayout"
			});

			window.ui = ui;
		};
	</script>
	`

	idx = strings.ReplaceAll(idx, "<!-- {{SWAGGER JS}} -->", injectJS)

	idx = strings.ReplaceAll(idx, "\n", " ")
	idx = strings.ReplaceAll(idx, "\t", " ")
	for strings.Contains(idx, "  ") {
		idx = strings.ReplaceAll(idx, "  ", " ")
	}
	idx = strings.ReplaceAll(idx, "/> ", "/>")
	idx = strings.ReplaceAll(idx, " <", "<")
	idx = strings.ReplaceAll(idx, " >", ">")

	s.index = []byte(idx)
}

func errorJSON(err error) []byte {
	return []byte(`
{
  "error": "` + err.Error() + `"
}
  `)
}
