package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/exec"
)

func Mount(mountpoint string) {
	fmt.Printf("%s", mountpoint)
	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func(){
	    // for _ := range ch {
	        // sig is a ^C, handle it
	        fmt.Printf("Control C")

	        exec.Command("umount", mountpoint)

	    //}
	}()

	fs.Serve(c, FS{})
}

type FS struct{}

func (FS) Root() (fs.Node, fuse.Error) {
	return Dir{}, nil
}

type Dir struct{}

func (Dir) Attr() fuse.Attr {
	return fuse.Attr{Inode: 1, Mode: os.ModeDir | 0555}
}

func (Dir) Lookup(name string, intr fs.Intr) (fs.Node, fuse.Error) {
	if name == "hello" {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: "hello", Type: fuse.DT_File},
}

func (Dir) ReadDir(intr fs.Intr) ([]fuse.Dirent, fuse.Error) {
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

func (File) Attr() fuse.Attr {
	return fuse.Attr{Mode: 0444}
}

func (File) ReadAll(intr fs.Intr) ([]byte, fuse.Error) {
	return []byte("hello, world\n"), nil
}