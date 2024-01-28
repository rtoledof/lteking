package mongo

func NewTestDB() *DB {
	return NewDB("mongodb://localhost:27017", "test")
}
