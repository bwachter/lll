package main

import (
	"fmt"
	readline "github.com/chzyer/readline"
	syscall "golang.org/x/sys/unix"
	"os"
	"strings"
)

func make_directory(directory string) {
	err := os.MkdirAll(directory, 0750)
	if err != nil {
		fmt.Printf("mkdir: error creating %s: %s\n", directory, err)
	}
}

func list_directory(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("ls: error: %s\n", err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func mount(source string, target string, fstype string, flags uintptr) {
	err := syscall.Mount(source, target, fstype, flags, "")
	if err != nil {
		fmt.Printf("mount: error mounting %s at %s as %s: %s\n", source, target, fstype, err)
	}
}

func main() {
	// The kernel mounts the root filesystem passed via boot parameters read only.
	// It is the job of init to make it available read/write. You can comment this
	// line to see how the behaviour changes if we don't do this.
	mount("root", "/", "auto", syscall.MS_REMOUNT)

	// While typically the job of bootstrap scripts, it is not unheard of of init
	// (or scripts started by init) to create or adjust basic directories. The
	// following creates the most common top level directories. Some of those
	// often just serve as mountpoints for separate filesystems
	fmt.Printf("creating directories\n")
	make_directory("/etc")
	make_directory("/bin")
	make_directory("/usr")
	make_directory("/var")
	make_directory("/home")
	make_directory("/mnt")
	make_directory("/opt")
	make_directory("/srv")
	make_directory("/boot")
	// the following 4 directories are used as mountpoint for special virtual
	// filesystems on a modern linux. Traditionally those filesystems would be
	// configured by /etc/fstab, which init or a script started by init would
	// parse - but some init systems (including systemd) nowadays just mount
	// some of them always as the system would be useless without it
	make_directory("/sys")
	make_directory("/tmp")
	make_directory("/proc")
	// dev used to just contain statically created device nodes - and can still
	// be used as such. For modern systems it typically is more sensible to create
	// the device nodes dynamically - either by a filesystem where device nodes
	// automatically show up, a tmpfs where something (like udev) dynamically
	// manages device nodes from userland, or a combination of both.
	make_directory("/dev")

	// this mounts the special filesystems mentioned above:
	// tmpfs - an empty filesystem for temporary files in RAM
	mount("tmpfs", "/tmp", "tmpfs", syscall.MS_NODEV)
	// devtmpfs, a tmpfs, prepopulated with some device nodes for the hyprid
	// approach mentioned above. For simple systems this is enough - for more
	// complex ones something like udev will be required on top of this
	mount("devtmpfs", "/dev", "devtmpfs", syscall.MS_NOSUID)
	// proc, showing process information, and - on Linux - a bunch of other
	// information. A lot of the stuff moved to sysfs over time
	mount("proc", "/proc", "proc", syscall.MS_NODEV|syscall.MS_NOSUID|syscall.MS_NOEXEC)
	// sysfs, showing a lot of system info, and allowing to manipulate some
	// of them. Most systems will have additional special filesystems mounted
	// to mountpoints inside of sysfs nowadays
	mount("sysfs", "/sys", "sysfs", syscall.MS_NODEV|syscall.MS_NOSUID|syscall.MS_NOEXEC)

	// drop into a minimalistic shell to inspect the system. This pretty much just
	// needs the ability to run external commands now, and it can bring up a
	// full system
	// useful commands to also add here are probably cat, mount, mkdir, mv, cp
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		arg := strings.Split(line, " ")

		switch {
		case arg[0] == "cd":
			if len(arg) > 1 {
				err := os.Chdir(arg[1])
				if err != nil {
					fmt.Println(err)
				}
			}
		case arg[0] == "pwd":
			path, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(path)
		case arg[0] == "ls":
			if len(arg) > 1 {
				list_directory(arg[1])
			} else {
				path, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
				}
				list_directory(path)
			}

		default:
			fmt.Printf("Invalid command: %s\n", line)
		}
	}
}
