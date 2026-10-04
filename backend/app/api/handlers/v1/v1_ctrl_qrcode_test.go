package v1

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteQRCodeProducesImage(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := writeQRCode(recorder, "https://catalog.example/t/test-token")
	require.NoError(t, err)
	assert.Equal(t, "image/jpeg", recorder.Header().Get("Content-Type"))
	assert.True(t, bytes.HasPrefix(recorder.Body.Bytes(), []byte{0xff, 0xd8, 0xff}))
}
