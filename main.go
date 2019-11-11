package main

import(
	"fmt"
	"net/http"
	"os"
	_ "io"
	"strconv"

	"github.com/joshuapohan/webapng/tools"
)


func uploadFile(w http.ResponseWriter, r *http.Request){
	// parse multi part form
	// 10 << 20 limits the file to 10MB
	r.ParseMultipartForm(10 << 20)

	apngEnc := &tools.APNGModel{}
	for key, values := range r.MultipartForm.File{
		fmt.Println(key, values)
		fmt.Println(values)
		
		if file, header, err := r.FormFile(key); err == nil{
			// multipartForm file implements reader , can pass directly to apngenc
			fmt.Println("Appended file of size : ", header.Size)
			apngEnc.AppendImage(file)
		}
	}

	var delays []int
	for _, formValue := range r.PostForm {
		delay, _ := strconv.Atoi(formValue[0])
		delays = append(delays, delay)
	}
	delayLen := len(delays)
	for i := 1; i <= delayLen; i++ {
		fmt.Println("Delay", delays[delayLen - i])
		apngEnc.AppendDelay(delays[delayLen - i])	
	}

	err := apngEnc.Encode()
	if err != nil {
		fmt.Println(err);
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	apngEnc.WriteBytes(w)
}

func GetPort() string {
 	var port = os.Getenv("PORT")
 	// Set a default port if there is nothing in the environment
 	if port == "" {
 		port = "8090"
 		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
 	}
 	return ":" + port
}

func rootPage(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./client/build/index.html")
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/", rootPage)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/build/static"))))
	http.ListenAndServe(GetPort(), mux)
}


