package interceptor;

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()
	if addr == "" {
		addr = "unknown"
	}
	i, e := handler(ctx, req)
	t := time.Now()
	ts := t.Format("2006/01/02 15:04:05.000000")
	if stat, ok := status.FromError(e); ok {
		fmt.Printf("%s %s %s %s\n", ts, addr, info.FullMethod, stat.Code())
	} else {
		fmt.Printf("%s %s %s %s\n", ts, addr, info.FullMethod, e)
	}
	return i, e
}
