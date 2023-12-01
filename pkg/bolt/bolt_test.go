package bolt

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{
			name:    "database open ok",
			in:      "/tmp/test.db",
			wantErr: false,
		}, {
			name:    "database open not ok",
			in:      "/tmp/noexists/test.db",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%s) error = %v, wantErr %v", tt.in, err, tt.wantErr)
				return
			}
			// Test cleanup
			if db != nil {
				db.Close()
			}
			os.Remove(tt.in)
		})
	}
}

func MustOpenDB(t *testing.T) (*DB, func(t *testing.T)) {
	t.Helper()
	dir, err := os.MkdirTemp("", "bolttest")
	checkError(err, t)
	f, err := os.CreateTemp(dir, "db")
	checkError(err, t)
	db, err := Open(f.Name())
	checkError(err, t)
	fn := func(t *testing.T) {
		t.Helper()
		db.Close()
		os.RemoveAll(dir)
	}
	return db, fn
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
