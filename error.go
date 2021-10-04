package golf

import "github.com/fhofherr/golf/internal/golferr"

// Error obtains an error that may have occurred during logging.
//
// In order to support this feature logger must implement the following method:
//
//     Err() error
//
// The Err method must return nil if no error occurred. Error itself returns
// whatever Err returned, or nil if logger does not implement Err.
func Error(logger Logger) error {
	e, ok := logger.(errorer)
	if !ok {
		Logf(logger, "%s: %T does not implement Err", golferr.MsgUnsupported, logger)
		return nil
	}
	return e.Err()
}

type errorer interface {
	Err() error
}
