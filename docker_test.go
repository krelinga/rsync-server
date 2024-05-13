package main

import (
    "bytes"
    "context"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "testing"

    "github.com/google/go-cmp/cmp"
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
    // Get a copy of all the logs before shutting down.
    cmd := exec.Command("docker", "container", "logs", tc.containerId)
    cmd.Stdout = log.Default().Writer()
    cmd.Stderr = log.Default().Writer()
    if err := cmd.Run(); err != nil {
        t.Fatalf("Could not get container logs before stopping.")
    }

    cmd = exec.Command("docker", "container", "stop", tc.containerId)
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

func (tc testContainer) Run(t *testing.T, dataDir string) {
    t.Helper()
    cmd := exec.Command("docker")
    mountCfg := fmt.Sprintf("type=bind,source=%s,target=/testdata", dataDir)
    userCfg := fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
    args := []string{
        "run",
        "--rm",
        "-d",
        "--name", tc.containerId,
        "-p", "25003:25003",
        "--mount", mountCfg,
        // This is needed so that generated files & directories have the correct ownership.
        "--user", userCfg,
        tc.containerId,
    }
    cmdOutput := captureOutput(cmd)
    cmd.Args = append(cmd.Args, args...)
    if err := cmd.Run(); err != nil {
        t.Fatalf("Could not run docker container: %s %s", err, cmdOutput)
    }
    t.Log("Started Docker container.")
}

func TestDocker(t *testing.T) {
    if testing.Short() {
        t.Skip()
        return
    }
    t.Parallel()
    dataDir, err := filepath.Abs("testdata")
    if err != nil {
        t.Fatal(err)
    }
    if err := os.Mkdir(dataDir, 0755); err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(dataDir)
    tc := newTestContainer()
    tc.Build(t)
    tc.Run(t, dataDir)
    defer tc.Stop(t)

    // Create a stub to the test server.
    const target = "docker-daemon:25003"
    creds := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err := grpc.DialContext(context.Background(), target, creds)
    if err != nil {
        t.Fatalf("Could not dial target: %s", err)
    }
    client := pb.NewRsyncClient(conn)
    t.Run("empty request", func(t *testing.T) {
        req := &pb.CopyRequest{}
        _, err := client.Copy(context.Background(), req)
        if err == nil {
            t.Error("Expected an error.")
        }
    })
    t.Run("directory is copied", func(t *testing.T) {
        inDir := filepath.Join(dataDir, "in")
        mkdir := func(p string) {
            if err := os.Mkdir(p, 0755); err != nil {
                t.Fatal(err)
            }
        }
        setContents := func(p, s string) {
            f, err := os.Create(p)
            if err != nil {
                t.Fatal(err)
            }
            _, err = fmt.Fprintf(f, "%s", s)
            if err != nil {
                t.Fatal(err)
            }
            f.Close()
        }
        mkdir(inDir)
        setContents(filepath.Join(inDir, "a.txt"), "a.txt")
        setContents(filepath.Join(inDir, "b.txt"), "b.txt")
        inSubDir := filepath.Join(inDir, "sub")
        mkdir(inSubDir)
        outDir := filepath.Join(dataDir, "out")
        mkdir(outDir)
        setContents(filepath.Join(inSubDir, "c.txt"), "c.txt")
        req := &pb.CopyRequest{
            InPath: "/testdata/in",
            OutSuperPath: "/testdata/out",
        }
        _, err = client.Copy(context.Background(), req)
        if err != nil {
            t.Error(err)
            return
        }
        ok := func(p, c string) {
            fullP := filepath.Join(outDir, p)
            data, err := os.ReadFile(fullP)
            if err != nil {
                t.Error(err)
                return
            }
            dataStr := string(data)
            if !cmp.Equal(c, dataStr) {
                t.Error(cmp.Diff(c, dataStr))
            }
        }
        ok("in/a.txt", "a.txt")
        ok("in/b.txt", "b.txt")
        ok("in/sub/c.txt", "c.txt")
    })
}
