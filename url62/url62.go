package url62

import (
	//"github.com/nu7hatch/gouuid"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const base62alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func FromUUID(shex string) (string, error) {
	u := strings.Replace(shex, "-", "", 4)

	converted, err := convertUp(u, base62alphabet)
	if err != nil {
		return "", err
	}

	return converted, nil
}

func ToUUID(s string) (string, error) {
	//.insert(8, '-').insert(13, '-').insert(18, '-').insert(23, '-')
	u, err := convertDown(s, base62alphabet)
	if err != nil {
		return "", err
	}

	u = insert(u, "-", 8)
	u = insert(u, "-", 13)
	u = insert(u, "-", 18)
	u = insert(u, "-", 23)

	return u, nil
}

// ConvertUp converts a hexadecimal UUID string to a base alphabet greater than
// 16. It is used here to compress a 32 character UUID down to 23 URL friendly
// characters.
func convertUp(oldNumber string, baseAlphabet string) (string, error) {
	n := big.NewInt(0)
	n.SetString(oldNumber, 16)

	base := big.NewInt(int64(len(baseAlphabet)))

	newNumber := make([]byte, 23) //converted size of max base-62 uuid
	i := len(newNumber)

	zero := big.NewInt(0)
	for n.Cmp(zero) != 0 {
		i--
		_, r := n.DivMod(n, base, big.NewInt(0))
		newNumber[i] = baseAlphabet[r.Int64()]
	}
	return fmt.Sprintf("%022s", newNumber[i:]), nil
}

func convertDown(oldNumber string, baseAlphabet string) (string, error) {
	// Return if oldNumber is less or bigger than 22
	if len(oldNumber) != 22 {
		return "", errors.New("Not valid base62-encoded UUID string")
	}

	n := big.NewInt(0)

	s := strings.Split(oldNumber, "")

	var i int

	for _, v := range s {
		i = strings.Index(base62alphabet, v)

		if i == -1 {
			return "", errors.New("Not valid base62 string. String does not contain a right Base62 alphabet")
		}

		n.Mul(n, big.NewInt(62))
		n.Add(n, big.NewInt(int64(i)))
	}

	return fmt.Sprintf("%032x", n), nil
}

func insert(s string, sep string, i int) string {
	return s[0:i] + sep + s[i:len(s)]
}
