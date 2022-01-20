package modules

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var Filename string

func DbInit() {

	Db, err := sql.Open("sqlite", Filename)
	if err != nil {
	}
	defer Db.Close()

	sqlCreateMemos := "CREATE TABLE IF NOT EXISTS `memos` " +
		"(`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, " +
		"`name`	TEXT,`date`	TEXT,`shop`	TEXT,`page`	TEXT,`play`	TEXT,`talk`	TEXT);"

	if _, err = Db.Exec(sqlCreateMemos); err != nil {
	}
}

func RowsToMemo(id string) ([]Memo, error) {
	var memos []Memo
	db, err := sql.Open("sqlite", Filename)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	var rows *sql.Rows
	if id == "" {
		const sql = "SELECT * FROM memos;"
		rows, err = db.Query(sql)
	} else {
		const sql = "SELECT * FROM memos WHERE id = ?"
		rows, err = db.Query(sql, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var memo Memo
		err = rows.Scan(&memo.Id, &memo.Name, &memo.Date, &memo.Shop, &memo.Page, &memo.Play, &memo.Talk)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	err = rows.Err()

	return memos, err
}

func UpsertMemo(b Memo) (err error) {
	Db, err := sql.Open("sqlite", Filename)
	if err == nil {
		if b.Id == "" {
			const sql = "INSERT INTO memos (name, date, shop, page, play, talk) VALUES (?, ?, ?, ?, ?, ?)"
			_, err = Db.Exec(sql, b.Name, b.Date, b.Shop, b.Page, b.Play, b.Talk)
		} else {
			const sql = "UPDATE memos SET name=?, date =?, shop=?, page=?, play=?, talk=? WHERE id = ?"
			_, err = Db.Exec(sql, b.Name, b.Date, b.Shop, b.Page, b.Play, b.Talk, b.Id)
		}
	}
	return err
}

func DeleteMemo(id string) error {
	Db, err := sql.Open("sqlite", Filename)
	if err == nil {
		const sql = "DELETE FROM memos WHERE id=?"
		_, err = Db.Exec(sql, id)
	}
	return err

}
