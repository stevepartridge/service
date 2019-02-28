package swagger

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/go-chi/chi"
	"github.com/go-openapi/spec"
)

const (
	DefaultPath = "/docs"
)

type Swagger struct {
	Title    string
	Version  string
	Schemes  []string
	JSONData []byte
	Path     string

	mux             *chi.Mux
	securitySchemes map[string]*spec.SecurityScheme
	security        []map[string][]string
}

// New initiates swagger server
func New(mux *chi.Mux) (*Swagger, error) {
	s := &Swagger{
		mux:     mux,
		Schemes: []string{"http", "https"},
		Path:    DefaultPath,
	}
	return s, nil
}

// AddSecurityScheme is a helper to add security schemes such as OAuth2, etc.
func (s *Swagger) AddSecurityScheme(name string, scheme spec.SecurityScheme) {

	if s.securitySchemes == nil {
		s.securitySchemes = map[string]*spec.SecurityScheme{}
	}

	if s.security == nil {
		s.security = []map[string][]string{}
	}

	s.security = append(s.security, map[string][]string{
		name: []string{},
	})

	s.securitySchemes[name] = &scheme

}

// Serve handles the docs, swagger.json, and server.js
func (s *Swagger) Serve() {

	s.mux.HandleFunc(fmt.Sprintf("%s/swagger.json", s.Path), func(w http.ResponseWriter, req *http.Request) {

		swag := spec.Swagger{}

		err := swag.UnmarshalJSON(s.JSONData)
		if err != nil {

			w.Write(errorJSON(err))
		}

		swag.Info.Title = s.Title
		swag.Info.Version = s.Version
		// swag.Info.Contact
		// swag.Info.License
		swag.Schemes = s.Schemes

		if s.securitySchemes != nil {
			swag.SecurityDefinitions = s.securitySchemes
		}

		if len(s.security) > 0 {

			if swag.Paths != nil {

				for k := range swag.Paths.Paths {

					if swag.Paths.Paths[k].Get != nil {
						swag.Paths.Paths[k].Get.Security = s.security
					}
					if swag.Paths.Paths[k].Post != nil {
						swag.Paths.Paths[k].Post.Security = s.security
					}
					if swag.Paths.Paths[k].Put != nil {
						swag.Paths.Paths[k].Put.Security = s.security
					}
					if swag.Paths.Paths[k].Delete != nil {
						swag.Paths.Paths[k].Delete.Security = s.security
					}
					if swag.Paths.Paths[k].Patch != nil {
						swag.Paths.Paths[k].Patch.Security = s.security
					}
				}
			}
		}

		data, err := swag.MarshalJSON()
		if err != nil {

			w.Write(errorJSON(err))
		}

		w.Write(data)
	})

	docTitle := fmt.Sprintf("%s%s", strings.ToUpper(s.Title[:1]), s.Title[1:])

	serviceJS := `
function service() {
	document.title = "` + docTitle + ` API Documentation";
};
	`

	s.mux.HandleFunc(fmt.Sprintf("%s/service.js", s.Path), func(w http.ResponseWriter, req *http.Request) {
		// fmt.Println("service.js")

		w.Write([]byte(serviceJS))
	})

	mime.AddExtensionType(".svg", "image/svg+xml")

	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger/ui",
	})

	fs := http.StripPrefix(s.Path, fileServer)

	if s.Path != "/" && s.Path[len(s.Path)-1] != '/' {
		s.mux.Get(s.Path, http.RedirectHandler(s.Path+"/", 301).ServeHTTP)
		s.Path += "/"
	}
	s.Path += "*"

	s.mux.Get(s.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

}

func errorJSON(err error) []byte {
	return []byte(`
{
  "error": "` + err.Error() + `"
}
  `)
}
