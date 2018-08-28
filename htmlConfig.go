package main

import (
	"html/template"
	"net/http"

	"github.com/joncalhoun/form"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/schema"
	"log"
	"encoding/json"
)

// Setup a template for how the received struct should be parsed
var inputTpl = `
<div class="form-group row">
		{{if eq .Type "checkbox"}}
			<div class="col-sm-2">{{.Name}}</div>
    		<div class="col-sm-10">
      			<div class="form-check">
					<input type="{{.Type}}" class="form-check-input" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" {{with .Value}}value="{{.}}"{{end}}>
      			</div>
    		</div>
		{{else}}
			<label {{with .ID}}for="{{.}}"{{end}} class="col-sm-2 col-form-label">{{.Label}}</label>
    		<div class="col-sm-10">
				<input type="{{.Type}}" class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
    		</div>
		{{end}}
</div>`

func main() {
	tpl := template.Must(template.New("").Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	pageTpl := template.Must(template.New("").Funcs(fb.FuncMap()).Parse(`
<html>
<head>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
</head>
<body>
	<div class="container">
		<div class="jumbotron" id="exampleModal" tabindex="-1">
			<div class="modal-dialog" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title" id="exampleModalLabel">Config</h5>
					</div>
					<div class="modal-body">
						<form action="/config" method="post">
							{{inputs_for .}}
							<button type="submit" class="btn btn-primary">Submit</button>
						</form>
					</div>
				</div>
			</div>
		</div>
	</div>
</body>
</html>
	`))

	router := httprouter.New()
	router.GET("/config", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html")
		pageTpl.Execute(w, Config{})
	})
	router.POST("/config", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var config Config
		err := dec.Decode(&config, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(config)
		w.Write(b)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Config will be used as the format for the HTML form in the view
type Config struct {
	AlarmTime string `form:"placeholder=16:30"`
	Enabled bool `form:"type=checkbox;"`
}