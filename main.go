package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{} {
	"Title": "Personal Web",
}

type Blog struct {
	Title 		string
	Duration	string
	Post_date 	string
	Author		string
	Content 	string
}

var Blogs = []Blog{
	{
		Title: "Dumbways Web App",
		Duration: "26 days",
		Post_date: "21 Oct 2022 22:22 WIB",
		Author: "Riki Wahyudi",
		Content: `Lorem ipsum dolor sit amet consectetur, adipisicing elit. Et quia quas magni
		molestiae amet. Accusantium labore harum tempore suscipit ex saepe
		cupiditate aperiam autem facere necessitatibus, architecto repellat a ad.`,
	},
}

func main() {

	route := mux.NewRouter()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer((http.Dir("./public")))))

	//routing
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/blog-detail/{index}", blogDetail).Methods("GET")
	route.HandleFunc("/blog", form).Methods("GET")
	route.HandleFunc("/process", process).Methods("POST")
	route.HandleFunc("/delete/{index}", deleted).Methods("GET")


	fmt.Println("Server running on port 5000");
	http.ListenAndServe("localhost:5000", route)
}

	//mengatur header
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println(Blogs)

	//membuat variabel memparsing template halaman index
	var tmpl, err  = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	respData := map[string]interface{}{
		"Blogs": Blogs,
	}	

	w.WriteHeader(http.StatusOK) 
	tmpl.Execute(w, respData)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var  tmpl, err = template.ParseFiles("views/form.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" +err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/blog-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" +err.Error()))
		return
	}
	var BlogDetail = Blog{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"]) 
	for i, data := range Blogs {
		if index == i {
			BlogDetail = Blog{
				Title: data.Title,
				Duration: data.Duration,
				Post_date: data.Post_date,
				Author: data.Author,
				Content: data.Content,
			}
		}
	}

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}
	
	fmt.Println(data)

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}


func form(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" +err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)

}

func process(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title :" + r.PostForm.Get("inputTitle"))
	fmt.Println("Content :" + r.PostForm.Get("inputContent"))

	var title = r.PostForm.Get("inputTitle")
	var content = r.PostForm.Get("inputContent")
	var newBlog = Blog{
		Title: title,
		Duration: "32 days",
		Content: content,
		Author: "Riki Wahyudi",
		Post_date: time.Now().String(),
	}

	Blogs = append(Blogs, newBlog)
	fmt.Println(Blogs)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleted(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	Blogs = append(Blogs[:index], Blogs[index+1:]...)
	fmt.Println(Blogs)
	http.Redirect(w, r, "/", http.StatusFound)
}