package storage

import "testing"

func TestLevelDBOpenClose(t *testing.T) {
	db, err := Open(t.TempDir())
	if err != nil { t.Fatal(err) }
	if err := db.Put("k", []byte("v")); err != nil { t.Fatal(err) }
	got, err := db.Get("k")
	if err != nil { t.Fatal(err) }
	if string(got) != "v" { t.Fatalf("got %q", got) }
	if err := db.Close(); err != nil { t.Fatal(err) }
}
