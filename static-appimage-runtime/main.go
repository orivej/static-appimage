package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/zipfs"
	"github.com/kardianos/osext"
	"github.com/orivej/e"
)

func main() {
	executable, err := osext.Executable()
	e.Exit(err)
	files, err := zipfs.NewZipTree(executable)
	e.Exit(err)

	mfs := zipfs.NewMemTreeFs(files)
	mfs.Name = fmt.Sprintf("fs(%s)", os.Args[0])

	opts := &nodefs.Options{
		AttrTimeout:  10 * time.Second,
		EntryTimeout: 10 * time.Second,
	}

	mnt, err := ioutil.TempDir("", ".mount_")
	e.Exit(err)

	server, _, err := nodefs.MountRoot(mnt, mfs.Root(), opts)
	e.Exit(err)

	go server.Serve()

	signals := make(chan os.Signal, 1)
	exitCode := 0
	go func() {
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		err = server.Unmount()
		e.Exit(err)
		err = os.Remove(mnt)
		e.Exit(err)

		os.Exit(exitCode)
	}()

	err = server.WaitMount()
	e.Exit(err)

	cmd := exec.Command(path.Join(mnt, "AppRun"), os.Args[1:]...) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if cmd.ProcessState != nil {
		if waitStatus, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			exitCode = waitStatus.ExitStatus()
			err = nil
		}
	}
	e.Print(err)

	signals <- syscall.SIGTERM
	runtime.Goexit()
}
