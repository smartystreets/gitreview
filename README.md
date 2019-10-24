# gitreview

gitreview facilitates visual inspection (code review) of git
repositories that meet any of the following criteria:

1. New content was fetched
2. Behind origin/master
3. Ahead of origin/master
4. Messy (have uncommitted state)
5. Throw errors for the required git operations (listed below)

We use variants of the following commands to ascertain the
status of each repository:

- `git remote`    (shows remote address)
- `git status`    (shows uncommitted files)
- `git fetch`     (finds new commits/tags/branches)
- `git rev-list`  (lists commits behind/ahead of master)

...all of which should be safe enough. 

Each repository that meets any criteria above will be
presented for review. After all reviews are complete a
concatenated report of all output from `git fetch` for
repositories that were behind their origin is printed to
stdout. Only repositories with "smartystreets" in their
path are included in this report.

Repositories are identified for consideration from path values
supplied as non-flag command line arguments or via the roots
flag (see details below).

Installation:

    go get -u github.com/smartystreets/gitreview


Skipping Repositories:

If you have repositories in your list that you would rather not review,
you can mark them to be skipped by adding a config variable to the
repository. The following command will produce this result:

    git config --add review.skip true


Omitting Repositories:

If you have repositories in your list that you would still like to audit
but aren`t responsible to sign off (it`s code from another team), you can 
mark them to be omitted from the final report by adding a config variable
to the repository. The following command will produce this result:

    git config --add review.omit true


CLI Flags:


```
  -fetch
    	When false, suppress all git fetch operations via --dry-run.
    	Repositories with updates will still be included in the review.
    	--> (default true)
  -gui string
    	The external git GUI application to use for visual reviews.
    	--> (default "smerge")
  -outfile string
    	The path or name of the environment variable containing the
    	path to your pre-existing code review file. If the file exists
    	the final log entry will be appended to that file instead of stdout.
    	--> (default "SMARTY_REVIEW_LOG")
  -repo-list string
    	A colon-separated list of file paths, where each file contains a
    	list of repositories to examine, with one repository on a line.
    	-->
  -review string
    	Letter code of repository statuses to review; where (a) is ahead,
    	origin/master (b) is behind origin/master, (e) has git errors,
    	(f) has new fetched contents, and (m) is messy with uncommitted
    	changes. (j) is like (f) except only 'smartystreets' repositories
    	are considered
    	--> (default "abejm")
  -roots string
    	The name of the environment variable containing colon-separated
    	path values to scan for any git repositories contained therein.
    	Scanning is NOT recursive.
    	NOTE: this flag will be ignored in the case that non-flag command
    	line arguments representing paths to git repositories are provided.
    	--> (default "CDPATH")
```
