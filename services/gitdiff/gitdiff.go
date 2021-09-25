	"time"
	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/analyze"
	"code.gitea.io/gitea/modules/util"
	Match       int
	IsGenerated             bool
	IsVendored              bool
			_, err := io.Copy(io.Discard, reader)
	lastLeftIdx := -1
			lastLeftIdx = -1
			diffLine := &DiffLine{Type: DiffLineAdd, RightIdx: rightLine, Match: -1}
				lastLeftIdx = -1
			}
			if lastLeftIdx > -1 {
				diffLine.Match = lastLeftIdx
				curSection.Lines[lastLeftIdx].Match = len(curSection.Lines)
				lastLeftIdx++
				if lastLeftIdx >= len(curSection.Lines) || curSection.Lines[lastLeftIdx].Type != DiffLineDel {
					lastLeftIdx = -1
				}
			diffLine := &DiffLine{Type: DiffLineDel, LeftIdx: leftLine, Match: -1}
				lastLeftIdx = -1
			}
			if len(curSection.Lines) == 0 || curSection.Lines[len(curSection.Lines)-1].Type != DiffLineDel {
				lastLeftIdx = len(curSection.Lines)
			lastLeftIdx = -1
				count, err := db.Count(m)
					lastLeftIdx = -1
func GetDiffRangeWithWhitespaceBehavior(gitRepo *git.Repository, beforeCommitID, afterCommitID string, maxLines, maxLineCharacters, maxFiles int, whitespaceBehavior string) (*Diff, error) {
	repoPath := gitRepo.Path
	ctx, cancel := context.WithTimeout(git.DefaultContext, time.Duration(setting.Git.Timeout.Default)*time.Second)


	var checker *git.CheckAttributeReader

	if git.CheckGitVersionAtLeast("1.7.8") == nil {
		indexFilename, deleteTemporaryFile, err := gitRepo.ReadTreeToTemporaryIndex(afterCommitID)
		if err == nil {
			defer deleteTemporaryFile()
			workdir, err := os.MkdirTemp("", "empty-work-dir")
			if err != nil {
				log.Error("Unable to create temporary directory: %v", err)
				return nil, err
			}
			defer func() {
				_ = util.RemoveAll(workdir)
			}()

			checker = &git.CheckAttributeReader{
				Attributes: []string{"linguist-vendored", "linguist-generated"},
				Repo:       gitRepo,
				IndexFile:  indexFilename,
				WorkTree:   workdir,
			}
			ctx, cancel := context.WithCancel(git.DefaultContext)
			if err := checker.Init(ctx); err != nil {
				log.Error("Unable to open checker for %s. Error: %v", afterCommitID, err)
			} else {
				go func() {
					err = checker.Run()
					if err != nil && err != ctx.Err() {
						log.Error("Unable to open checker for %s. Error: %v", afterCommitID, err)
					}
					cancel()
				}()
			}
			defer func() {
				cancel()
			}()
		}
	}


		gotVendor := false
		gotGenerated := false
		if checker != nil {
			attrs, err := checker.CheckPath(diffFile.Name)
			if err == nil {
				if vendored, has := attrs["linguist-vendored"]; has {
					if vendored == "set" || vendored == "true" {
						diffFile.IsVendored = true
						gotVendor = true
					} else {
						gotVendor = vendored == "false"
					}
				}
				if generated, has := attrs["linguist-generated"]; has {
					if generated == "set" || generated == "true" {
						diffFile.IsGenerated = true
						gotGenerated = true
					} else {
						gotGenerated = generated == "false"
					}
				}
			} else {
				log.Error("Unexpected error: %v", err)
			}
		}

		if !gotVendor {
			diffFile.IsVendored = analyze.IsVendor(diffFile.Name)
		}
		if !gotGenerated {
			diffFile.IsGenerated = analyze.IsGenerated(diffFile.Name)
		}

func GetDiffCommitWithWhitespaceBehavior(gitRepo *git.Repository, commitID string, maxLines, maxLineCharacters, maxFiles int, whitespaceBehavior string) (*Diff, error) {
	return GetDiffRangeWithWhitespaceBehavior(gitRepo, "", commitID, maxLines, maxLineCharacters, maxFiles, whitespaceBehavior)