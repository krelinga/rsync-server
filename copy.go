package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os/exec"
    "syscall"
    "time"

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
    pipe, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("Could not get stdout pipe: %w", err)
    }
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("Could not start process: %w", err)
    }
    s := bufio.NewScanner(pipe)
    proc := cmd.Process
    time.Sleep(time.Second)
    requestOutput := func() {
        if proc != nil {
            // there's a race where cmd.Process could become nil before we can
            // grab it for very short-running commands, so we need to check
            // for that here.

            // Send a signal to force rsync to output a progress line after the
            // current file is completed.
            // Swallow the error, if any.
            _ = proc.Signal(syscall.SIGVTALRM)
        }
    }
    requestOutput()
    for s.Scan() {
        line := s.Text()
        log.Println(line)
        requestOutput()
    }
    if err := s.Err(); err != nil {
        return nil, fmt.Errorf("Could not finish scanner: %w", err)
    }
    if err := cmd.Wait(); err != nil {
        return nil, fmt.Errorf("Error waiting for command to finish: %w", err)
    }

    return &pb.CopyReply{}, nil
}
