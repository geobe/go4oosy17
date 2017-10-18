package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"github.com/geobe/go4oosy17/person"
	"github.com/jinzhu/gorm"
)

const html = `<!DOCTYPE html>
<html>
<head>
<meta content="text/html;
charset=windows-1252" http-equiv="content-type">
<title>Person Info</title>
</head>
<body>
Name: %s %s<br>
Login: %s<br>
</p>
</body>
</html>
`
const pri = `Name: %s %s
Login: %s
password: %s
Person: %+v
`

var templates *template.Template
var db *gorm.DB

// Parse die angegebenen Template-Files in ein Template
// templatedir	Verzeichnis mit den Template-Files
// fn...	ein oder mehrere Filenamen (ohne .html)
func prepareTemplates(templatedir string, fn ...string) (t *template.Template) {
	var files []string
	for _, file := range fn {
		files = append(files, fmt.Sprintf("%s/%s.html", templatedir, file))
	}
	t = template.Must(template.ParseFiles(files...))
	return
}

func parsePerson(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	lastname := r.PostFormValue("lastname")
	firstname := r.PostFormValue("firstname")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	fmt.Fprintf(w, html, firstname, lastname, username)
	p := person.Person{
		Lastname:  lastname,
		Firstname: firstname,
		Username:  username,
		Password:  password,
	}
	fmt.Printf(pri, firstname, lastname, username, password, p)
	db.Create(&p)
	//db.Create(&person.Person{
	//	Lastname:  lastname,
	//	Firstname: firstname,
	//	Username:  username,
	//	Password:  password,
	//})
}

func personForm(w http.ResponseWriter, r *http.Request) {
	//citynames := make([]string, len(poi.GermanCities))
	//for i, c := range poi.GermanCities {
	//	citynames[i] = c.Name()
	//}
	if templates != nil {
		templates.ExecuteTemplate(w, "personform", &person.Person{})
	}
}

func personList(w http.ResponseWriter, r *http.Request) {
	var perscount int
	var persons []person.Person
	db.Model(&person.Person{}).Count(&perscount)
	var perslist map[uint]string
	perslist = make(map[uint]string)
	db.Find(&persons)
	for _, p := range persons {
		perslist[p.ID] = p.Username
	}
	perslist[0] = "neuen Benutzer anlegen"
	if templates != nil {
		templates.ExecuteTemplate(w, "personlist", &perslist)
	}
}

func main() {
	// Connect to database
	db = person.Dbopen()
	// Close at program end
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&person.Person{})
	//create test data
	person.Populate(db)

	// get a multiplexer for request handling
	mux := http.NewServeMux()

	// get path to the go templates
	pwd, _ := os.Getwd()
	tpl := pwd + "/src/github.com/geobe/go4oosy17/tpl"
	// parse templates into a variable
	templates = prepareTemplates(tpl, "Personform")
	// route requests to handler functions
	mux.HandleFunc("/persondata", parsePerson)
	mux.HandleFunc("/", personForm)
	mux.HandleFunc("/list", personList)
	// configure and start server
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}