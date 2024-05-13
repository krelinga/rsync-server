package main

import (
    "log"
    "os/exec"
)

func logRsyncVersion() error {
    cmd := exec.Command("rsync", "-V")
    cmd.Stdout = log.Default().Writer()
    cmd.Stderr = log.Default().Writer()
    return cmd.Run()
}
