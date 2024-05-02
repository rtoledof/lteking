package mongo

import "testing"

func NewTestDB(t *testing.T) *DB {
	t.Helper()
	return NewDB("mongodb://localhost:27017", "test")
}
