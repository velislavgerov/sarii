package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"math/rand"
	"io/ioutil"
	"time"
	"text/template"
)

type SariPageData struct {
	PageTitle string
	BoxColor string
	BtnColor string
	Verb string
	Adjective string
	Emoji Emoji
	Chore Chore
}

type Emoji struct {
	Character string
	Name string
}

type Chore struct {
	Text string
	Drawn bool
	dateSelected string
}

var chore Chore = Chore{
	Text: "...",
}

var colors [10]string = [...]string{
	"navy",
	"green",
	"teal",
	"blue",
	"purple",
	"pink",
	"red",
	"orange",
	"yellow",
	"grey",
}

const templt string = `
<!DOCTYPE html>
<html>
<head>
  <title>{{.PageTitle}}</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/siimple@3.1.0/dist/siimple.min.css">
  <style>
  	#reload_btn {
  		margin-top: 1%;
  	}

  	#choose_a_chore {
  		font-size: 1.5em;
  	}
  </style>
  <script>
  	function ChooseAChore() {
  		var chooseAChoreBtn = document.getElementById("choose_a_chore_btn");
		chooseAChoreBtn.setAttribute('class', 'siimple-spinner siimple-spinner--primary');

  		var xmlhttp = new XMLHttpRequest();
		xmlhttp.onreadystatechange=function() {
		  if (xmlhttp.readyState==4 && xmlhttp.status==200) {
			chooseAChoreBtn.remove();
		    var response = xmlhttp.responseText; //if you need to do something with the returned value
		  	console.log(response);
		  	var element = document.getElementById("chore");
		  	element.innerHTML = response;
		  }
		}

		setTimeout(function() {
			xmlhttp.open("GET","/chore",true);
			xmlhttp.send();
		}, 2000)

  	}
  </script>
</head>

<body>
  <div class="siimple-jumbotron siimple-jumbotron--{{.BoxColor}} siimple-jumbotron--fluid" align="center">
    <div class="siimple-jumbotron-title">{{.Emoji.Character}}️</div>
    <div class="siimple-jumbotron-subtitle"><strong>Сари</strong> е {{.Adjective}} {{.Emoji.Name}}!</div>
    <div class="siimple-jumbotron-detail">А <strong>velislav</strong> го {{.Verb}}.</div>
    <div class="siimple-btn siimple-btn--{{.BtnColor}}" id="reload_btn" onClick="window.location.reload()">Ново</div>
  </div>
  <div class="siimple-content siimple-content--fluid" align="center">
    <div class="siimple-p" id="choose_a_chore">Днес <strong>velislav</strong> ще <span id="chore">{{.Chore.Text}}</span></div>
  	{{if not .Chore.Drawn}}
  	  <div class="siimple-btn siimple-btn--primary" id="choose_a_chore_btn" onClick="ChooseAChore()">Изтегли</div>
  	{{end}}
  </div>
  <div class="siimple-footer" align="center">
    Made with love by <strong>velislav</strong>
  </div>
</body>

</html>
`

func sayHelloSari(w http.ResponseWriter, r *http.Request) {
	adjs := readFileLines("static/adjectives.txt")
	verbs := readFileLines("static/verbs.txt")
	emjs := getEmojisFromFile("static/emojis.txt")

	date := time.Now().Local().Format("2006-01-02")
  	data := SariPageData{
  		PageTitle: "Страничката на Сари!",
  		BoxColor: colors[rand.Intn(len(colors))],
  		BtnColor: colors[rand.Intn(len(colors))],
  		Adjective: adjs[rand.Intn(len(adjs))],
  		Verb: verbs[rand.Intn(len(verbs))],
  		Emoji: emjs[rand.Intn(len(emjs))],
  		Chore: Chore{
  			Text: chore.Text,
  			Drawn: !(chore.Drawn == false ||
  				chore.dateSelected != date),
  		},
  	}

  	// no same colors for box and btn
  	for c := data.BtnColor; c == data.BoxColor; c = data.BtnColor {
  		data.BtnColor = colors[rand.Intn(len(colors))]
  	} 

  	t := template.Must(template.New("page").Parse(templt))
	err := t.Execute(w, data)
  	if err != nil {
  		log.Fatal(err)
  	}
}

func chooseAChore(w http.ResponseWriter, r *http.Request) {
	date := time.Now().Local().Format("2006-01-02")
	if chore.Drawn == false ||
		chore.dateSelected != date {
		rand.Seed(time.Now().UTC().UnixNano())
		chores := readFileLines("static/chores.txt")
		chore = Chore{
			Text: chores[rand.Intn(len(chores))],
			Drawn: true,
			dateSelected: time.Now().Local().Format("2006-01-02"),
		}
	}

	fmt.Fprintf(w, chore.Text)
}

func readFileLines(path string) []string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}

func getEmojisFromFile(path string) []Emoji {
	var emojis []Emoji
	lines := readFileLines(path)
	for _, line := range lines {
		e := strings.Split(line, " ")
		emojis = append(emojis, Emoji{
			Character: e[0],
			Name: e[1],
		})
	}
	return emojis
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", sayHelloSari)
	http.HandleFunc("/chore", chooseAChore)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}