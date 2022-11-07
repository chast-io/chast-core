package filetree

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"

	"chast.io/core/internal/post_processing/pipelinereport/internal/diff"
	"github.com/ttacon/chalk"
	"github.com/xlab/treeprint"
)

const unionFsHiddenPathSuffix = "_HIDDEN~"

func ToString(rootPath string, changeDiff *diff.ChangeDiff, printFullRootPath bool, colorize bool) (string, error) {
	var rootName string
	if printFullRootPath {
		rootName = rootPath
	} else {
		rootName = filepath.Base(rootPath)
	}

	tree, fileTreeBuildError := buildFileTree(rootPath, changeDiff, rootName, colorize)
	if fileTreeBuildError != nil {
		return "", errors.Wrap(fileTreeBuildError, "failed to build file tree")
	}

	return tree.String(), nil
}

func buildFileTree( //nolint:ireturn // return tree of tree library
	rootPath string,
	changeDiff *diff.ChangeDiff,
	rootName string,
	colorize bool,
) (treeprint.Tree, error) {
	tree := treeprint.NewWithRoot(rootName)
	nodeMap := make(map[string]treeprint.Tree)
	nodeMap[rootPath] = tree

	visit := func(path string, info os.FileInfo, err error) error {
		if path == rootPath {
			return nil
		}

		name := nameWithChangeType(path, info, changeDiff, colorize)

		parent := getParent(path)
		if info.IsDir() {
			handleDir(path, parent, nodeMap, name, tree)
		} else {
			handleFile(parent, nodeMap, name, tree)
		}

		return nil
	}

	if err := filepath.Walk(rootPath, visit); err != nil {
		return nil, errors.Wrap(err, "failed to walk file tree")
	}

	return tree, nil
}

func handleFile(parent string, nodeMap map[string]treeprint.Tree, name string, tree treeprint.Tree) {
	if parent != "" {
		nodeMap[parent].AddNode(name)
	} else {
		tree.AddNode(name)
	}
}

func handleDir(path string, parent string, nodeMap map[string]treeprint.Tree, name string, tree treeprint.Tree) {
	if parent != "" {
		nodeMap[path] = nodeMap[parent].AddBranch(name)
	} else {
		nodeMap[path] = tree.AddBranch(name)
	}
}

func getParent(path string) string {
	parent := ""
	segments := strings.Split(path, string(filepath.Separator))

	if len(segments) > 1 {
		parent = strings.Join(segments[:len(segments)-1], string(filepath.Separator))
	}

	return parent
}

func nameWithChangeType(path string, info os.FileInfo, changeDiff *diff.ChangeDiff, colorize bool) string {
	name := info.Name()

	switch changeDiff.Diffs[path].FileStatus {
	case diff.Deleted:
		name = strings.TrimSuffix(name, unionFsHiddenPathSuffix)
		name = "[-] " + name

		if colorize {
			name = chalk.Red.Color(name)
		}
	case diff.Added:
		name = "[+] " + name

		if colorize {
			name = chalk.Green.Color(name)
		}
	case diff.Modified:
		name = "[~] " + name

		if colorize {
			name = chalk.Blue.Color(name)
		}
	}

	return name
}
