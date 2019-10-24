package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GitReviewer struct {
	config    *Config
	repoPaths []string

	erred   map[string]string
	messy   map[string]string
	ahead   map[string]string
	behind  map[string]string
	fetched map[string]string
	journal map[string]string
	omitted map[string]string
	skipped map[string]string
}

func NewGitReviewer(config *Config) *GitReviewer {
	return &GitReviewer{
		config: config,
		repoPaths: append(
			collectGitRepositories(config.GitRepositoryRoots),
			filterGitRepositories(config.GitRepositoryPaths)...,
		),
		erred:   make(map[string]string),
		messy:   make(map[string]string),
		ahead:   make(map[string]string),
		behind:  make(map[string]string),
		fetched: make(map[string]string),
		journal: make(map[string]string),
		omitted: make(map[string]string),
		skipped: make(map[string]string),
	}
}

func (this *GitReviewer) GitAnalyzeAll() {
	log.Printf("Analyzing %d git repositories...", len(this.repoPaths))
	log.Println("Legend: [!] = error; [M] = messy; [A] = ahead; [B] = behind; [F] = fetched; [O] = omitted; [S] = skipped;")
	reports := NewAnalyzer(workerCount).AnalyzeAll(this.repoPaths)
	for _, report := range reports {
		if len(report.StatusError) > 0 {
			this.erred[report.RepoPath] += report.StatusError
			log.Println(report.RepoPath, report.StatusError)
		}
		if len(report.FetchError) > 0 {
			this.erred[report.RepoPath] += report.FetchError
			log.Println(report.RepoPath, report.FetchError)
		}
		if len(report.RevListError) > 0 {
			this.erred[report.RepoPath] += report.RevListError
			log.Println(report.RepoPath, report.RevListError)
		}

		if len(report.StatusOutput) > 0 {
			this.messy[report.RepoPath] += report.StatusOutput
		}
		if len(report.RevListAhead) > 0 {
			this.ahead[report.RepoPath] += report.RevListAhead
		}
		if len(report.RevListBehind) > 0 {
			this.behind[report.RepoPath] += report.RevListBehind
		}
		if len(report.SkipOutput) > 0 {
			this.skipped[report.RepoPath] += report.SkipOutput
		}
		if len(report.OmitOutput) > 0 {
			this.omitted[report.RepoPath] += report.OmitOutput
		}

		if this.config.GitFetch && len(report.FetchOutput) > 0 {
			this.fetched[report.RepoPath] += report.FetchOutput + report.RevListOutput

			if this.canJournal(report) {
				this.journal[report.RepoPath] += report.FetchOutput + report.RevListOutput
			}
		}
	}
}
func (this *GitReviewer) canJournal(report *GitReport) bool {
	if !strings.Contains(report.RemoteOutput, "smartystreets") { // Exclude externals from code review journal.
		return false
	}
	if _, found := this.omitted[report.RepoPath]; found {
		return false
	}
	return true
}

func (this *GitReviewer) ReviewAll() {
	var review []map[string]string
	if this.config.ReviewError {
		review = append(review, this.erred)
	}
	if this.config.ReviewMessy {
		review = append(review, this.messy)
	}
	if this.config.ReviewAhead {
		review = append(review, this.ahead)
	}
	if this.config.ReviewBehind {
		review = append(review, this.behind)
	}
	if this.config.ReviewFetched {
		review = append(review, this.fetched)
	}
	if this.config.ReviewJournal {
		review = append(review, this.journal)
	}
	reviewable := sortUniqueKeys(review...)
	if len(reviewable) == 0 {
		log.Println("Nothing to review at this time.")
		return
	}

	printMapKeys(this.erred, "Repositories with git errors: %d")
	printMapKeys(this.messy, "Repositories with uncommitted changes: %d")
	printMapKeys(this.ahead, "Repositories ahead of origin master: %d")
	printMapKeys(this.behind, "Repositories behind origin master: %d")
	printMapKeys(this.fetched, "Repositories with new content since the last review: %d")
	printMapKeys(this.journal, "Repositories to be included in the final report: %d")
	printMapKeys(this.skipped, "Repositories that were skipped: %d")
	printStrings(reviewable, "Repositories to be reviewed: %d")

	in := prompt(fmt.Sprintf("Press <ENTER> to initiate the review process (will open %d review windows), or 'q' to quit...", len(reviewable)))
	if len(in) > 0 && in[0] == 'q' {
		os.Exit(0)
	}

	for _, path := range reviewable {
		log.Printf("Opening %s at %s", this.config.GitGUILauncher, path)
		var err error
		if this.config.GitGUILauncher == "gitk" {
			tmp, _ := os.Getwd()
			os.Chdir(path)
			err = exec.Command(this.config.GitGUILauncher, "--all").Run()
			os.Chdir(tmp)
		} else {
			err = exec.Command(this.config.GitGUILauncher, path).Run()
		}
		if err != nil {
			log.Println("Failed to open git GUI:", err)
		}
		time.Sleep(time.Millisecond * 250)
	}
}

//noinspection GoUnhandledErrorResult
func (this *GitReviewer) PrintCodeReviewLogEntry() {
	if len(this.journal) == 0 {
		return
	}

	prompt("Press <ENTER> to conclude review process and print code review log entry...")

	writer := this.config.OpenOutputWriter()
	defer func() { _ = writer.Close() }()

	fmt.Fprintln(writer)
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "##", time.Now().Format("2006-01-02"))
	fmt.Fprintln(writer)
	for _, review := range this.journal {
		fmt.Fprintln(writer, review)
	}
}

const workerCount = 16
