package main

import (
	"fmt"
	"github.com/containers/storage/pkg/reexec"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func init() {
	reexec.Register("nsInitialisation", nsInitialisation)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialisation() {
	upperPath := os.Args[1]
	workdirPath := os.Args[2]
	mergedFsPath := os.Args[3]
	rootfsPath := os.Args[4]

	mergerFs := newMergerFs("/", mergedFsPath)
	overlayFs := newOverlayFs(mergerFs.Target, rootfsPath, upperPath, workdirPath)
	changeRoot := newChangeRoot(overlayFs.Target)

	if err := mergerFs.mount(); err != nil {
		log.Fatalf("Error mounting mergerfs - %s", err)
	}

	//if err := mountMergerFs("/", mergedFsPath); err != nil {
	//	fmt.Printf("Error setting up mergedfs - %s\n", err)
	//	os.Exit(1)
	//}
	//
	//if err := mountOverlayFs(mergedFsPath, rootfsPath, upperPath, workdirPath); err != nil {
	//	fmt.Printf("Error setting up overlayfs - %s\n", err)
	//	os.Exit(1)
	//}

	if err := mountDev(mergedFsPath, rootfsPath); err != nil {
		fmt.Printf("Error mounting /dev - %s\n", err)
		os.Exit(1)
	}

	if err := mountProc(mergedFsPath, rootfsPath); err != nil {
		fmt.Printf("Error mounting /proc - %s\n", err)
		os.Exit(1)
	}

	//exitFn, err := changeRoot(rootfsPath)
	//if err != nil {
	//	fmt.Printf("Error running changeRoot - %s\n", err)
	//	os.Exit(1)
	//}

	//if err := pivotRoot(rootfsPath); err != nil {
	//	fmt.Printf("Error running pivot_root - %s\n", err)
	//	os.Exit(1)
	//}

	fmt.Println("Setup done!")

	nsRun()

	//if err := exitFn(); err != nil {
	//	fmt.Printf("Error exiting chroot environment - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Exit change root")
	//}

	if err := unmountProc(rootfsPath); err != nil {
		fmt.Printf("Error unmounting proc - %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Unmount proc")
	}

	if err := unmountDev(rootfsPath); err != nil {
		fmt.Printf("Error unmounting dev - %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Unmount dev")
	}

	//if err := unmountOverlayFs(rootfsPath); err != nil {
	//	fmt.Printf("Error unmounting overlayfs - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Unmount overlayfs")
	//}
	//
	//if err := unmountMergerFs(mergedFsPath); err != nil {
	//	fmt.Printf("Error unmounting mergerfs - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Unmount mergerfs")
	//}
	//
	//// should be empty
	//if err := syscall.Rmdir(rootfsPath); err != nil {
	//	fmt.Printf("Error removing rootfs dir - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Removed rootfs dir")
	//}
	//
	//// should be empty
	//if err := syscall.Rmdir(mergedFsPath); err != nil {
	//	fmt.Printf("Error removing mergerfs dir - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Removed mergerfs dir")
	//}
	//
	//// does not need to be empty
	//if err := os.RemoveAll(workdirPath); err != nil {
	//	fmt.Printf("Error removing overlayfs workdir - %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Println("Removed overlayfs workdir")
	//}

}

func nsRun() {
	cmd := exec.Command("/bin/sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[ns-process]- # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func main() {
	paths := []string{
		"/tmp/overlay-auto-test/upper",
		"/tmp/overlay-auto-test/workdir",
		"/tmp/overlay-auto-test/mergedfs",
		"/tmp/overlay-auto-test/mount",
	}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	rootfsPath := paths[3]

	exitIfRootfsNotFound(rootfsPath)

	cmd := reexec.Command("nsInitialisation", paths[0], paths[1], paths[2], paths[3])

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	userToRootUidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		},
	}

	userGroupToRootGroupMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		},
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: userToRootUidMappings,
		GidMappings: userGroupToRootGroupMappings,
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the reexec.Command - %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for the reexec.Command - %s\n", err)
		os.Exit(1)
	}

	fmt.Println("DONE")
}

func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")

	// bind mount newroot to itself - this is a slight hack needed to satisfy the
	// pivot_root requirement that newroot and putold must not be on the same
	// filesystem as the current root
	if err := syscall.Mount(newroot, newroot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// create putold directory
	if err := os.MkdirAll(putold, 0700); err != nil {
		return err
	}

	// call pivot_root
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	// ensure current working directory is set to new root
	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount putold, which now lives at /.pivot_root
	putold = "/.pivot_root"
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err
	}

	// remove putold
	if err := os.RemoveAll(putold); err != nil {
		return err
	}

	return nil
}

func changeRoot(newroot string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(newroot); err != nil {
		closeErr := root.Close()
		if closeErr != nil {
			return nil, closeErr
			//return nil, errors.Wrap(err, closeErr.Error())
		}
		return nil, err
	}
	if err := syscall.Chdir("/"); err != nil {
		return nil, err
	}

	return func() error {
		defer func(root *os.File) {
			err := root.Close()
			if err != nil {
				log.Fatalln("Could not close root")
			}
		}(root)

		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

func mountMergerFs(source string, target string) error {
	command := "/usr/bin/mergerfs"
	args := []string{
		source,
		target,
	}

	fmt.Println("trying to merge ", source, " into ", target)

	if _, err := exec.Command(command, args...).CombinedOutput(); err != nil {
		return err
	}

	fmt.Println("mergerfs was successfully mounted")

	return nil
}

func unmountMergerFs(target string) error {
	if err := syscall.Unmount(target, 0); err != nil {
		return err
	}
	return nil
}

func mountOverlayFs(source string, target string, upperdir string, workdir string) error {
	fstype := "overlay"

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", source, upperdir, workdir)
	if err := syscall.Mount("none", target, fstype, syscall.MS_NOSUID, opts); err != nil {
		return err
	}

	return nil
}

func unmountOverlayFs(target string) error {
	if err := syscall.Unmount(target, 0); err != nil {
		return err
	}
	return nil
}

func mountProc(root string, newroot string) error {
	source := filepath.Join(root, "/proc")
	target := filepath.Join(newroot, "/proc")
	fstype := ""
	flags := unix.MS_BIND | unix.MS_REC | unix.MS_PRIVATE
	data := ""

	if err := syscall.Mount(source, target, fstype, uintptr(flags), data); err != nil {
		return err
	}

	return nil
}

func unmountProc(newroot string) error {
	target := filepath.Join(newroot, "/proc")
	if err := syscall.Unmount(target, 0); err != nil {
		return err
	}
	return nil
}

func mountDev(root string, newroot string) error {
	source := filepath.Join(root, "/dev")
	target := filepath.Join(newroot, "/dev")
	fstype := ""
	flags := unix.MS_BIND | unix.MS_REC | unix.MS_PRIVATE
	data := ""

	if err := syscall.Mount(source, target, fstype, uintptr(flags), data); err != nil {
		return err
	}

	return nil
}

func unmountDev(newroot string) error {
	target := filepath.Join(newroot, "/dev")
	if err := syscall.Unmount(target, 0); err != nil {
		return err
	}
	return nil
}

func exitIfRootfsNotFound(rootfsPath string) {
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		usefulErrorMsg := fmt.Sprintf(`
"%s" does not exist.
Please create this directory and unpack a suitable root filesystem inside it.
An example rootfs, BusyBox, can be downloaded and unpacked as follows:
wget "https://raw.githubusercontent.com/teddyking/ns-process/4.0/assets/busybox.tar"
mkdir -p %s
tar -C %s -xf busybox.tar
`, rootfsPath, rootfsPath, rootfsPath)

		fmt.Println(usefulErrorMsg)
		os.Exit(1)
	}
}
