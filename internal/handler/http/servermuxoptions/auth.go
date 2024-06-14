package servermuxoptions

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	_ "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

func WithAuthCookieToAuthMetadata(authCookieName, authMetadataName string) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			return nil
		}

		return metadata.New(map[string]string{
			authMetadataName: cookie.Value,
		})
	})
}

func WithAuthMetadataToAuthCookie(authCookieName, authMetadataName string, expiresInDuration time.Duration) runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, msg proto.Message) error {
		metadata, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return nil
		}

		authMetadataValues := metadata.Get(authMetadataName)
		if len(authMetadataValues) == 0 {
			return nil
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    authMetadataValues[0],
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(expiresInDuration),
		})

		return nil
	})
}
