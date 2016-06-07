package goshim

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/Songmu/wrapcommander"
	"github.com/mitchellh/go-homedir"
)

const version = "0.0.1"

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage:
  goshim [-fhv] ./path/to/pkg [args...]

Verion: %s

Better `+"`go run`"+`. Build go codes transparently and exec
`, version)
	fmt.Println()
	flag.PrintDefaults()
}

var (
	helpFlag    = flag.Bool("h", false, "help")
	forceFlag   = flag.Bool("f", false, "force rebuild")
	verboseFlag = flag.Bool("v", false, "output verbosely")
	execFunc    = syscall.Exec
)

func RunWithFlag() int {
	flag.Parse()
	return Run(flag.Args())
}

// Run the goshim
func Run(args []string) int {
	if len(args) < 1 || *helpFlag {
		printUsage()
		return 1
	}
	subcmdOrSrcdir := args[0]

	// no subcommands for now
	srcdir := subcmdOrSrcdir
	fi, err := os.Stat(srcdir)
	if err != nil || !fi.IsDir() {
		log.Fatalf("not a directory: %s", srcdir)
		return 1
	}
	list, err := srcList(srcdir)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	dst := binDst(srcdir)
	if *forceFlag || isRebuildRequired(srcdir, dst, list) {
		if *verboseFlag {
			fmt.Printf("rebuild\n\n")
		}
		err := build(srcdir, dst)
		if err != nil {
			log.Fatal(err)
			return 1
		}
	} else {
		if *verboseFlag {
			fmt.Printf("use cache\n\n")
		}
	}
	err = execFunc(dst, append([]string{dst}, args[1:]...), os.Environ())
	return wrapcommander.ResolveExitCode(err)
}

func isRebuildRequired(dir, dst string, list []string) bool {
	dstFi, err := os.Stat(dst)
	if os.IsNotExist(err) {
		return true
	}
	dstMtime := dstFi.ModTime()
	for _, f := range list {
		f = filepath.Join(dir, f)
		fi, err := os.Stat(f)
		if err != nil {
			continue
		}
		if fi.ModTime().After(dstMtime) {
			return true
		}
	}
	return false
}

var cacheDir = func() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Panic(err)
	}
	return filepath.Join(home, ".goshim/bin")
}()

func binDst(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Panic(err)
	}
	if runtime.GOOS == "windows" {
		dir = strings.Replace(dir, ":", "", -1) + ".exe"
	}
	return filepath.Join(cacheDir, dir)
}

func normalizeDir(dir string) string {
	if filepath.IsAbs(dir) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		dir, err = filepath.Rel(cwd, dir)
		if err != nil {
			log.Panic(err)
		}
	}
	if !strings.HasPrefix(dir, ".") {
		dir = "." + string(os.PathSeparator) + dir
	}
	return dir
}

func build(dir, dst string) error {
	dir = normalizeDir(dir)
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}
	_, stderr, err := runCmd("go", "build", "-o", dst, dir)
	if err != nil {
		return fmt.Errorf(stderr)
	}
	return nil
}

func srcList(dir string) ([]string, error) {
	dir = normalizeDir(dir)
	stdout, stderr, err := runCmd("go", "list", "-f", `{{ join .GoFiles "\x00" }}`, dir)
	if err != nil {
		return nil, fmt.Errorf(stderr)
	}
	return strings.Split(strings.TrimSpace(stdout), "\x00"), err
}

func runCmd(cmdArgs ...string) (string, string, error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
