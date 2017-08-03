package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

/*
InitDb initializes the DB for CommCommServer based on the passed in credentials.
The db object is non-exported so all db calls must be in the database.go file.
If the connection cannot be opened or the db cannot be pinged, returns an error.
*/
func InitDb(user, pass, ip, port string) error {
	var err error
	db, err = sql.Open("mysql", user+":"+pass+"@tcp("+ip+":"+port+")/commcomm?parseTime=true")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

/*
InsertUser inserts a user of the CommComm ecosystem. If the user name already exists,
the function will return an error. If no error, returns the user.
*/
func InsertUser(username, password string, created time.Time) (*User, error) {
	stmt, err := db.Prepare("INSERT users SET username=?,password=?,created_date=?,active=1")
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(username, password, created.Format(time.RFC1123))
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user, err := GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

/*
UpdateUsername updates the user name of the passed in user. If the username
is taken, the function returns an error.
*/
func UpdateUsername(username, newname string) error {
	stmt, err := db.Prepare("UPDATE users set username=? where username=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newname, username)
	if err != nil {
		return err
	}

	return nil
}

/*
UpdateUserPassword updates the password of the passed in username.
*/
func UpdateUserPassword(password, username string) error {
	stmt, err := db.Prepare("update users set password=? where username=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(password, username)
	if err != nil {
		return err
	}
	return nil
}

/*
GetUserByEmail retrieves a particular user based on their username.
*/
func GetUserByEmail(username string) (*User, error) {
	var u User

	stmt, err := db.Prepare("SELECT * FROM users where username=?")
	if err != nil {
		return &u, err
	}

	row := stmt.QueryRow(username)
	if err != nil {
		return &u, err
	}

	err = row.Scan(&u.ID, &u.Email, &u.Password, &u.Date, &u.Active)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

/*
GetUserByID gets a user by their integer ID.
*/
func GetUserByID(id int64) (*User, error) {
	var u User
	stmt, err := db.Prepare("SELECT * FROM users where id=?")
	if err != nil {
		return &u, err
	}

	row := stmt.QueryRow(id)
	if err != nil {
		return &u, err
	}

	err = row.Scan(&u.ID, &u.Email, &u.Password, &u.Date, &u.Active)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

/*
GetAllUsers Retrieves all the users of the CommComm system
*/
func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT * FROM users WHERE active=1")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Date, &u.Active); err != nil {
			return nil, err
		}
		u.Password = ""
		users = append(users, u)
	}

	return users, nil
}

func getAllReports() ([]Report, error) {
	rows, err := db.Query("SELECT * FROM reports where active=1")
	if err != nil {
		return nil, err
	}

	var reports []Report

	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.ID, &r.ReporterID, &r.Date, &r.Long, &r.Lat, &r.Description, &r.LocationInfo, &r.ImageLocation, &r.Active); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}

func getReportsByUser(id int64) ([]Report, error) {
	stmt, err := db.Prepare("SELECT * FROM reports where active=1 AND reporter_id=?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var reports []Report

	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.ID, &r.ReporterID, &r.Date, &r.Long, &r.Lat, &r.Description, &r.LocationInfo, &r.ImageLocation, &r.Active); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}

func getSpecificReport(id int64) (*Report, error) {
	stmt, err := db.Prepare("SELECT * FROM reports where active=1 AND id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)
	if err != nil {
		return nil, err
	}

	var r Report
	if err := row.Scan(&r.ID, &r.ReporterID, &r.Date, &r.Long, &r.Lat, &r.Description, &r.LocationInfo, &r.ImageLocation, &r.Active); err != nil {
		return nil, err
	}

	return &r, nil
}

func closeDb() {
	db.Close()
}

func getReportComments(id int64) ([]Comment, error) {
	stmt, err := db.Prepare("SELECT * FROM comments where active=1 AND report_id=?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	var comments []Comment

	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.ReportID, &c.AuthorID, &c.Date, &c.Message, &c.Active); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func insertReport(r *Report) (*Report, error) {
	stmt, err := db.Prepare("INSERT reports SET reporter_id=?,report_date=?,longitude=?,latitude=?,description=?,location_info=?,image_location=?,active=1")
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(r.ReporterID, time.Now().Format(time.RFC1123), r.Long, r.Lat, r.Description, r.LocationInfo, "")
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	report, err := getReportByID(id)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func insertComment(c *Comment) (*Comment, error) {
	stmt, err := db.Prepare("INSERT comments SET report_id=?,author_id=?,comment_date=?,message=?,active=1")
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(c.ReportID, c.AuthorID, time.Now().Format(time.RFC1123), c.Message)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	comment, err := getCommentByID(id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func getReportByID(id int64) (*Report, error) {
	var r Report

	stmt, err := db.Prepare("SELECT * FROM reports where id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&r.ID, &r.ReporterID, &r.Date, &r.Long, &r.Lat, &r.Description, &r.LocationInfo, &r.ImageLocation, &r.Active)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func getCommentByID(id int64) (*Comment, error) {
	var c Comment

	stmt, err := db.Prepare("SELECT * FROM comments where id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&c.ID, &c.ReportID, &c.AuthorID, &c.Date, &c.Message, &c.Active)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func deactivateUserByID(id int64) (*User, error) {
	var u User

	stmt, err := db.Prepare("SELECT * FROM users where id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&u.ID, &u.Email, &u.Password, &u.Date, &u.Active)
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE users set active=-1 where id=?")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func deactivateReportByID(id int64) (*Report, error) {
	var r Report

	stmt, err := db.Prepare("SELECT * FROM reports where id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&r.ID, &r.ReporterID, &r.Date, &r.Long, &r.Lat, &r.Description, &r.LocationInfo, &r.ImageLocation, &r.Active)
	if err != nil {
		return nil, err
	}

	stmt, err = db.Prepare("UPDATE reports set active=-1 where id=?")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func deactivateCommentByID(reportID, commentID int64) (*Comment, error) {
	var c Comment

	stmt, err := db.Prepare("SELECT * FROM comments where id=? and report_id=?")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(commentID, reportID)

	err = row.Scan(&c.ID, &c.ReportID, &c.AuthorID, &c.Date, &c.Message, &c.Active)
	if err != nil {
		return nil, err
	}

	_, err = db.Prepare("UPDATE comments set active=-1 where id=? and report_id=?")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(commentID, reportID)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
