package republique

import "errors"

var (
	errUnauthorised   = errors.New("Unauthorised")
	errSessionExpired = errors.New("Session Expired")
)
