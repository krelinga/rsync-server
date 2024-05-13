package main

import (
    "context"
    "errors"
    "log"
    "net"

    "github.com/krelinga/rsync-server/pb"
    "google.golang.org/grpc"
)

type RsyncServer struct {
    pb.UnimplementedRsyncServer
}

func (rs *RsyncServer) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
    return nil, errors.New("Not implemented.")
}

func MainOrError() error {
    lis, err := net.Listen("tcp", ":25002")
    if err != nil {
        return err
    }
    grpcServer := grpc.NewServer()
    pb.RegisterRsyncServer(grpcServer, &RsyncServer{})
    grpcServer.Serve(lis)  // Runs as long as the server is alive.

    return nil
}

func main() {
    if err := MainOrError(); err != nil {
        log.Fatal(err)
    }
}
