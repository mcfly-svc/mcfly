package client

import "fmt"

func NewBodyDecodeError(body string) error {
	return fmt.Errorf("Unable to decode Body `%s`", body)
}
