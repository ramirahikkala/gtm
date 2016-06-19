package metric

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/git-time-metric/gtm/project"
	"github.com/git-time-metric/gtm/scm"
	"github.com/git-time-metric/gtm/util"
)

func TestFullCommit(t *testing.T) {
	repo := util.NewTestRepo(t, false)
	defer repo.Remove()

	curDir, err := os.Getwd()
	util.CheckFatal(t, err)
	defer os.Chdir(curDir)

	os.Chdir(repo.PathIn(""))

	repo.SaveFile("event.go", "event", "")
	repo.SaveFile("event_test.go", "event", "")
	repo.SaveFile("1458496803.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496811.event", project.GTMDir, filepath.Join("event", "event_test.go"))
	repo.SaveFile("1458496818.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496943.event", project.GTMDir, filepath.Join("event", "event.go"))

	treeID := repo.Stage(filepath.Join("event", "event.go"), filepath.Join("event", "event_test.go"))
	commitID := repo.Commit(treeID)

	_, err = Process(false)
	if err != nil {
		t.Fatalf("Process(false) - test full commit, want error nil, got %s", err)
	}

	n, err := scm.ReadNote(commitID.String(), "gtm-data")
	util.CheckFatal(t, err)

	want := []string{`total:300.*`, `event.go:280.*,m`, `event_test.go:20.*,m`}
	for _, s := range want {
		matched, err := regexp.MatchString(s, n.Note)
		util.CheckFatal(t, err)
		if !matched {
			t.Errorf("Process(false) - test full commit, \nwant:\n%s,\ngot:\n%s", s, n.Note)
		}

	}
}

func TestPartialCommit(t *testing.T) {
	repo := util.NewTestRepo(t, false)
	defer repo.Remove()

	curDir, err := os.Getwd()
	util.CheckFatal(t, err)
	defer os.Chdir(curDir)

	os.Chdir(repo.PathIn(""))

	repo.SaveFile("event.go", "event", "")
	repo.SaveFile("event_test.go", "event", "")
	treeID := repo.Stage(filepath.Join("event", "event.go"), filepath.Join("event", "event_test.go"))
	commitID := repo.Commit(treeID)

	repo.SaveFile("event_test.go", "event", "update")
	repo.SaveFile("1458496803.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496811.event", project.GTMDir, filepath.Join("event", "event_test.go"))
	repo.SaveFile("1458496818.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496943.event", project.GTMDir, filepath.Join("event", "event.go"))

	treeID = repo.Stage(filepath.Join("event", "event_test.go"))
	commitID = repo.Commit(treeID)

	_, err = Process(false)
	if err != nil {
		t.Fatalf("Process(false) - test full commit, want error nil, got %s", err)
	}

	n, err := scm.ReadNote(commitID.String(), "gtm-data")
	util.CheckFatal(t, err)

	want := []string{`total:300`, `event_test.go:20.*,m`, `event.go:280.*,r`}
	for _, s := range want {
		matched, err := regexp.MatchString(s, n.Note)
		util.CheckFatal(t, err)
		if !matched {
			t.Errorf("Process(false) - test partial commit, \nwant:\n%s,\ngot:\n%s", s, n.Note)
		}
	}
}

func TestInterim(t *testing.T) {
	repo := util.NewTestRepo(t, false)
	defer repo.Remove()

	curDir, err := os.Getwd()
	util.CheckFatal(t, err)
	defer os.Chdir(curDir)

	os.Chdir(repo.PathIn(""))

	repo.SaveFile("event.go", "event", "")
	repo.SaveFile("event_test.go", "event", "")
	repo.SaveFile("1458496803.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496811.event", project.GTMDir, filepath.Join("event", "event_test.go"))
	repo.SaveFile("1458496818.event", project.GTMDir, filepath.Join("event", "event.go"))
	repo.SaveFile("1458496943.event", project.GTMDir, filepath.Join("event", "event.go"))

	treeID := repo.Stage(filepath.Join("event", "event.go"), filepath.Join("event", "event_test.go"))
	commitID := repo.Commit(treeID)

	commitNote, err := Process(true)
	if err != nil {
		t.Fatalf("Process(false) - test full commit, want error nil, got %s", err)
	}

	n, err := scm.ReadNote(commitID.String(), "gtm-data")
	util.CheckFatal(t, err)

	if n.Note != "" {
		t.Errorf("Process(true) - test interm, notes is note blank, \n%s\n", n.Note)
	}

	if commitNote.Total() != 300 {
		t.Errorf("Process(true) - test interm, want total 300, got %d", commitNote.Total())
	}
}
