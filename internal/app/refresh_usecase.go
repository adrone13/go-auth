package app

import (
	"auth/internal/app/common/repositories"
	"auth/internal/app/sessions"
	"auth/internal/app/tokens"
	"auth/internal/app/users"
	"auth/internal/config"
	"auth/internal/logger"
	"auth/internal/utils/arrays"
	"context"
	"errors"
	"fmt"
	jwt "github.com/adrone13/gojwt"
)

/*
https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps-05#section-8

In particular, authorization servers:

o  MUST rotate refresh tokens on each use, in order to be able to
detect a stolen refresh token if one is replayed (described in
[oauth-security-topics] section 4.12)

o  MUST either set a maximum lifetime on refresh tokens OR expire if
the refresh token has not been used within some amount of time

o  upon issuing a rotated refresh token, MUST NOT extend the lifetime
of the new refresh token beyond the lifetime of the initial
refresh token if the refresh token has a preestablished expiration
time

For example:

o  A user authorizes an application, issuing an access token that
lasts 1 hour, and a refresh token that lasts 24 hours

o  After 1 hour, the initial access token expires, so the application
uses the refresh token to get a new access token

o  The authorization server returns a new access token that lasts 1
hour, and a new refresh token that lasts 23 hours

o  This continues until 24 hours pass from the initial authorization

o  At this point, when the application attempts to use the refresh
token after 24 hours, the request will fail and the application
will have to involve the user in a new authorization request

By limiting the overall refresh token lifetime to the lifetime of the
initial refresh token, this ensures a stolen refresh token cannot be
used indefinitely.

https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/

However, we can reduce the absolute token expiration time of tokens
to reduce the security risks of storing tokens in local storage.
This reduces the impact of a reflected XSS attack (but not of a
persistent one). A refresh token may have a long lifespan by
configuration. However, the defined long lifespan of a refresh token
is cut short with refresh token rotation. The refresh is only valid
within the lifespan of the access token, which would be short-lived.

Authorization servers MAY revoke refresh tokens automatically in case of a security event, such as:

password change
logout at the authorization server

Errors:
invalid_request
	The request is missing a required parameter, includes an
	unsupported parameter value (other than grant type),
	repeats a parameter, includes multiple credentials,
	utilizes more than one mechanism for authenticating the
	client, or is otherwise malformed.

invalid_client
	Client authentication failed (e.g., unknown client, no
	client authentication included, or unsupported
	authentication method).  The authorization server MAY
	return an HTTP 401 (Unauthorized) status code to indicate
	which HTTP authentication schemes are supported.  If the
	client attempted to authenticate via the "Authorization"
	request header field, the authorization server MUST
	respond with an HTTP 401 (Unauthorized) status code and
	include the "WWW-Authenticate" response header field
	matching the authentication scheme used by the client.

invalid_grant
	The provided authorization grant (e.g., authorization
	code, resource owner credentials) or refresh token is
	invalid, expired, revoked, does not match the redirection
	URI used in the authorization request, or was issued to
	another client.

unauthorized_client
	The authenticated client is not authorized to use this
	authorization grant type.

unsupported_grant_type
	The authorization grant type is not supported by the
	authorization server.

Refreshing an Access Token
	If the authorization server issued a refresh token to the client, the
	client makes a refresh request to the token endpoint by adding the
	following parameters using the "application/x-www-form-urlencoded"
	format per Appendix B with a character encoding of UTF-8 in the HTTP
	request entity-body:

	grant_type
		REQUIRED.  Value MUST be set to "refresh_token".

	refresh_token
		REQUIRED.  The refresh token issued to the client.

For example, the client makes the following HTTP request using
transport-layer security (with extra line breaks for display purposes
only):
POST /token HTTP/1.1
	Host: server.example.com
	Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
	Content-Type: application/x-www-form-urlencoded

	grant_type=refresh_token&refresh_token=tGzv3JOkF0XG5Qx2TlKWIA

The authorization server MUST verify the binding between the refresh
token and client identity whenever the client identity can be
authenticated.  When client authentication is not possible, the
authorization server SHOULD deploy other means to detect refresh
token abuse.

For example, the authorization server could employ refresh token
rotation in which a new refresh token is issued with every access
token refresh response.  The previous refresh token is invalidated
but retained by the authorization server.  If a refresh token is
compromised and subsequently used by both the attacker and the
legitimate client, one of them will present an invalidated refresh
token, which will inform the authorization server of the breach.
The authorization server MUST ensure that refresh tokens cannot be
generated, modified, or guessed to produce valid refresh tokens by
unauthorized parties.
*/

/*
Plan

1. Implement refresh token rotation when new access token is requested
2. Implement automatic replay detection (check for use of invalidated
refresh tokens in token family)
3. Enable absolute (e.g. 1 week) and inactivity lifetimes for refresh tokens
(https://auth0.com/blog/achieving-a-seamless-user-experience-with-refresh-token-inactivity-lifetimes/)

Absolute lifetime   30 days
Inactivity lifetime 15 days

Request
grant_type=refresh_token
client_id — The application’s client ID
client_secret — The application’s client secret. The client secret must be kept confidential.
refresh_token — Include the refresh token.

Response
{
  "token_type": "bearer",
  "team_id": 3074457358607431473,
  "access_token": "eyJtaXJvLm9yaWdpbiI6ImV1MDEifQ_o-P91OccaII0A63CDSK--x21xiI",
  "refresh_token": "eyJtaXJvLm9yaWdpbiI6ImV1MDEifQ_-PIBKmE9rzQuL3bUeAvUEGFEhLk",
  "scope": "boards:write boards:read identity:read",
  "expires_in": 3599
}
*/

type RefreshAuthUseCase struct {
	UserRepository     repositories.UserRepository
	SessionsRepository repositories.SessionsRepository
}

func (u *RefreshAuthUseCase) Execute(ctx context.Context, refreshToken string) (*Auth, error) {
	token, err := jwt.ParseCustomClaims(refreshToken, tokens.RefreshTokenClaims{})
	if err != nil {
		return nil, err
	}
	secret := config.Values.JwtSecret
	if !token.IsValid(secret) {
		return nil, errors.New("invalid_request")
	}

	user, err := u.UserRepository.FindById(ctx, users.UserId(token.Claims.Subject))
	if err != nil {
		logger.Error(fmt.Sprintf("user (id: %s) does not exist", token.Claims.Subject))

		return nil, err
	}
	session, err := u.SessionsRepository.FindById(ctx, sessions.SessionId(token.Claims.SessionId))
	if err != nil {
		logger.Error(fmt.Sprintf("session (id: %s) does not exist", token.Claims.SessionId))

		return nil, err
	}
	if session.Expired() {
		logger.Error("Logging out. Session has expired")

		// Logout user if session expired
		err = u.SessionsRepository.Delete(ctx, session.Id)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("access_denied")
	}

	if !arrays.Contains(session.RefreshTokens, refreshToken) {
		return nil, errors.New("invalid_request")
	}

	currentRefresh := arrays.Last(session.RefreshTokens)
	if refreshToken != currentRefresh {
		logger.Error(fmt.Sprintf("Logging out. Invalidated token provided: %s", refreshToken))

		// Logout user if invalidated refresh token received
		err = u.SessionsRepository.Delete(ctx, session.Id)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("access_denied")
	}

	access := tokens.CreateAccessToken(user)
	refresh := tokens.CreateRefreshToken(user, session)

	session.AddRefreshToken(refresh)
	err = u.SessionsRepository.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	auth := &Auth{
		AccessToken:  fmt.Sprintf("Bearer %s", access),
		TokenType:    "bearer",
		ExpiresIn:    config.Values.AccessTokenTtl,
		RefreshToken: refresh,
	}

	return auth, nil
}
