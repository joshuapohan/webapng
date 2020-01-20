package main

import(
	"fmt"
	"net/http"
	"os"
	_ "io"
	"strings"
	"strconv"
	"bytes"
	"encoding/json"
	"encoding/hex"

	"github.com/joshuapohan/webapng/tools"
)

type resAPNG struct{
	Status int
	Image string
}

/******************************************************
	uploadFile

	returns blob of apng
*******************************************************/
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

	err := apngEnc.Encode()
	if err != nil {
		fmt.Println(err);
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	apngEnc.WriteBytes(w)
}

/******************************************************
	JSONGenerateAPNG

	returns json containing Image (hexstring format)
*******************************************************/
func JSONGenerateAPNG(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(10 << 20)

	apngEnc := &tools.APNGModel{}
	for key, _ := range r.MultipartForm.File{
		if file, _, err := r.FormFile(key); err == nil{
			apngEnc.AppendImage(file)
		}
	}

	delays := []int{}
	for _, formValue := range r.PostForm{
		delay, _ := strconv.Atoi(formValue[0])
		delays = append(delays, delay)
	}

	delayLen := len(delays)
	for i := 1; i <= delayLen; i++{
		apngEnc.AppendDelay(delays[delayLen - i])
	}

	apngEnc.Encode()

	res := &resAPNG{}
	res.Status = 0
	buf := &bytes.Buffer{}
	apngEnc.WriteBytes(buf)
	res.Image = "0x" + hex.EncodeToString(buf.Bytes())
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

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
	//mux.HandleFunc("/upload", uploadFile)
	mux.HandleFunc("/upload", JSONGenerateAPNG)
	mux.HandleFunc("/", rootPage)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./client/build/static"))))
	http.ListenAndServe(GetPort(), mux)
}


