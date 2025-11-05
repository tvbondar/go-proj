package unpacker

import (
	"errors"
	"unicode"
)

// создаем переменную, содержащую ошибку
var Err = errors.New("Invalid string")

// создаем основную функцию
func Unpack(s string) (string, error) {
	str := []rune(s)
	var s2 string
	var n int
	var backslash bool //для решения доп. задания

	for i, item := range str {
		//если строка начинается с цифры, возвращаем ошибку
		if unicode.IsDigit(item) && i == 0 {
			return "", Err
		}
		//если подряд идет два числа, возвращаем ошибку
		if unicode.IsDigit(item) && unicode.IsDigit(str[i-1]) {
			return "", Err
		}
		//если стоит две черты и ... возвращаем ошибку
		if item == '\\' && !backslash {
			return "", Err
		}
		//если есть одна черта, то
		if backslash {
			s2 += string(item)
			backslash = false
			continue
		}

		if unicode.IsDigit(item) {
			n = int(item - '0')
			if n == 0 {
				s2 = s2[:len(s2)-1]
				continue
			}

			for j := 0; j < n-1; j++ {
				s2 += string(str[i-1])
			}
			continue
		}
		s2 += string(item)
	}
	return s2, nil
}
