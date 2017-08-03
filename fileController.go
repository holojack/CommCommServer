package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

/*
GetImage Returns the image corresponding to the
TODO: possibly susceptible to attack since the passed in string is
not checked as of yet.
*/
func GetImage(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s := v["imageId"]
	http.ServeFile(w, r, "/home/ec2-user/images/"+s)
}

/*
UploadFile processes and saves an uploaded image.
*/
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	n, err := splitString(handler.Filename, "/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := os.Create("/home/ec2-user/images/" + n[len(n)-1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
