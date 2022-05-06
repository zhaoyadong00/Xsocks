package internet

import (
	"context"
	"net"
)

type Dialer interface {
	// Dial dials a system connection to the given destination.
	Dial(ctx context.Context, dest net.Addr) (net.Conn, error)

	// Address returns the address used by this Dialer. Maybe nil if not known.
	Address() net.Addr
}
