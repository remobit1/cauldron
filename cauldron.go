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

type cauldronPot strings.Builder

var (
	config cauldronConfig
)

func init() {
	// initializing function to make sure the environment variables are set and such and such, the environment is capable of running the program
}

func (pot cauldronPot) appendHandlerSnippets(name string) {
	b := strings.Builder{}
	fmt.Fprintf(&b,
		`func %[1]sHandler(response http.Responseappendr, request *http.Request, title string) {
			renderTemplate(response, "%[1]s", page)
}

`, name)
	fmt.Print(b.String())
}

func (pot cauldronPot) appendTypeSnippets(name string) {
	b := strings.Builder{}
	fmt.Fprintf(&b,
		`type %sPage struct {
		Title string
		}

	`, name)
}

func (pot cauldronPot) appendMainSnippet(name string) {

}

func (pot cauldronPot) appendPackageAndImportSnippet(imports []string) {
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

	for _, i := range imports {
		fmt.Fprintf(&b,
			`"%s"
		`, i)

		// change this to Fprintf later and have the parse files be generate by an array passed through channels from the appendHandler snippet. Ditto for MustCompile.
		fmt.Fprint(&b,
			`
		)

		var (
			templates = template.Must(template.ParseFiles("welcome.html"))
			validPath = regexp.MustCompile("^/(welcome|login|signUp|success)/$")
		)

	`)
	}
}

func finalizeBoilerPlate(b boilerPlater) {

}

func main() {
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

		func welcomeHandler(response http.Responseappendr, request *http.Request, title string) {
			page := &Page{Title: title}
			renderTemplate(response, "welcome", page)
		}

		func makeHandler(fn func(http.Responseappendr, *http.Request, string)) http.HandlerFunc {
			return func(response http.Responseappendr, request *http.Request) {
				validatedPath := validPath.FindStringSubmatch(request.URL.Path)
				if validatedPath == nil {
					http.NotFound(response, request)
					return
				}
				fn(response, request, validatedPath[1])
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

	cauldronDoc, err := ioutil.ReadFile("ExampleCauldron.toml")
	doc := bytes.NewBuffer(cauldronDoc).String()
	_, err = toml.Decode(doc, &config)
	if err != nil {
		log.Println(err)
	}

	pot := cauldronPot{}
	for _, i := range config.Pages {
		pot.appendHandlerSnippets(i.Name)
	}

}
