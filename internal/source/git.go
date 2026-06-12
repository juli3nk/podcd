package source

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitSource struct {
	repoURL string
	branch  string
	path    string

	lastRevision string
}

func NewGit(repoURL, branch, path string) (*GitSource, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	g := &GitSource{
		repoURL: repoURL,
		branch:  branch,
		path:    path,
	}

	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if err := g.clone(); err != nil {
			return nil, err
		}
	}

	rev, err := g.getRevision()
	if err != nil {
		return nil, err
	}

	g.lastRevision = rev

	return g, nil
}

func (g *GitSource) Fetch() error {
	if _, err := os.Stat(g.path); os.IsNotExist(err) {
		return g.clone()
	}
	return g.pull()
}

func (g *GitSource) Diff() ([]Change, error) {
	cmd := exec.Command(
		"git", "-C", g.path,
		"diff", "--name-status",
		g.lastRevision, "HEAD",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseDiff(out), nil
}

func (g *GitSource) HasChanged() (bool, error) {
	rev, err := g.getRevision()
	if err != nil {
		return false, err
	}

	if rev != g.lastRevision {
		g.lastRevision = rev
		return true, nil
	}

	return false, nil
}

func (g *GitSource) Path() string {
	return g.path
}

func (g *GitSource) Revision() string {
	return g.lastRevision
}

func (g *GitSource) clone() error {
	cmd := exec.Command("git", "clone", "-b", g.branch, g.repoURL, g.path)
	return cmd.Run()
}

func (g *GitSource) pull() error {
	cmd := exec.Command("git", "-C", g.path, "pull")
	return cmd.Run()
}

func (g *GitSource) getRevision() (string, error) {
	cmd := exec.Command("git", "-C", g.path, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func parseDiff(data []byte) []Change {
	lines := strings.Split(string(data), "\n")

	var changes []Change

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		var t ChangeType

		switch parts[0] {
		case "A":
			t = Added
		case "M":
			t = Modified
		case "D":
			t = Deleted
		}

		changes = append(changes, Change{
			Path: parts[1],
			Type: t,
		})
	}

	return changes
}
