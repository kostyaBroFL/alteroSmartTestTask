// Package errors provides helper functions and structures for more natural
// error handling in some situations.
package errors

import (
	"context"
	"errors"
	"fmt"

	log_context "alteroSmartTestTask/common/log/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Convenience function to call errors.New() from the standard library.
func New(text string) error {
	return errors.New(text)
}

// Newf is a convenient function for creating formatted errors.
func Newf(format string, arguments ...interface{}) error {
	return fmt.Errorf(format, arguments...)
}

// This function hides the server error from the user and returns a readable replacement.
func ToFrontendError(
	context context.Context,
	originalError error,
	code codes.Code, newErrorMessage string,
) error {
	log_context.FromContext(context).Error(originalError)
	return status.Error(code, newErrorMessage)
}
