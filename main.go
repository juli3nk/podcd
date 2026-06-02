package main

import (
    "log"
    "os/exec"
    "time"
)

func run(cmd string, args ...string) error {
    c := exec.Command(cmd, args...)
    c.Stdout = nil
    c.Stderr = nil
    return c.Run()
}

func getCommit(ref string) string {
    out, err := exec.Command("git", "rev-parse", ref).Output()
    if err != nil {
        log.Fatal(err)
    }
    return string(out)
}

func main() {
    repoPath := "/opt/gitops/repo"

    if err := run("git", "-C", repoPath, "fetch"); err != nil {
        log.Fatal(err)
    }

    lastCommit := getCommit("HEAD")

    for {
        time.Sleep(30 * time.Second)

        run("git", "-C", repoPath, "fetch")

        newCommit := getCommit("origin/main")

        if lastCommit != newCommit {
            log.Println("New version detected")

            run("git", "-C", repoPath, "pull")

            changedFiles := 
            changedServices := 

            for f, _ := changedFiles {
                run("ln", "-sf", fmt.Sprintf("%s/containers/%s", repoPath, f), "/etc/containers/systemd/")
            } 

            run("systemctl", "daemon-reload")

            for svc, _ := changedServices{
                run("systemctl", "restart", fmt.Sprintf("%s.service", svc))
            } 

            lastCommit = newCommit
        }
    }
}
