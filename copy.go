package main

import (
    "context"
    "errors"

    "github.com/krelinga/rsync-server/pb"
)

func copyImpl(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
    return nil, errors.New("Not implemented.")
}
