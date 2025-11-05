package unpacker

import (
	"errors"
	"unicode"
)

// создаем rконстанту, содержащую ошибку
var ErrInvalidString = errors.New("Invalid string")

// создаем основную функцию Unpack
// она осуществляет распаковку строки последующим правилам:
// 1. Символ, за которым следует цифра n, повторяется n раз, т. е. к нему добавляется n-1 таких же символов
// 2. Если за символом находится цифра 0, то этот символ удаляется
// 3. Символ "\" экранирует следующий символ, включая цифры и сам символ "\"
// 4. Если строка начинается с цифры, или подряд идущих неэкранированных цифр, возвращаем ErrInvalidString
// 5. "\" без следующего символа == ошибка
func Unpack(s string) (string, error) {
	//если пустая строка, возвращаем nil и саму строку
	if s == "" {
		return "", nil
	}
	//преобразуем строку в слайс рун, для упрощения работы
	RuneSlice := []rune(s)
	//проверка первого символа
	//Если строка начинается с цифры, возвращаем ErrInvalidString
	if unicode.IsDigit(RuneSlice[0]) {
		return "", ErrInvalidString
	}

	var res []rune        //срез рун, содержащий результирующую строку
	backslash := false    // флаг
	lastWasDigit := false // флаг

	//основной цикл
	for _, r := range RuneSlice {

		//Случай 1.
		// Если встретили "\" - он экранирует следующий символ; сбрасываем флаги, идем дальше
		if backslash {
			res = append(res, r)
			backslash = false
			lastWasDigit = false
			continue
		}

		//Случай 2.
		// Символ "\" экранирует себя
		if r == '\\' {
			backslash = true
			continue
		}

		//Случай 3.
		//Поиск цифр в строке
		if unicode.IsDigit(r) {
			//Если результирующий слайс пуст, он не может начинаться с цифры, следовательно, ErrInvalidString
			if len(res) == 0 {
				return "", ErrInvalidString
			}
			//Есл подряд идут две неэкранированные цифры, ErrInvalidString
			if lastWasDigit {
				return "", ErrInvalidString
			}

			//переводим ASCII-цифру в число
			count := int(r - '0')
			// Если count == 0: удаляем последний рун
			if count == 0 {
				res = res[:len(res)-1]
			} else {
				//Иначе: берем последний рун
				// и добавляем его count-1 раз в res
				last := res[len(res)-1]
				for k := 0; k < count-1; k++ {
					res = append(res, last)
				}
			}
			lastWasDigit = true
			continue
		}

		//Случай 4.
		//если обычный символ, просто добавляем его в результат
		res = append(res, r)
		lastWasDigit = false
	}

	// Если завершающий слеш без экранируемого символа - ErrInvalidString
	if backslash {
		return "", ErrInvalidString
	}

	return string(res), nil
}
