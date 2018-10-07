package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

type boilerPlater interface {
	appendHandlerSnippets()
	appendTypeSnippets()
	appendMainSnippet()
	appendPackageAndImportSnippet()
}

type server struct {
	Port int `toml:"port"`
}

type page struct {
	Homepage       bool   `toml:"Homepage,omitempty"`
	Name           string `toml:"name"`
	JavascriptBool bool   `toml:"javascript"`
	CSSBool        bool   `toml:"css"`
	Description    string `toml:"description,omitempty"`
	IsPublic       bool   `toml:"PublicPage"`
}

// will be used to import toml/json/yaml fields
type cauldronConfig struct {
	Pages   []page   `toml:"page"`
	Server  server   `toml:"server"`
	Imports []string `toml:"imports"`
}

var (
	config cauldronConfig
)

/*func (pot cauldronPot) appendHandlerSnippets(name string) {
	b := strings.Builder{}
	fmt.Fprintf(&b,
		`func %[1]sHandler(response http.Responseappendr, request *http.Request, title string) {
			renderTemplate(response, "%[1]s", page)
}

`, name)
	fmt.Print(b.String())
}
*/
func main() {
	cauldronDoc, err := ioutil.ReadFile("ExampleCauldron.toml")
	doc := bytes.NewBuffer(cauldronDoc).String()
	_, err = toml.Decode(doc, &config)
	if err != nil {
		log.Println(err)
	}
	/* make a directory to house the web app
	err := os.Mkdir("Random name", 0777)

	if err != nil {
		panic(err)
	}

	err = os.Chdir("Random name")

	if err != nil {
		panic(err)
	}

	/* boilerPlate := []byte(`

		package main

		import (
			"database/sql"
			"fmt"
			"html/template"
			"io/ioutil"
			"log"
			"net/http"
			"regexp"
			"strings"
		)

		var (
			templates = template.Must(template.ParseFiles("welcome.html"))
			validPath = regexp.MustCompile("^/(welcome|login|signUp|success)/$")
		)

		// Page is the typical format of a webpage on this site.
		type Page struct {
		Title string
		}

		func renderTemplate(response http.Responseappendr, tmpl string, page *Page) {
			err := templates.ExecuteTemplate(response, tmpl+".html", page)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		func checkInternalServerError(err Error) {
			if err != nil {
				return http.Error(response, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		func %sHandler(response http.Responseappendr, request *http.Request) {
			err := templates.ExecuteTemplate(response, %s+".html", *%spage)
			checkInternalServerError(err)
		}

		func makeHandler(fn func(http.Responseappendr, *http.Request)) http.HandlerFunc {
			return func(response http.Responseappendr, request *http.Request) {
				validatedPath := validPath.FindStringSubmatch(request.URL.Path)
				if validatedPath == nil {
					http.NotFound(response, request)
					return
				}
				fn(response, request)
			}
		}

		func main() {
			http.HandleFunc("/", redirectToWelcome)
			http.HandleFunc("/welcome/", makeHandler(welcomeHandler))
			log.Fatal(http.ListenAndServe(":8080", nil))
		}
		`) */

	/* below will be the actual 'boilerPlate.' This strings.Builder type will be passed to multiple function that will concatenate
	their respective []bytes to it to fully form the boilerplate code for a web server */
	/* var boilerPlate strings.Builder
	err = ioutil.appendFile("TestFile.go", boilerPlate, 0777)

	if err != nil {
		panic(err)
	} */
	b := strings.Builder{}
	fmt.Fprint(&b,
		`
	package main

	import (
		"database/sql"
		"fmt"
		"html/template"
		"io/ioutil"
		"log"
		"net/http"
		"regexp"
		"strings"
	`)

	for _, i := range config.Imports {
		fmt.Fprintf(&b,
			`"%s"
		`, i)
	}

	fmt.Fprint(&b,
		`
	)

	var (
		templates = template.Must(template.ParseFiles(`)

	// anonymous function to return a slice of strings to pass to template.ParseFiles
	parseFiles := func() []string {
		var files []string
		for _, i := range config.Pages {
			files = append(files, i.Name)
		}
		return files
	}
	validPath := strings.Join(parseFiles(), "|")

	fmt.Fprintf(&b,
		`%s...))
		validPath = regexp.MustCompile("^/(%s)/$")
			)

		`, parseFiles(), validPath)

	for _, i := range config.Pages {
		fmt.Fprintf(&b,
			`
			//%s
			type %sPage struct {
			Title string
			}
	
		`, i.Description, i.Name)
	}

	fmt.Fprint(&b,
		`	func checkInternalServerError(err Error) {
	if err != nil {
		return http.Error(response, err.Error(), http.StatusInternalServerError)				
	}
	return
}

`)

	fmt.Fprint(&b,
		`		func makeHandler(fn func(http.Responseappendr, *http.Request)) http.HandlerFunc {
		return func(response http.Responseappendr, request *http.Request) {
			validatedPath := validPath.FindStringSubmatch(request.URL.Path)
			if validatedPath == nil {
				http.NotFound(response, request)
				return
			}
			fn(response, request)
		}
	}`)

	for _, i := range config.Pages {
		fmt.Fprintf(&b,
			`func %[1]sHandler(response http.Responseappendr, request *http.Request) {
			err := templates.ExecuteTemplate(response, %[1]s+".html", *%[1]spage)
			checkInternalServerError(err)
		}

		`, i.Name)
	}

	for _, i := range config.Pages {
		if i.Homepage != false {
			fmt.Fprintf(&b,
				`func main() {
				http.HandleFunc("/", %[1]sHandler)
				http.HandleFunc("/%[1]s/", makeHandler(%[1]sHandler))
				`, i.Name)
		}
	}

	for _, i := range config.Pages {
		fmt.Fprintf(&b,
			`
		http.HandleFunc("/%[1]s/", %[1]sHandler)
		`, i.Name)
	}

	fmt.Fprint(&b, `
	}`)
	fmt.Print(b.String())

}
