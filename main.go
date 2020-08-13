package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	os.Exit(_main())
}

func gitRepositoryRoot() (string, error) {
	var output bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func gitCollectFiles(command []string) ([]string, error) {
	var output bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var files []string
	scanner := bufio.NewScanner(&output)
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func filterCFamilyLanguage(files []string) []string {
	// Check C, C++, Objective-C, Objective-C++ files
	re := regexp.MustCompile(`\.(?:m?m|c?c|cpp|h|hpp)$`)

	var ret []string
	for _, file := range files {
		if re.MatchString(file) {
			ret = append(ret, file)
		}
	}

	return ret
}

func applyClangFormat(files []string, runAtRoot bool) error {
	args := make([][]string, 0)

	// If there are too many C family languages, then passing them to clang-format
	// at once causes error(too many arguments).
	j := 0
	k := 0
	loopSize := 50
	for i := 0; i < len(files); i += loopSize {
		j += loopSize
		if j > len(files) {
			j = len(files)
		}

		args = append(args, make([]string, 0))
		args[k] = append(args[k], "-i")
		args[k] = append(args[k], files[i:j]...)
	}

	dir, err := gitRepositoryRoot()
	if err != nil {
		return err
	}

	for _, arg := range args {
		cmd := exec.Command("clang-format", arg...)
		if runAtRoot {
			cmd.Dir = dir
		}

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func usage() {
	fmt.Printf("Usage: git-clang-format [-all] [-modified] [-staged] [-verbose] dirs...\n")
}

func _main() int {
	applyAll := flag.Bool("all", false, "apply for all files in repository")
	modified := flag.Bool("modified", false, "apply only modified files")
	staged := flag.Bool("staged", false, "apply only staged files")
	verbose := flag.Bool("verbose", false, "verbose output")
	help := flag.Bool("help", false, "show help message")

	flag.Parse()

	if *help {
		usage()
		return 0
	}

	var collectCommand []string
	runAtRoot := false
	if *applyAll {
		root, err := gitRepositoryRoot()
		if err != nil {
			fmt.Println(err)
			return 1
		}

		collectCommand = []string{"git", "ls-files", root}
	} else if *modified {
		collectCommand = []string{"git", "diff", "--name-only"}
		runAtRoot = true
	} else if *staged {
		collectCommand = []string{"git", "diff", "--cached", "--name-only"}
		runAtRoot = true
	} else {
		collectCommand = []string{"git", "ls-files"}
	}

	if flag.NArg() > 0 {
		collectCommand = append(collectCommand, flag.Args()...)
	}

	files, err := gitCollectFiles(collectCommand)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	cFiles := filterCFamilyLanguage(files)
	if len(cFiles) == 0 {
		fmt.Println("There are no C/C++/Objective-C files in this repository")
		return 1
	}

	if *verbose {
		for _, file := range cFiles {
			fmt.Println(file)
		}
	}

	if err := applyClangFormat(cFiles, runAtRoot); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
