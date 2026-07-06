package server

import (
	"github.com/BenDowswell/go-kv-store/kvstore"
	"github.com/BenDowswell/go-kv-store/protocol"
)

// execute runs cmd's against s (store) and returns the plain-text response string.
// Used by both the TCP and UDP routes.
func execute(s kvstore.Store, cmd protocol.Command) string {
	switch cmd.Op {
	case protocol.OpGet:
		v, ok := s.Get(cmd.Key)
		if !ok {
			return "NOT FOUND"
		}
		return v
	case protocol.OpSet:
		s.Set(cmd.Key, cmd.Value)
		return "OK"
	case protocol.OpDelete:
		s.Delete(cmd.Key)
		return "OK"
	case protocol.OpHelp:
		return protocol.HelpText
	default:
		return "ERR unknown command"
	}
}
