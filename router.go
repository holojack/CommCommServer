package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
Route contains information to pass in to the mux router
*/
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var userRoutes = []Route{
	Route{
		"Get User by ID",
		"GET",
		"/user/{userId}",
		GetSpecificUserByID,
	},
	Route{
		"Get User by Email",
		"GET",
		"/user",
		GetSpecificUserByEmail,
	},
	Route{
		"UserCreate",
		"POST",
		"/user",
		UserCreate,
	},
	Route{
		"Deactivate user",
		"DELETE",
		"/user/{userId}",
		DeactivateUser,
	},
}

var reportRoutes = []Route{
	Route{
		"Get User Reports",
		"GET",
		"/user/{userId}/report",
		UserReports,
	},
	Route{
		"Get Report Details",
		"GET",
		"/report/{reportId}",
		ReportDetails,
	},
	Route{
		"ReportCreate",
		"POST",
		"/report",
		ReportCreate,
	},
	Route{
		"RportIndex",
		"GET",
		"/report",
		ReportIndex,
	},
}

var commentRoutes = []Route{
	Route{
		"Get Report Comments",
		"GET",
		"/report/{reportId}/comment",
		ReportComments,
	},
	Route{
		"Create comment",
		"POST",
		"/report/{reportId}/comment",
		CommentCreate,
	},
	Route{
		"Get specific comment",
		"GET",
		"/report/{reportId}/comment/{commentId}",
		GetSpecificReportComment,
	},
}

var otherRoutes = []Route{
	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},
	Route{
		"Store image",
		"POST",
		"/report/{reportId}/image",
		UploadFile,
	},
	Route{
		"Get image from report",
		"GET",
		"report/{reportId}/image",
		GetImage,
	},
}

/*
InitRouter initializes the mux router for the CommComm API
*/
func InitRouter() (r *mux.Router) {
	var routes = []Route{}
	routes = append(routes, userRoutes...)
	routes = append(routes, reportRoutes...)
	routes = append(routes, commentRoutes...)
	routes = append(routes, otherRoutes...)

	r = mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		r.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}

	return
}
