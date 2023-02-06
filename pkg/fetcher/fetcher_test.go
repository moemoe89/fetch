package fetcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		args    args
		wantErr error
	}

	tests := map[string]func(t *testing.T) test{
		"Successfully init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: defaultOptions,
				},
				wantErr: nil,
			}
		},
		"Failed init New": func(t *testing.T) test {
			t.Helper()

			return test{
				args: args{
					opts: []Option{
						WithHTTPClient(nil),
					},
				},
				wantErr: errFailedSetHTTPClient,
			}
		},
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			tt := fn(t)

			_, err := New(tt.args.opts...)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
