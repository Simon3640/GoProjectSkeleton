package contractsobservability

import "time"

// Time abstraction for testing
type Clock interface {
	Now() time.Time
}
