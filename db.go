package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var con *sql.DB

func init() {
	dbPath := os.Getenv("DB_LOCATION")
	if dbPath == "" {
		flag.StringVar(&dbPath, "db", "./crowd.db", "database path (including filename)")
	}
	log.Printf("Using db path %s\n", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	con = db

	doMigrations(con)
}

func FindCategory(category string) int {
	var catogoryId int
	stmt, err := con.Prepare("select id from category where category = ?")
	defer stmt.Close()
	if !dbError("couldn't create query to check for categories", err) {
		result, err := stmt.Query(category)
		defer result.Close()
		if !dbError("couldn't run query to check for categories", err) {
			if result.Next() {
				err := result.Scan(&catogoryId)
				if !dbError("problem reading column value", err) {
					return catogoryId
				}
			}
		}
	}
	return -1
}

func InsertCategory(category string) int {
	log.Printf("Adding category %s", category)
	stmt, err := con.Prepare("insert or ignore into category (category) values(?)")
	defer stmt.Close()
	if !dbError("couldn't prepare call for insert", err) {
		_, err := stmt.Exec(category)
		if !dbError("insert failed", err) {
			return FindCategory(category)
		}
	}
	return -1
}

func InsertContent(categoryId int, keyword string, detail string) bool {
	log.Printf("Adding content %s %s", keyword, detail)
	stmt, err := con.Prepare("insert into content (category_id, keyword, detail, status) values(?, ?,?, 'ACTIVE')")
	defer stmt.Close()
	if !dbError("couldn't prepare call for upsert", err) {
		_, err := stmt.Exec(categoryId, keyword, detail)
		if !dbError("insert failed", err) {
			return true
		}
	}
	return false
}

func countContent(category int, searchString string) int {
	var cnt int
	stmt, err := con.Prepare("select count(1) from content where category_id = ? and keyword like ? and status = 'ACTIVE'")
	defer stmt.Close()
	if !dbError("couldn't create query to check for details", err) {
		result, err := stmt.Query(category, searchString+"%")
		defer result.Close()
		if !dbError("couldn't run query to check for details", err) {
			if result.Next() {
				err := result.Scan(&cnt)
				if !dbError("problem reading detail value", err) {
					return cnt
				}
			}
		}
	}
	return 0
}

func FindContent(category int, searchString string) []string {
	cnt := countContent(category, searchString)
	log.Printf("count is %d", cnt)
	if cnt > 0 {
		var detail []string
		detail = make([]string, cnt, cnt)

		var tmp string
		stmt, err := con.Prepare("select detail from content where category_id = ? and keyword like ? and status = 'ACTIVE'")
		defer stmt.Close()
		if !dbError("couldn't create query to check for details", err) {
			result, err := stmt.Query(category, searchString+"%")
			defer result.Close()
			if !dbError("couldn't run query to check for details", err) {
				cnt = 0
				for result.Next() {
					err := result.Scan(&tmp)
					if !dbError("problem reading detail value", err) {
						log.Printf("found %s", tmp)
						detail[cnt] = tmp
						cnt++
					}
				}
			}
		}
		return detail
	}
	return make([]string, 0)
}

var migrations = []struct {
	ver int
	sql string
}{
	{1, "create table migrations (ver int null); insert into migrations values(0);"},
	{2, "create table category (id INTEGER PRIMARY KEY AUTOINCREMENT, category varchar(80) not null);"},
	{3, "create index idx_category on category(category);"},
	{4, "insert into category (category) values (\"mission\");"},
	{5, "insert into category (category) values (\"rss\");"},
	{6, "create table content (id INTEGER PRIMARY KEY AUTOINCREMENT, category_id INTEGER(5), status varchar(10) default 'PENDING', keyword varchar(80) not null, detail varchar(1000) not null, active integer(1) default 1);"},
	{7, "create index idx_content on content(category_id, status, keyword);"},
	{8, "insert into content(category_id, keyword, detail) values (2, 'ORE', 'Kepler Grade 3');"},
	{9, "create unique index idx_category_type on category(category);"},
}

func doMigrations(con *sql.DB) {
	var ver int
	row := con.QueryRow("select ver from migrations")
	row.Scan(&ver)

	for _, migration := range migrations {
		if migration.ver > ver {
			log.Printf("SQL: %s", migration)
			_, error := con.Exec(migration.sql)
			if !dbError("Couldn't run migration", error) {
				con.Exec("update migrations set ver = ?", migration.ver)
				ver = migration.ver
			}
		}
	}
}

func dbError(msg string, err error) bool {
	if err != nil {
		log.Fatalf("%s %s\n", msg, err)
		panic(err)
		return true
	}
	return false
}
