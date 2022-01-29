package aws

const (
	ERR_FAILED_START_INSTANCE = iota
	ERR_INVALID_RESPONSE_START_INSTANCE
	ERR_INSTANCE_ALREADY_STARTED
	ERR_STARTING_INSTANCE
	ERR_FAILED_WAIT_START_INSTANCE
	ERR_FAILED_GET_IP_ADDRESS
	ERR_INVALID_RESPONSE_GET_IP_ADDRESS
	ERR_FAILED_STOP_INSTANCE
	ERR_INVALID_RESPONSE_STOP_INSTANCE
	ERR_INSTANCE_ALREADY_STOPPED
	ERR_STOPPING_INSTANCE
	ERR_FAILED_WAIT_STOP_INSTANCE
)

type StatusError struct {
	Code int // ERR_で始まる定数のみ
	Err  error
}

func NewStatusError(code int, err error) StatusError {
	return StatusError{
		Code: code,
		Err:  err,
	}
}

func (e StatusError) IsEmpty() bool {
	return e == StatusError{}
}
