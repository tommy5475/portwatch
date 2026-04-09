package filter

import "fmt"

// Builder provides a fluent API for constructing a Filter.
type Builder struct {
	rules []Rule
	err   error
}

// NewBuilder returns a new Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// DenyPort adds a rule that denies a single port for the given protocol.
func (b *Builder) DenyPort(protocol string, port uint16) *Builder {
	b.addRule(protocol, port, port, false)
	return b
}

// AllowPort adds a rule that allows a single port for the given protocol.
func (b *Builder) AllowPort(protocol string, port uint16) *Builder {
	b.addRule(protocol, port, port, true)
	return b
}

// DenyRange adds a rule that denies a port range for the given protocol.
func (b *Builder) DenyRange(protocol string, min, max uint16) *Builder {
	b.addRule(protocol, min, max, false)
	return b
}

// AllowRange adds a rule that allows a port range for the given protocol.
func (b *Builder) AllowRange(protocol string, min, max uint16) *Builder {
	b.addRule(protocol, min, max, true)
	return b
}

// Build constructs the Filter or returns the first validation error.
func (b *Builder) Build() (*Filter, error) {
	if b.err != nil {
		return nil, b.err
	}
	return New(b.rules)
}

func (b *Builder) addRule(protocol string, min, max uint16, allow bool) {
	if b.err != nil {
		return
	}
	if min > max {
		b.err = fmt.Errorf("filter builder: invalid range %d-%d", min, max)
		return
	}
	b.rules = append(b.rules, Rule{
		Protocol: protocol,
		PortMin:  min,
		PortMax:  max,
		Allow:    allow,
	})
}
