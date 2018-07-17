import path from 'path';
import sqlite3 from 'sqlite3';
const sqlite3Verbose = sqlite3.verbose();
const db = new sqlite3Verbose.Database(path.join(process.cwd(), './data/db.sqlite'));

const sqlCreateMemos = "CREATE TABLE IF NOT EXISTS `memos` (`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `name`	TEXT,`date`	TEXT,`shop`	TEXT,`page`	TEXT,`play`	TEXT,`talk`	TEXT);";
db.serialize(() => {
  console.log("initialize database")
  db.run(sqlCreateMemos);
});

export default db;