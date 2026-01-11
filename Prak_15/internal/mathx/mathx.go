package mathx

import "errors"

func Sum(a, b int) int { return a + b }

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

// MustDivide делит a на b и panic-ит при делении на 0.
func MustDivide(a, b int) int {
	if b == 0 {
		panic("divide by zero")
	}
	return a / b
}
