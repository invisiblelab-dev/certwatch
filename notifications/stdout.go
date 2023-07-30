package notifications

import "fmt"

type StdoutNotifier struct{}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (s *StdoutNotifier) Notify(_ string, message string, _ ...string) error {
	fmt.Println(message)
	return nil
}

func (s *StdoutNotifier) Recipient() string {
	return ""
}
