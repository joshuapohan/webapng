package main

import(
	"fmt"
	"net/http"
	_ "io"
	"strings"
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
	for formKey, formValue := range r.PostForm {
		if strings.Index(formKey, "input") > - 1 {
			delay, _ := strconv.Atoi(formValue[0])
			delays = append(delays, delay)
		}
	}
	delayLen := len(delays)
	for i := 1; i <= delayLen; i++ {
		fmt.Println("Delay", delays[delayLen - i])
		apngEnc.AppendDelay(delays[delayLen - i])	
	}

	apngEnc.Encode()
	w.Header().Set("Access-Control-Allow-Origin", "*")	
	apngEnc.WriteBytes(w)
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8090", mux)
}


