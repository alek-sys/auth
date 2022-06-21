package provider_test

import (
	"context"
	"fmt"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewOpenID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	devPort := rand.Intn(10_000) + 50_000
	expectedTestUserSub := fmt.Sprintf("test-user-%d", devPort)
	svc := auth.NewService(auth.Opts{
		SecretReader: token.SecretFunc(func(aud string) (string, error) {
			return "some-signing-key", nil
		}),
		Logger:      logger.Std,
		AvatarStore: avatar.NewNoOp(),
	})

	svc.AddDevOpenIDProvider(devPort)
	devAuth, err := svc.DevAuth()
	require.NoError(t, err)

	devAuth.Automatic = true
	devAuth.CustomizeIdTokenFn = func(m map[string]interface{}) map[string]interface{} {
		m["sub"] = expectedTestUserSub
		return m
	}

	go devAuth.Run(context.Background())
	defer devAuth.Shutdown()

	authHandler, _ := svc.Handlers()
	server := httptest.NewServer(authHandler)
	defer server.Close()

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	client := http.Client{
		Jar: jar,
	}

	time.Sleep(200 * time.Millisecond)

	resp, err := client.Get(server.URL + "/auth/dev/login")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	var cookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "JWT" {
			cookie = c
			break
		}
	}

	require.NotNil(t, cookie)
	claims, err := devAuth.ParseToken(cookie.Value)
	require.NoError(t, err)

	// check user details are from the ID token
	assert.Equal(t, expectedTestUserSub, claims.User.ID)
}