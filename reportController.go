package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

/*
Report contains information about a given report.
The reporter must
*/
type Report struct {
	ID            int       `json:"-"`
	ReporterID    int       `json:"reporter"`
	Date          time.Time `json:"created"`
	Long          string    `json:"long"`
	Lat           string    `json:"lat"`
	Description   string    `json:"description"`
	LocationInfo  string    `json:"locInfo"`
	ImageLocation string    `json:"image"`
	Active        int       `json:"-"`
}

/*
Comment contains information of a comment made on a report.
*/
type Comment struct {
	ID       int       `json:"-"`
	ReportID int       `json:"ReportId"`
	AuthorID int       `json:"AuthorId"`
	Date     time.Time `json:"created"`
	Message  string    `json:"Message"`
	Active   int       `json:"-"`
}

/*
ReportIndex retireves all of the reports for a given system.
*/
func ReportIndex(w http.ResponseWriter, r *http.Request) {
	reports, err := getAllReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

/*
UserReports Retrieves all of the reports made by a particular user.
*/
func UserReports(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reports, err := getReportsByUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(reports); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

/*
ReportDetails Handler function to get the details for a specific report
*/
func ReportDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["reportId"]

	id, err := strconv.ParseInt(reportID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := getSpecificReport(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

/*
ReportCreate handler function for the creation of a report.
*/
func ReportCreate(w http.ResponseWriter, r *http.Request) {
	var report Report
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.Unmarshal(body, &report); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	created, err := insertReport(&report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

/*
ReportComments handler function to get the comments for a given report
*/
func ReportComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["reportId"]

	id, err := strconv.ParseInt(reportID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := getReportComments(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

/*
CommentCreate creates a comment for a specified report
*/
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := json.Unmarshal(body, &comment); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	created, err := insertComment(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

/*
GetSpecificReportComment is the handler function for getting a comment off of a report
*/
func GetSpecificReportComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["commentId"]

	id, err := strconv.ParseInt(reportID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := getCommentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

/*
DeactivateReport is the handler function for deactivating a report.
*/
func DeactivateReport(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s := v["reportId"]
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := getReportByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if &u == nil || u.Active == -1 {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}
	u, err = deactivateReportByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
DeactivateComment is the handler function for deactivating a report.
*/
func DeactivateComment(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s := v["reportId"]
	c := v["commentId"]
	reportID, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	commentID, err := strconv.ParseInt(c, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := getCommentByID(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if &u == nil || u.Active == -1 {
		http.Error(w, "Comment
		 not found", http.StatusNotFound)
		return
	}
	u, err = deactivateCommentByID(reportID, commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
