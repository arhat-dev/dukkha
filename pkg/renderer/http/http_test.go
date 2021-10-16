package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"arhat.dev/pkg/tlshelper"
	"arhat.dev/rs"
	"github.com/stretchr/testify/assert"

	dukkha_test "arhat.dev/dukkha/pkg/dukkha/test"
)

func TestDriver_RenderYaml(t *testing.T) {
	t.Run("TLS Basic Auth", func(t *testing.T) {
		expectPassword := true
		srv := httptest.NewUnstartedServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				user, password, hasBasicAuth := r.BasicAuth()
				assert.True(t, hasBasicAuth)

				assert.Equal(t, "foo", user)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(r.RequestURI))

				if expectPassword {
					assert.Equal(t, "bar", password)
				} else {
					assert.Equal(t, "", password)
				}
			},
		))

		srv.EnableHTTP2 = true
		srv.StartTLS()
		defer srv.Close()

		d := &driver{
			DefaultConfig: rendererHTTPConfig{
				User:     "foo",
				Password: "bar",

				TLS: tlshelper.TLSConfig{
					Enabled:            true,
					InsecureSkipVerify: true,
				},
			},
		}

		rc := dukkha_test.NewTestContext(context.TODO())
		assert.NoError(t, d.Init(rc))

		result, err := d.RenderYaml(rc, srv.URL+"/with-password")
		assert.NoError(t, err)
		assert.EqualValues(t, "/with-password", string(result.([]byte)))

		expectPassword = false
		result, err = d.RenderYaml(rc, rs.Init(&inputHTTPSpec{
			URL: srv.URL + "/no-password",
			Config: rendererHTTPConfig{
				User: "foo",

				TLS: tlshelper.TLSConfig{
					Enabled:            true,
					InsecureSkipVerify: true,
				},
			},
		}, nil))

		assert.NoError(t, err)
		assert.EqualValues(t, "/no-password", string(result.([]byte)))
	})
}
