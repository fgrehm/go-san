package sanast

import (
	token "github.com/fgrehm/go-san/token"
)

// ReachabilityDefinition represents the reachability information defined on
// the SAN file representing the reachable state space of the SAN model
type ReachabilityDefinition struct {
	Tokens     []token.Token
	Expression *Expression
}
