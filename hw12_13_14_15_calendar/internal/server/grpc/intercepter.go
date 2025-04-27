package internalgrpc

import (
	"context"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	p, _ := peer.FromContext(ctx)
	ip := middleware.ExtractIP(p.Addr.String())

	userAgent := "-"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
	}

	resp, err := handler(ctx, req)

	logLine := middleware.CommonLogFormat(
		ip,
		"gRPC",
		info.FullMethod,
		"HTTP/2",
		int(status.Code(err)),
		time.Since(start),
		userAgent,
	)

	middleware.WriteLogToFile(logLine)

	return resp, err
}
