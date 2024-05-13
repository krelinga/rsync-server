package main

import (
    "context"
    "log"
    "os/exec"

    "github.com/krelinga/rsync-server/pb"
)

func copyImpl(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
    args := []string{
        "-a",
        "--info=progress2",
        "-r",
        req.InPath,
        req.OutSuperPath,
    }
    cmd := exec.CommandContext(ctx, "rsync", args...)
    cmd.Stdout = log.Default().Writer()
    cmd.Stderr = log.Default().Writer()
    if err := cmd.Run(); err != nil {
        return nil, err
    }
    return &pb.CopyReply{}, nil
}
