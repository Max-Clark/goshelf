package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetCliPrompt(prompt string) (*string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	ret := strings.TrimSpace(in)

	return &ret, nil
}
