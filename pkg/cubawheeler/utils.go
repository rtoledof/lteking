package cubawheeler

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cr "crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/oklog/ulid/v2"
)

type identifier []byte

func (id identifier) MarshalText() ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(id)))
	base64.RawURLEncoding.Encode(dst, id[:])
	return dst, nil
}

func (id *identifier) UnmarshalText(b []byte) error {
	*id = make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	_, err := base64.RawURLEncoding.Decode(*id, b)
	return err
}

func (id identifier) String() string {
	return base64.RawURLEncoding.EncodeToString(id[:])
}

// ID represents a random identifier for objects.
type ID identifier

func (id ID) MarshalText() ([]byte, error) {
	return identifier(id).MarshalText()
}

func (id *ID) UnmarshalText(b []byte) error {
	return ((*identifier)(id)).UnmarshalText(b)
}

func (id ID) String() string {
	return identifier(id).String()
}

func (id *ID) UnmarshalBson(b []byte) error {
	return id.UnmarshalText(b)
}

func (id ID) MarshalBson() ([]byte, error) {
	return id.MarshalText()
}

func NewID() ID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return ID(id[:])
}

func MustParseID(strID string) ID {
	var id ID
	id.UnmarshalText([]byte(strID))
	return id
}

type Time struct {
	time.Time
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Format(time.RFC822)), nil
}

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}
	str := fmt.Sprintf("%s", value)
	t.Time, err = time.Parse(time.RFC822, str)
	return err
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return t.Format(time.RFC822), nil
}

func Now() Time {
	return Time{time.Now().UTC()}
}

func ParseDate(date string) (Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Time{}, err
	}
	return Time{t}, nil
}

func NewKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, cr.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, nil, fmt.Errorf("unable to generate private key: %w", err)
	}
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubkey, nil
}

func NewReferalCode() string {
	var referCode = rand.Int63n(999999)
	return fmt.Sprintf("%06d", referCode)
}

func NewOtp() string {
	var otp = rand.Int63n(9999)
	return fmt.Sprintf("%04d", otp)
}

type ChannelEvent string

const (
	ChannelEventUpdateLocation ChannelEvent = "UPDATE_LOCATION"
	ChannelEventUpdateStatus   ChannelEvent = "UPDATE_STATUS"
	ChannelEventDriverArriving ChannelEvent = "DRIVER_ARRIVING"
	ChannelEventTripEnding     ChannelEvent = "TRIP_ENDING"
	ChannelEventNewOrder       ChannelEvent = "NEW_ORDER"
)

var AllChannelEvent = []ChannelEvent{
	ChannelEventUpdateLocation,
	ChannelEventUpdateStatus,
	ChannelEventDriverArriving,
	ChannelEventTripEnding,
	ChannelEventNewOrder,
}

func (e ChannelEvent) IsValid() bool {
	switch e {
	case ChannelEventUpdateLocation, ChannelEventUpdateStatus, ChannelEventDriverArriving, ChannelEventTripEnding, ChannelEventNewOrder:
		return true
	}
	return false
}

func (e ChannelEvent) String() string {
	return string(e)
}

func (e *ChannelEvent) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChannelEvent(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChannelEvent", str)
	}
	return nil
}

func (e ChannelEvent) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Response struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
}
