package unpacker

import "testing"

func TestUnpacker(t *testing.T){
	tests := []struct{
		packed string
		unpacked string
		unpackedErr bool
	}
	{
		//в качестве тестов берем примеры из текста задания
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\4\5", "qwe45", false},
		{"qwe\45", "qwe44444", false},
		{"\\", "", true}, 

	}

	for _, test := range tests{
		got, err := Unpack(test.in)
		if (err != nil) != test.unpackedErr{
			t.Fatalf("Несоответствие статуса ошибки при распаковке строки %q: ожидаемая ошибка = %v,  полученная ошибка = %v", t.packed, t.unpackedErr, err)	
		}
		if got != t.unpacked{
			t.Fatalf("Ошибка при распаковке строки %q. Ожидаемый результат: %q. Полученный результат: %q", t.packed, unpacked, got)
		}	
	}
}