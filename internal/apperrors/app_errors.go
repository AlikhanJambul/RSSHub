package apperrors

import "errors"

var (
	ErrInvalidURL     = errors.New("Invalid URL")
	ErrNameExists     = errors.New("Name already exists")
	ErrInvalidName    = errors.New("Invalid Name")
	ErrInvalidFeed    = errors.New("Invalid Feed")
	ErrCountWorker    = errors.New("Count of workers should be less than 15")
	ErrAggregatorStop = errors.New("Aggregator stopped")
)

func CheckError(err error) int {
	if err == ErrInvalidURL || err == ErrNameExists || err == ErrCountWorker {
		return 401
	}

	return 500
}
