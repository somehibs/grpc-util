package interceptor;

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type AuthProvider interface {
	Check(string, map[string][]string) error
}

var NoAuthMethods = map[string]bool{}
var AuthChecker AuthProvider

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

func AuthenticationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if NoAuthMethods[info.FullMethod] == true {
		return handler(ctx, req)
	}
	p, _ := peer.FromContext(ctx)
	peerAddr := p.Addr.String()
	data, _ := metadata.FromIncomingContext(ctx)
	err := AuthChecker.Check(peerAddr, data)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
