package bolt

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"time"

	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

var (
	userBucket        = []byte("users")
	userBucketByEmail = []byte("user_email")
)

type DB struct {
	*bolt.DB
}

const initialMmapSize = 10 * 1 << 30

func Open(path string) (*DB, error) {
	opts := &bolt.Options{
		Timeout:         1 * time.Second,
		InitialMmapSize: initialMmapSize,
	}
	db, err := bolt.Open(path, 0600, opts)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func nextID() []byte {
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return id[:]
}

func nextKeys() ([]byte, []byte) {
	private, err := ecdh.P521().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	public := private.PublicKey()
	priv := sha256.Sum256(private.Bytes())
	pub := sha256.Sum256(public.Bytes())
	return priv[:], pub[:]
}
