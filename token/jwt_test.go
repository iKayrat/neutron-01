package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"neutron0.1/util"
)

func TestJwtMaker(t *testing.T) {
	maker := NewJwtManager()

	id := util.RandomInt(0, 10)
	duration := 15 * time.Minute

	expiredAt := time.Now().Add(duration)

	token, err := maker.Create(id)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token.AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// req, err := http.NewRequest("GET", "/health-check", nil)
	// require.NoError(t, err)

	// req.Header.Set("Authorization", token.AccessToken)
	// svr := httptest.NewRecorder()
	// svr.Header().Set("Authorization", token.AccessToken)

	uuid, err := getUuid(payload)
	require.NoError(t, err)

	tm := time.Unix(token.AtExpires, 0)

	require.NotZero(t, payload)
	require.Equal(t, id, int64(uuid.Userid))
	require.WithinDuration(t, expiredAt, tm, 60*time.Minute)
}
