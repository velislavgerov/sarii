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

func sayHelloSari(w http.ResponseWriter, r *http.Request) {
    adjs := readFileLines("internal/db/adjectives.txt")
    verbs := readFileLines("internal/db/verbs.txt")
    emjs := getEmojisFromFile("internal/db/emojis.txt")

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

    t := template.Must(template.ParseFiles("web/template/index.tmpl.html"))
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
        chores := readFileLines("internal/db/chores.txt")
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
    fs := http.FileServer(http.Dir("assets/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    if err := http.ListenAndServe(":" + port, nil); err != nil {
        log.Fatal(err)
    }
}