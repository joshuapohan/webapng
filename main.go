package main

import(
	"fmt"
	"net/http"
	_ "io"

	"github.com/joshuapohan/apng"
)


func uploadFile(w http.ResponseWriter, r *http.Request){
	// 10 << 20 limits the file to 10MB
	r.ParseMultipartForm(10 << 20)

	apngEnc := &apng.APNGModel{}
	fmt.Println(r)
	for key, values := range r.MultipartForm.File{
		fmt.Println(key, values)
		fmt.Println(values)
		
		if file, header, err := r.FormFile(key); err == nil{
			// multipartForm file implements reader , can pass directly to apngenc
			fmt.Println("Appended file of size : ", header.Size)
			apngEnc.AppendImage(file)
			apngEnc.AppendDelay(64)

			// example of reading to byte slice first before passing
			/*
			fileBuffer := make([]byte, header.Size)
			file.Read(fileBuffer)
			fmt.Println(fileBuffer)
			*/

			// example of reading nb of bytes (4 bytes each chunk)
			/*
			fileBytes := []byte{}
			readBuffer := make([]byte, 4)
			readSize := 0
			for{
				curSize, err := file.Read(readBuffer)
				fmt.Println("Read chunk of, ", curSize, " bytes")
				fileBytes = append(fileBytes, readBuffer...)
				if err == io.EOF{
					break
				}
			}

			fmt.Println("Read ", readSize)
			fmt.Println("Total Size ", header.Size)
			fmt.Println(fileBytes)
			*/	
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	apngEnc.Encode()
	apngEnc.WriteBytes(w)
	apngEnc.SavePNGData("result.png")
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8090", mux)
}


