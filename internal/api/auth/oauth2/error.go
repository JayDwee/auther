package oauth2

type InvalidRequestError struct{}

func (e InvalidRequestError) Error() string {
	return "invalid_request"
}

type InvalidClientError struct{}

func (e InvalidClientError) Error() string {
	return "invalid_client"
}

type InvalidGrantError struct{}

func (e InvalidGrantError) Error() string {
	return "invalid_grant"
}

type UnauthorizedClientError struct{}

func (e UnauthorizedClientError) Error() string {
	return "unauthorized_client"
}

type UnsupportedGrantTypeError struct{}

func (e UnsupportedGrantTypeError) Error() string {
	return "unsupported_grant_type"
}

type InvalidScopeError struct{}

func (e InvalidScopeError) Error() string {
	return "invalid_scope"
}
