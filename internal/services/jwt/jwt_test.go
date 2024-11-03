package jwt

import (
	"testing"
	"time"
	"verni/internal/services/logging"

	"github.com/google/uuid"
)

func createConfig() DefaultConfig {
	return DefaultConfig{
		RefreshTokenLifetimeHours: 24 * 30,
		AccessTokenLifetimeHours:  1,
		RefreshTokenSecret:        "RefreshTokenSecret",
		AccessTokenSecret:         "AccessTokenSecret",
	}
}

func TestIssuedRefreshTokenIsValid(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	if err := service.ValidateRefreshToken(token); err != nil {
		t.Fatalf("ValidateRefreshToken err: %v", err)
	}
}

func TestIssuedRefreshTokenIsNotAnAccessToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	validateAccessTokenError := service.ValidateAccessToken(AccessToken(token))
	if validateAccessTokenError == nil {
		t.Fatalf("refresh token is recognized as a valid access token")
	} else if validateAccessTokenError.Code != CodeTokenInvalid {
		t.Fatalf("ValidateAccessToken unexpected err: %v", validateAccessTokenError)
	}
}

func TestIssuedRefreshTokenSubjectInaccessibleAsAnAccessToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	subjectFromTokenAsAccessToken, err := service.GetAccessTokenSubject(AccessToken(token))
	if err == nil {
		t.Fatalf("unexpected valid subject from access token %s", subjectFromTokenAsAccessToken)
	} else if err.Code != CodeTokenInvalid {
		t.Fatalf("unexpected err getting subject from access token %v", err)
	}
}

func TestIssuedRefreshTokenSubject(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	subjectFromToken, err := service.GetRefreshTokenSubject(token)
	if err != nil {
		t.Fatalf("GetRefreshTokenSubject err: %v", err)
	}
	if subject != subjectFromToken {
		t.Fatalf("subjects did not match %s != %s", subject, subjectFromToken)
	}
}

func TestExpiredRefreshToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now().Add(-(time.Hour*time.Duration(createConfig().RefreshTokenLifetimeHours) + time.Hour))
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	err = service.ValidateRefreshToken(token)
	if err == nil {
		t.Fatalf("outdated token should not be valid")
	} else if err.Code != CodeTokenExpired {
		t.Fatalf("outdated token unexpected validation err %v", err)
	}
}

func TestRefreshTokenValidOnTheLastMinute(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now().Add(-(time.Hour*time.Duration(createConfig().RefreshTokenLifetimeHours) - time.Minute))
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueRefreshToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	if err := service.ValidateRefreshToken(token); err != nil {
		t.Fatalf("ValidateRefreshToken err: %v", err)
	}
}

func TestIssuedAccessTokenIsValid(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueAccessToken err: %v", err)
	}
	if err := service.ValidateAccessToken(token); err != nil {
		t.Fatalf("ValidateAccessToken err: %v", err)
	}
}

func TestIssuedAccessTokenIsNotARefreshToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueAccessToken err: %v", err)
	}
	validateRefreshTokenError := service.ValidateRefreshToken(RefreshToken(token))
	if validateRefreshTokenError == nil {
		t.Fatalf("access token is recognized as a valid refresh token")
	} else if validateRefreshTokenError.Code != CodeTokenInvalid {
		t.Fatalf("ValidateRefreshToken unexpected err: %v", validateRefreshTokenError)
	}
}

func TestIssuedAccessTokenSubjectInaccessibleAsAnRefreshToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueAccessToken err: %v", err)
	}
	subjectFromTokenAsRefreshToken, err := service.GetRefreshTokenSubject(RefreshToken(token))
	if err == nil {
		t.Fatalf("unexpected valid subject from refresh token %s", subjectFromTokenAsRefreshToken)
	} else if err.Code != CodeTokenInvalid {
		t.Fatalf("unexpected err getting subject from refresh token %v", err)
	}
}

func TestIssuedAccessTokenSubject(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now()
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueAccessToken err: %v", err)
	}
	subjectFromToken, err := service.GetAccessTokenSubject(token)
	if err != nil {
		t.Fatalf("GetAccessTokenSubject err: %v", err)
	}
	if subject != subjectFromToken {
		t.Fatalf("subjects did not match %s != %s", subject, subjectFromToken)
	}
}

func TestExpiredAccessToken(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now().Add(-(time.Hour + time.Hour))
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueAccessToken err: %v", err)
	}
	err = service.ValidateAccessToken(token)
	if err == nil {
		t.Fatalf("outdated token should not be valid")
	} else if err.Code != CodeTokenExpired {
		t.Fatalf("outdated token unexpected validation err %v", err)
	}
}

func TestAccessTokenValidOnTheLastMinute(t *testing.T) {
	service := DefaultService(
		createConfig(),
		logging.TestService(),
		func() time.Time {
			return time.Now().Add(-(time.Hour - time.Minute))
		},
	)
	subject := Subject(uuid.New().String())
	token, err := service.IssueAccessToken(subject)
	if err != nil {
		t.Fatalf("IssueRefreshToken err: %v", err)
	}
	if err := service.ValidateAccessToken(token); err != nil {
		t.Fatalf("ValidateRefreshToken err: %v", err)
	}
}
