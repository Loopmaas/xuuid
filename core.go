package xuuid

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// ================================================================
// UUID
// ================================================================
type UUID uuid.UUID

func New() UUID {
	return UUID(uuid.New())
}

func Parse(s string) (UUID, error) {
	u, error := uuid.Parse(s)

	return UUID(u), error
}

func (xu UUID) IsZero() bool {
	return (uuid.UUID)(xu) == uuid.Nil
}

func (xu *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data[:], &s); err != nil {
		return err
	} else if s == "" {
		return nil
	} else if u, err := uuid.Parse(s); err != nil {
		return err
	} else {
		*xu = UUID(u)
		return nil
	}
}

func (xu *UUID) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(xu).UnmarshalBinary(data)
}

func (xu UUID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(xu).MarshalBinary()
}

func (xu *UUID) UnmarshalText(data []byte) error {
	parsedUUID, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*xu = UUID(parsedUUID)
	return nil
}

func (xu UUID) MarshalText() ([]byte, error) {
	return uuid.UUID(xu).MarshalText()
}

func (xu UUID) String() string {
	return uuid.UUID(xu).String()
}

func (xu *UUID) Scan(src interface{}) error {
	return (*uuid.UUID)(xu).Scan(src)
}

func (xu UUID) Value() (driver.Value, error) {
	return xu.MarshalBinary()
}

func (xu UUID) ToBase62() string {
	noDashUuidStr := strings.ReplaceAll(xu.String(), "-", "")
	num := new(big.Int)
	num.SetString(noDashUuidStr, 16)
	return base62Encode(num)
}

// ================================================================
// Wildcard
// ================================================================
type Wildcard []byte

func (w *Wildcard) UnmarshalJSON(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		// string
		*w = data
		return nil
	}

	if b, err := u.MarshalBinary(); err != nil {
		return err
	} else {
		// uuid
		*w = b
		return nil
	}
}

func (w *Wildcard) UnmarshalBinary(data []byte) error {
	u, err := uuid.FromBytes(data)
	if err != nil {
		// string
		*w = data
		return nil
	}

	if b, err := u.MarshalBinary(); err != nil {
		return err
	} else {
		// uuid
		*w = b
		return nil
	}
}

func (w Wildcard) MarshalBinary() ([]byte, error) {
	if u, err := uuid.FromBytes(w); err != nil {
		// string
		return w, nil
	} else {
		// uuid
		return u.MarshalBinary()
	}
}

func (w *Wildcard) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		// string
		*w = data
		return nil
	}

	if b, err := u.MarshalBinary(); err != nil {
		return err
	} else {
		// uuid
		*w = b
		return nil
	}
}

func (w Wildcard) MarshalText() ([]byte, error) {
	if u, err := uuid.FromBytes(w); err != nil {
		// string
		return w, nil
	} else {
		// uuid
		return u.MarshalText()
	}
}

func (w Wildcard) String() string {
	if u, err := uuid.FromBytes(w); err != nil {
		// string
		return string(w)
	} else {
		// uuid
		return u.String()
	}
}

func (w Wildcard) Value() (driver.Value, error) {
	if u, err := uuid.FromBytes(w); err != nil {
		// string
		return string(w), nil
	} else {
		// uuid
		return UUID(u).Value()
	}
}

// ================================================================
var Nil UUID

// ----------------------------------------------------------------
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(num *big.Int) string {
	if num.Sign() == 0 {
		return "0"
	}

	var encoded string
	base := big.NewInt(62)
	zero := big.NewInt(0)

	for num.Cmp(zero) > 0 {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		encoded = string(base62Chars[mod.Int64()]) + encoded
	}

	return encoded
}
