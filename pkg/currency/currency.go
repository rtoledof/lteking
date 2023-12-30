package currency

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/currency"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

var roundingCrypto = map[string]int{
	BTC:  8,
	BTG:  8,
	BCC:  8,
	USDT: 2,
	USDC: 2,
	EURC: 2,
	EURT: 2,
	DASH: 8, // DASH
	LTC:  8,
	LSK:  8,
	XRP:  6,
	ZEC:  8,
	DOGE: 8,
	TRX:  6,
	CUP:  2,
}

const (
	BTC  = "BTC"
	BTG  = "BTG"
	BCC  = "BCC"
	USDT = "USDT"
	USDC = "USDC"
	EURC = "EURC"
	EURT = "EURT"
	DASH = "DASH"
	LSK  = "LSK"
	LTC  = "LTC"
	XRP  = "XRP"
	ZEC  = "ZEC"
	DOGE = "DOGE"
	TRX  = "TRX"
	CUP  = "CUP"
)

type Currency struct {
	Unit currency.Unit
	unit string
}

func (c *Currency) IsCrypto() bool {
	_, ok := roundingCrypto[c.String()]
	return ok
}

func (c *Currency) Rounding() (int, int) {
	if v, ok := roundingCrypto[c.unit]; ok {
		return v, 0
	}
	return currency.Kind{}.Rounding(c.Unit)
}

func (c *Currency) Parse(cur string) (err error) {
	c.Unit = currency.XXX
	cu := strings.ToUpper(cur)
	if _, ok := roundingCrypto[cu]; ok {
		c.unit = cu
		return nil
	}
	c.Unit, err = currency.ParseISO(cur)
	return err
}

func (c Currency) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	s := c.String()
	_, err := fmt.Fprintln(&b, s)
	return b.Bytes(), err
}

func (c Currency) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Currency) Scan(src any) error {
	return c.Parse(fmt.Sprintf("%s", src))
}

func (c *Currency) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	var d string
	_, err := fmt.Fscanln(b, &d)
	if err != nil {
		return err
	}
	return c.Parse(d)
}

func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Currency) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return c.Parse(str)
}

func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Currency) UnmarshalText(b []byte) error {
	return c.Parse(string(b))
}

func (c Currency) MarshalBSON() (interface{}, error) {
	return c.String(), nil
}

func (c *Currency) UnmarshalBSON(src interface{}) error {
	return c.Parse(fmt.Sprintf("%s", src))
}

func (c *Currency) String() string {
	if c.unit != "" {
		return c.unit
	}
	if c.Unit == currency.XXX {
		return ""
	}
	return c.Unit.String()
}

func (c *Currency) Equals(cur string) bool {
	return c.String() == strings.ToUpper(cur)
}

func (c Currency) Equal(cur Currency) bool {
	return c.String() == cur.String()
}

func Parse(src string) (Currency, error) {
	var c Currency
	err := c.Parse(src)

	return c, err
}

func MustParse(cur string) Currency {
	c, err := Parse(cur)
	if err != nil {
		panic(err)
	}
	return c
}

var XXX = Currency{Unit: currency.XXX}
