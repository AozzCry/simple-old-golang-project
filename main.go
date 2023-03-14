package main

import (
	"encoding/json"
	"html/template"
	"math"
	"net/http"
	"path"
	"strconv"
)

var aparaty = []Aparat{
	{
		Nazwa:      "Canon",
		Megapixels: 10.5,
		Matryca:    19.5,
		Zoom:       20.0,
		Grade:      0,
	},
	{
		Nazwa:      "Canon1",
		Megapixels: 12.5,
		Matryca:    19.2,
		Zoom:       15.0,
		Grade:      2,
	},
	{
		Nazwa:      "Canon2",
		Megapixels: 5.5,
		Matryca:    190.5,
		Zoom:       200.0,
		Grade:      3,
	},
	{
		Nazwa:      "Canon4",
		Megapixels: 38.5,
		Matryca:    20.0,
		Zoom:       19.0,
		Grade:      4,
	},
	{
		Nazwa:      "Canon5",
		Megapixels: 0.5,
		Matryca:    10.5,
		Zoom:       1.0,
		Grade:      5,
	},
}

type Aparat struct {
	Nazwa      string
	Megapixels float64
	Matryca    float64
	Zoom       float64
	Grade      int
}

func diffrence(apA *Aparat, apB *Aparat) float64 {
	return (math.Abs(apA.Megapixels-apB.Megapixels) / 50) + (math.Abs(apA.Matryca-apB.Matryca) / 10) + (math.Abs(apA.Zoom-apB.Zoom) / 10)
}

func ShowCameras(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("pages/cameras.html")
	tmpl.Execute(w, &aparaty)
}

type fuzzy struct {
	x, o float64
}

func (f *fuzzy) calc(x float64) float64 {
	return math.Exp((x - f.x) * (x - f.x) / f.o)
}

var f1, f2, f3 = fuzzy{15, 8}, fuzzy{3, 6}, fuzzy{10, 10}

func (a *Aparat) calcAparat() float64 {
	return f1.calc(a.Megapixels) * f2.calc(a.Matryca) * f3.calc(a.Zoom)
}

func ShowOneCamera(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err == nil && id >= 0 && id < len(aparaty) {
		aparat := aparaty[id]
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(aparat)
	}
}
func main() {
	bestAparat := 0
	for i := 1; i < len(aparaty); i++ {
		if aparaty[bestAparat].calcAparat() < aparaty[i].calcAparat() {
			bestAparat = i
		}
	}
	println(aparaty[bestAparat].Nazwa)

	http.HandleFunc("/show", ShowCameras)
	http.HandleFunc("/", ShowOneCamera)
	http.ListenAndServe(":8080", nil)
}
