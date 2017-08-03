CREATE DATABASE IF NOT EXISTS commcomm CHARACTER SET utf8 COLLATE utf8_unicode_ci;

GRANT ALL PRIVILEGES ON commcomm.* TO 'commadmin'@'localhost' IDENTIFIED BY 'CommComm20!6';

CREATE TABLE IF NOT EXISTS commcomm.users (id BIGINT(20) UNSIGNED AUTO_INCREMENT NOT NULL, username varchar(255) NOT NULL, password varchar(255) NOT NULL, created_date DATETIME NOT NULL, active int NOT NULL,UNIQUE(id), UNIQUE(username), PRIMARY KEY(id));

CREATE TABLE IF NOT EXISTS commcomm.reports (id BIGINT(20) UNSIGNED AUTO_INCREMENT NOT NULL, reporter_id BIGINT(20) NOT NULL, report_date DATETIME NOT NULL, longitude decimal(10,6) NOT NULL, latitude decimal(10,6) NOT NULL, description varchar(255) NOT NULL, location_info varchar(255), image_location varchar(255), active int NOT NULL, UNIQUE(id), PRIMARY KEY(id));

CREATE TABLE IF NOT EXISTS commcomm.comments (id BIGINT(20) UNSIGNED AUTO_INCREMENT NOT NULL, report_id BIGINT(20) UNSIGNED NOT NULL, author_id BIGINT(20) UNSIGNED NOT NULL, comment_date DATETIME NOT NULL, message varchar(255) NOT NULL, active int NOT NULL, UNIQUE(id), PRIMARY KEY(id), FOREIGN KEY(report_id) REFERENCES commcomm.reports(id));