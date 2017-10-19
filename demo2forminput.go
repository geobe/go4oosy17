package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"github.com/geobe/go4oosy17/person"
	"github.com/jinzhu/gorm"
	"strconv"
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

// allgemeine Map als Schnittstelle zu den Templates
type Viewmodel map[string]interface{}

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

func personListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In personListHandler")
	r.ParseForm()
	key, err := strconv.Atoi(r.PostFormValue("persons"))
	if err != nil {
		fmt.Fprintf(w, "Error, index \"%s\" is no number!", r.PostFormValue("persons"))
	} else {
		var p person.Person
		db.First(&p, key)
		model := Viewmodel{
			"person": &p,
			"error":  "",
			"targeturl": "/personedit",
		}
		fmt.Printf("Gefunden: %v\n", p)
		templates.ExecuteTemplate(w, "personform", model)
	}

}

func editPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In editPerson")
	p, pwok := extractPerson(r)
	model := Viewmodel{
		"person":    p,
		"error":     "",
		"targeturl": "/personedit",
	}
	var person person.Person
	db.First(&person, p.ID)
	person.Firstname = p.Firstname
	person.Lastname = p.Lastname
	if pwok {
		if p.Password != ""  {
			person.Password = p.Password
		}
		fmt.Printf("from DB: %+v\n", &person)
		fmt.Printf("from Form: %+v\n", p)
		db.Save(&person)
		http.Redirect(w,r, "/list", http.StatusSeeOther)
	} else {
		model["error"] = "Passworte stimmen nicht überein"
		templates.ExecuteTemplate(w, "personform", model)
	}
}

func parsePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In parsePerson")
	p, pwok := extractPerson(r)
	model := Viewmodel{
		"person":    p,
		"error":     "",
		"targeturl": "/persondata",
	}
	var dup int
	var err bool
	db.Model(&person.Person{}).Where("username = ?", p.Username).Count(&dup)
	if dup > 0 {
		model["error"] = "Login existiert schon "
		p.Username = ""
		err = true
	}
	if !pwok {
		model["error"] = model["error"].(string) + "Passworte stimmen nicht überein"
		err = true
	}
	if err {
		p.Password = ""
		fmt.Printf("Error %s\n", model["error"])
		templates.ExecuteTemplate(w, "personform", model)
	} else {
		db.Create(&p)
		fmt.Printf(pri, p.Firstname, p.Lastname, p.Username, p.Password, p)
		fmt.Fprintf(w, html, p.Firstname, p.Lastname, p.Username)
	}
}
func extractPerson(r *http.Request) (*person.Person, bool) {
	r.ParseForm()
	lastname := r.PostFormValue("lastname")
	firstname := r.PostFormValue("firstname")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	pwrepeat := r.PostFormValue("repeat_password")
	id, err := strconv.Atoi(r.PostFormValue("id"))
	p := person.Person{
		Lastname:  lastname,
		Firstname: firstname,
		Username:  username,
		Password:  password,
	}
	if err == nil {
		p.ID = uint(id)
	} else {
		fmt.Printf("id = %v\n", r.PostFormValue("id"))
	}
	return &p, password == pwrepeat
}

func personForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In personForm")
	model := Viewmodel{
		"person": &person.Person{},
		"targeturl": "/persondata",
	}
	if templates != nil {
		templates.ExecuteTemplate(w, "personform", model)
	}
}

func personList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In personList")
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
	mux.HandleFunc("/personlist", personListHandler)
	mux.HandleFunc("/personedit", editPerson)
	// configure and start server
	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
