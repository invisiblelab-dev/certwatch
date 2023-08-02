package notifications

import (
	"fmt"
	"os"
)

type StdoutNotifier struct{}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (s *StdoutNotifier) Notify(_ string, message MessageData, _ ...string) error {
	fmt.Fprintln(os.Stdout, message)
	return nil
}

func (s *StdoutNotifier) Recipient() string {
	return ""
}
