package schema

type V interface {
	INJ() RealV
}

type RealV interface {
	Check() error
}
