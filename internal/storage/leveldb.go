package storage

import "github.com/syndtr/goleveldb/leveldb"

type LevelDB struct {
	db *leveldb.DB
}

func Open(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{db: db}, nil
}

func (l *LevelDB) Put(key string, value []byte) error {
	return l.db.Put([]byte(key), value, nil)
}

func (l *LevelDB) Get(key string) ([]byte, error) {
	return l.db.Get([]byte(key), nil)
}

func (l *LevelDB) Write(batch *leveldb.Batch) error {
	return l.db.Write(batch, nil)
}

func (l *LevelDB) Close() error {
	if l == nil || l.db == nil {
		return nil
	}
	return l.db.Close()
}
