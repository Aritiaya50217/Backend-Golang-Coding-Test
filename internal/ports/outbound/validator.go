package outbound

type Validator interface {
	Validate(i interface{}) error
}
