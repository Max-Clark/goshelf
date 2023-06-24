package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func GetCliPrompt(prompt string, reader io.Reader) (*string, error) {
	fmt.Print(prompt)

	rdr := bufio.NewReader(reader)
	in, err := rdr.ReadString('\n')

	if err != nil {
		return nil, err
	}

	ret := strings.TrimSpace(in)

	return &ret, nil
}
