package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os/exec"
    "bytes"

    "github.com/krelinga/rsync-server/pb"
)

func copyImpl(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
    args := []string{
        "-a",
        "--info=progress2",
        "-r",
    }
    if req.BwLimitKbps > 0 {
        args = append(args, fmt.Sprintf("--bwlimit=%d", req.BwLimitKbps))
    }
    args = append(args, req.InPath, req.OutSuperPath)
    cmd := exec.CommandContext(ctx, "rsync", args...)
    pipe := &bytes.Buffer{}
    cmd.Stdout = pipe
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("Could not start process: %w", err)
    }
    if err := cmd.Wait(); err != nil {
        return nil, fmt.Errorf("Error waiting for command to finish: %w", err)
    }
    contents, err := io.ReadAll(pipe)
    if err != nil {
        return nil, fmt.Errorf("Could not read from pipe: %w", err)
    }
    log.Printf("Read %d bytes from pipe.\n", len(contents))
    log.Printf("Contents: %q\n", contents)

    return &pb.CopyReply{}, nil
}
