package schema

type V interface {
	INJ() RealV
}

type RealV interface {
	Check() error
}

/*
	V 和 realV 成对使用

	详见 phone_number.go 和 phone_number_test.go
*/
