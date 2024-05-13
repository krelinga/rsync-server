package main

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
    "testing"

    "github.com/google/uuid"
    "github.com/krelinga/rsync-server/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type testContainer struct {
    containerId string
}

func newTestContainer() testContainer {
    return testContainer{
        containerId: fmt.Sprintf("rsync-server-docker-test-%s", uuid.NewString()),
    }
}

func captureOutput(cmd *exec.Cmd) *bytes.Buffer {
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    cmd.Stderr = cmdOutput
    return cmdOutput
}

func (tc testContainer) Build(t *testing.T) {
    t.Helper()
    cmd := exec.Command("docker", "image", "build", "-t", tc.containerId, ".")
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        t.Fatalf("could not build docker container: %s %s", err, cmdOutput)
    }
    t.Log("Finished building docker container.")
}

func (tc testContainer) Stop(t *testing.T) {
    t.Helper()
    cmd := exec.Command("docker", "container", "stop", tc.containerId)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        t.Fatalf("could not stop & delete docker container: %s %s", err, cmdOutput)
    }
    t.Log("Finished stopping & deleting docker container.")

    cmd = exec.Command("docker", "image", "rm", tc.containerId)
    cmdOutput = captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        t.Fatalf("could not delete docker image: %s %s", err, cmdOutput)
    }
    t.Log("Finished deleting docker image.")
}

func (tc testContainer) Run(t *testing.T) {
    t.Helper()
    cmd := exec.Command("docker")
    args := []string{
        "run",
        "--rm",
        "-d",
        "--name", tc.containerId,
        "-p", "25003:25003",
        tc.containerId,
    }
    cmdOutput := captureOutput(cmd)
    cmd.Args = append(cmd.Args, args...)
    if err := cmd.Run(); err != nil {
        t.Fatalf("Could not run docker container: %s %s", err, cmdOutput)
    }
    t.Log("Started Docker container.")
}

func testCopy(t *testing.T, c pb.RsyncClient) {
    req := &pb.CopyRequest{}
    _, err := c.Copy(context.Background(), req)
    if err == nil {
        t.Error("Expected an error.")
    }
}

func TestDocker(t *testing.T) {
    if testing.Short() {
        t.Skip()
        return
    }
    t.Parallel()
    tc := newTestContainer()
    tc.Build(t)
    tc.Run(t)
    defer tc.Stop(t)

    // Create a stub to the test server.
    const target = "docker-daemon:25003"
    creds := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err := grpc.DialContext(context.Background(), target, creds)
    if err != nil {
        t.Fatalf("Could not dial target: %s", err)
    }
    client := pb.NewRsyncClient(conn)
    t.Run("testCopy", func(t *testing.T) {
        testCopy(t, client)
    })
}
