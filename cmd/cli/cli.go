package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Returns all characters up to '\n' from a reader
func GetCliPrompt(prompt *string, reader io.Reader) (*string, error) {
	fmt.Print(*prompt)

	rdr := bufio.NewReader(reader)
	in, err := rdr.ReadString('\n')

	if err != nil {
		return nil, err
	}

	ret := strings.TrimSpace(in)

	return &ret, nil
}

// Returns an int from a reader. Will perform retries if user
// enters invalid character. User entering /n will return nil.
func GetIntFromCli(prompt *string, r io.Reader, w io.Writer) (*int, error) {
	for {
		valStr, err := GetCliPrompt(prompt, r)

		if err != nil {
			return nil, err
		}

		if *valStr == "" {
			return nil, nil
		}

		valInt64, err := strconv.ParseInt(*valStr, 10, 32)

		if err != nil {
			fmt.Fprintln(w, "invalid integer, try again")
			continue
		}

		id := int(valInt64)

		return &id, nil
	}
}
