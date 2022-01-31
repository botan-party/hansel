package aws

import "fmt"

var (
	ErrFailedGetIpAddress           error = &StatusError{message: "IPアドレス取得時、コマンド実行に失敗 :"}
	ErrFailedStartInstance          error = &StatusError{message: "起動に失敗した :"}
	ErrFailedStopInstance           error = &StatusError{message: "停止に失敗した :"}
	ErrFailedWaitStartInstance      error = &StatusError{message: "起動待ちに失敗した"}
	ErrFailedWaitStopInstance       error = &StatusError{message: "停止待ちに失敗した :"}
	ErrInstanceAlreadyStarted       error = &StatusError{message: "既に起動している"}
	ErrInstanceAlreadyStopped       error = &StatusError{message: "既に停止している"}
	ErrInvalidResponseGetIpAddress  error = &StatusError{message: "IPアドレス取得時のレスポンスに異常 :"}
	ErrInvalidResponseStartInstance error = &StatusError{message: "起動時のレスポンスに異常 :"}
	ErrInvalidResponseStopInstance  error = &StatusError{message: "停止時のレスポンスに異常 :"}
	ErrStartingInstance             error = &StatusError{message: "起動処理実行中"}
	ErrStoppingInstance             error = &StatusError{message: "停止処理実行中"}
)

type StatusError struct {
	base    error
	message string
}

func (e *StatusError) Error() string {
	baseMessage := ""
	if e.base != nil {
		baseMessage = e.base.Error()
	}
	return fmt.Sprintf("%s %s", e.message, baseMessage)
}

func (e *StatusError) Unwrap() error {
	return e.base
}

func WrapError(err error, statusError error) error {
	e := statusError.(*StatusError)
	e.base = err
	return e
}
