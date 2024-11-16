package server

import (
	"github.com/kaium123/order/internal/cache"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
)

func TestNewAPI(t *testing.T) {
	tests := []struct {
		name    string
		opts    TodoAPIServerOpts
		wantErr bool
	}{
		{
			name: "Valid configuration",
			opts: TodoAPIServerOpts{
				ListenPort: 8080,
				Config: model.Config{
					SQLite: model.SQLite{
						DBFilename: ":memory:",
					},
					UI: model.UI{
						URL: "http://localhost:3000",
					},
					SwaggerServer: model.Server{
						Enable: true,
						Port:   8081,
					},
					Redis: &cache.Config{
						Addr:     "localhost:6379",
						Password: "",
						DB:       5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid database configuration",
			opts: TodoAPIServerOpts{
				ListenPort: 8080,
				Config: model.Config{
					SQLite: model.SQLite{
						DBFilename: "/invalid/path/to/db",
					},
					UI: model.UI{
						URL: "http://localhost:3000",
					},
					SwaggerServer: model.Server{
						Enable: true,
						Port:   8081,
					},
					Redis: &cache.Config{
						Addr:     "localhost:6379",
						Password: "",
						DB:       5,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New()
			server, err := NewAPI(context.Background(), &InitNewAPI{
				TodoAPIServerOpts: tt.opts, Log: logger,
			})
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, server)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tt.opts.ListenPort, server.(*todoAPIServer).port)
				assert.IsType(t, &echo.Echo{}, server.(*todoAPIServer).engine)
				assert.IsType(t, &log.Logger{}, server.(*todoAPIServer).log)
				assert.IsType(t, &gorm.DB{}, server.(*todoAPIServer).db)
			}
		})
	}
}
