package protocol

import (
	"fmt"

	"strings"
)

// Op identifies which KV operation a Command represents.

type Op int

const (
	OpGet Op = iota

	OpSet

	OpDelete

	OpHelp

	OpUnknown
)

// Command is the parsed form of a single line of input.

type Command struct {
	Op Op

	Key string

	Value string // only populated for OpSet

}

// HelpText is returned for HELP and for the HTTP root path.

const HelpText = "Commands: GET <key> | SET <key> <value> | DELETE <key> | HELP"

// Parse reads a single line and returns the corresponding Command.

// used by both the TCP and UDP routes

func Parse(line string) (Command, error) {

	line = strings.TrimSpace(line)

	if line == "" {

		return Command{Op: OpUnknown}, fmt.Errorf("empty command")

	}

	// SplitN to 3 so a SET value containing spaces is preserved intact.

	parts := strings.SplitN(line, " ", 3)

	switch strings.ToUpper(parts[0]) {

	case "GET":

		if len(parts) < 2 {

			return Command{Op: OpUnknown}, fmt.Errorf("GET requires a key")

		}

		return Command{Op: OpGet, Key: parts[1]}, nil

	case "SET":

		if len(parts) < 3 {

			return Command{Op: OpUnknown}, fmt.Errorf("SET requires a key and a value")

		}

		return Command{Op: OpSet, Key: parts[1], Value: parts[2]}, nil

	case "DELETE":

		if len(parts) < 2 {

			return Command{Op: OpUnknown}, fmt.Errorf("DELETE requires a key")

		}

		return Command{Op: OpDelete, Key: parts[1]}, nil

	case "HELP":

		return Command{Op: OpHelp}, nil

	default:

		return Command{Op: OpUnknown}, fmt.Errorf("unknown command: %q", parts[0])

	}

}
