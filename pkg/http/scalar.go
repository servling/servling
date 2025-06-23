package http

import (
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
)

// Config stores httpSwagger configuration variables.
type Config struct {
	URL string
}

// ScalarHandler wraps `http.ScalarHandler` into `http.HandlerFunc`.
func ScalarHandler(specUrl string) http.HandlerFunc {
	// create a template with name
	index, _ := template.New("swagger_index.html").Parse(indexTempl)

	re := regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		matches := re.FindStringSubmatch(r.RequestURI)

		path := matches[2]

		switch filepath.Ext(path) {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}

		switch path {
		case "index.html":
			_ = index.Execute(w, Config{URL: specUrl})
		case "":
			http.Redirect(w, r, matches[1]+"/"+"index.html", http.StatusMovedPermanently)
		default:
			var err error
			r.URL, err = url.Parse(matches[2])
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

				return
			}
		}
	}
}

const indexTempl = `<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="{{.URL}}"></script>

    <!-- Optional: You can set a full configuration object like this: -->
    <script>
      var configuration = {
        theme: 'purple',
      }

      document.getElementById('api-reference').dataset.configuration =
        JSON.stringify(configuration)
    </script>

    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`
