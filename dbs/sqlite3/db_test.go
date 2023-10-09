package sqlite3

import (
	"log"
	"strings"
	"testing"
	"time"
)

type Post struct {
	// db tag lets you specify the column name
	// if it differs from the struct field
	Id      int64 `db:"post_id"`
	Created int64
	Title   string
	Body    string
}
type MetaInfo struct {
	Id       int64
	Created  int64
	Title    string
	Tags     string
	Category string
	NewTags  []string `db:"-"`
}

func init() {
	// Must Have key `Id`
	RegisterModels = map[string]interface{}{
		"posts":      Post{},
		"meta_infos": MetaInfo{},
	}
}

func TestInitDb(t *testing.T) {

	dbPath = "../../data/dbs/data.db"

	dbMap, err := GetDbMap(dbPath)
	if err != nil {
		t.Fatalf("Failed to init db, err: %v\n", err)
	}
	defer dbMap.Db.Close()

	err = InitDb(dbPath)
	if err != nil {
		t.Fatalf("failed to init db, err: %v\n", err)
	}

	// delete any existing rows
	err = dbMap.TruncateTables()
	checkErr(err, "TruncateTables failed")
}

func TestInsertToDb(t *testing.T) {
	dbPath = "../../data/dbs/data.db"

	dbMap, err := GetDbMap("")
	if err != nil {
		t.Fatalf("Failed to init db, err: %v\n", err)
	}
	defer dbMap.Db.Close()

	err = InitDb(dbPath)
	if err != nil {
		t.Fatalf("failed to init db, err: %v\n", err)
	}

	// create two posts
	p1 := newPost("Go 1.1 released!", "Lorem ipsum lorem ipsum")
	p2 := newPost("Go 1.2 released!", "Lorem ipsum lorem ipsum")

	// insert rows - auto increment PKs will be set properly after the insert
	err = dbMap.Insert(&p1, &p2)
	if err != nil {
		t.Fatalf("Failed to insert, err: %v\n", err)
	}
}

func TestSelectFromDb(t *testing.T) {
	dbPath = "../../data/dbs/data.db"

	dbMap, err := GetDbMap("")
	if err != nil {
		t.Fatalf("Failed to init db, err: %v\n", err)
	}
	defer dbMap.Db.Close()

	err = InitDb(dbPath)
	if err != nil {
		t.Fatalf("failed to init db, err: %v\n", err)
	}

	// create two posts
	p1 := newPost("Go 1.1 released!", "Lorem ipsum lorem ipsum")
	p2 := newPost("Go 1.2 released!", "Lorem ipsum lorem ipsum")

	// insert rows - auto increment PKs will be set properly after the insert
	err = dbMap.Insert(&p1, &p2)
	if err != nil {
		t.Fatalf("Failed to insert, err: %v\n", err)
	}

	// fetch all rows
	var posts []Post
	_, err = dbMap.Select(&posts, "select * from posts order by post_id")
	checkErr(err, "Select failed")
	log.Println("All rows:")
	for x, p := range posts {
		log.Printf("    %d: %+v\n", x, p)
	}
}

func TestOperationForMetaInfo(t *testing.T) {
	dbPath = "../../data/dbs/data.db"

	dbMap, err := GetDbMap("")
	if err != nil {
		t.Fatalf("Failed to init db, err: %v\n", err)
	}
	defer dbMap.Db.Close()

	err = InitDb(dbPath)
	if err != nil {
		t.Fatalf("failed to init db, err: %v\n", err)
	}

	// delete any existing rows
	err = dbMap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	// create two posts
	p1 := newMetaInfo("Go 1.1 released!", "default/123", []string{"123", "234"})
	p2 := newMetaInfo("Go 1.2 released!", "default/234", []string{"345", "234"})

	// insert rows - auto increment PKs will be set properly after the insert
	err = dbMap.Insert(&p1, &p2)
	if err != nil {
		t.Fatalf("Failed to insert, err: %v\n", err)
	}

	// fetch all rows
	var metaInfos []MetaInfo
	_, err = dbMap.Select(&metaInfos, "select * from meta_infos order by id")
	checkErr(err, "Select failed")
	log.Println("All rows:")
	for i := range metaInfos {
		metaInfos[i].NewTags = strings.Split(metaInfos[i].Tags, ",")
		log.Printf("    %d: %+v\n", i, metaInfos[i])
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func newPost(title, body string) Post {
	return Post{
		Created: time.Now().UnixNano(),
		Title:   title,
		Body:    body,
	}
}

func newMetaInfo(title, category string, tags []string) MetaInfo {
	return MetaInfo{
		Created:  time.Now().UnixNano(),
		Title:    title,
		Category: category,
		Tags:     strings.Join(tags, ","),
	}
}
