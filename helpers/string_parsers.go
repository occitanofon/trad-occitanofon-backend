package helpers

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmailValid checks if email has valid length, match regex and host exists
func IsEmailValid(email string) error {
	if len(email) < 3 {
		return &HelperError{
			Code:    409,
			Message: ErrEmailTooShort,
			Wrapped: fmt.Errorf("email is too short: %s", email),
		}
	}

	if len(email) > 254 {
		return &HelperError{
			Code:    409,
			Message: ErrEmailTooLong,
			Wrapped: fmt.Errorf("email is too long: %s", email),
		}
	}

	if !emailRegex.MatchString(email) {
		return &HelperError{
			Code:    409,
			Message: ErrEmailBadFormat,
			Wrapped: fmt.Errorf("email is in bad format: %s", email),
		}
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return &HelperError{
			Code:    409,
			Message: ErrEmailDomainNotExist,
			Wrapped: fmt.Errorf("domain email doesn't exist: %s", email),
		}
	}
	return nil
}

func UsernameValidity(username string) error {
	if len(username) < 3 {
		return &HelperError{
			Code:    409,
			Message: ErrUsernameTooShort,
			Wrapped: errors.New("username too short"),
		}
	}

	if err := onlyListedCaracters(username); err != nil {
		return err
	}

	if err := onlyOneSpace(username); err != nil {
		return err
	}

	return nil
}

func Normalize(str string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, str)
	if err != nil {
		return "", &HelperError{
			Code:    500,
			Message: ErrDefault,
			Wrapped: err,
		}
	}

	return result, nil
}

// onlyOneSpace checks if username has only one space
func onlyOneSpace(username string) error {
	if count := strings.Count(username, " "); count > 1 {
		return &HelperError{
			Code:    409,
			Message: ErrOnlyOneSpaceByUsername,
			Wrapped: errors.New("only one space by username"),
		}
	}
	return nil
}

// onlyListedCaracters checks if username letters match occitan alphabet
func onlyListedCaracters(username string) error {
	// lettres autorisées
	// à è ò á é í ó ú ï ü ç
	// À È Ò Á É Í Ó Ú Ï Ü Ç
	// a b c d e f g h i j l m n o p q r s t u v x z
	// A B C D E F G H I J L M N O P Q R S T U V X Z
	for _, c := range username {
		if !(c == 'à' || c == 'è' || c == 'ò' || c == 'á' || c == 'é' ||
			c == 'í' || c == 'ó' || c == 'ú' || c == 'ï' || c == 'ü' ||
			c == 'ç' ||
			c == 'À' || c == 'È' || c == 'Ò' || c == 'Á' || c == 'É' ||
			c == 'Í' || c == 'Ó' || c == 'Ú' || c == 'Ï' || c == 'Ü' ||
			c == 'Ç' ||
			c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' ||
			c == 'f' || c == 'g' || c == 'h' || c == 'i' || c == 'j' ||
			c == 'l' || c == 'm' || c == 'n' || c == 'o' || c == 'p' ||
			c == 'q' || c == 'r' || c == 's' || c == 't' || c == 'u' ||
			c == 'v' || c == 'x' || c == 'z' ||
			c == 'A' || c == 'B' || c == 'C' || c == 'D' || c == 'E' ||
			c == 'F' || c == 'G' || c == 'H' || c == 'I' || c == 'J' ||
			c == 'L' || c == 'M' || c == 'N' || c == 'O' || c == 'P' ||
			c == 'Q' || c == 'R' || c == 'S' || c == 'T' || c == 'U' ||
			c == 'V' || c == 'X' || c == 'Z' || c == ' ') {
			return &HelperError{
				Code:    409,
				Message: fmt.Sprintf("Lo caractère “%c“ es pas defendut", c),
				Wrapped: fmt.Errorf("%c is not allowed", c),
			}
		}
	}
	return nil
}
