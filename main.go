package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/trichner/gitc0ffee/pkg/commit"
	"github.com/trichner/gitc0ffee/pkg/git"
	"github.com/trichner/gitc0ffee/pkg/solver"
	"github.com/trichner/gitc0ffee/pkg/solver/model"
)

const gitCommitObjectType = "commit"

var cli struct {
	UpdateRef bool   `help:"also update the current HEAD revision" default:"false"`
	Prefix    string `help:"a hex prefix to find a collision for" default:"c0ffee"`
	Solver    string `help:"the solver to use for brute-forcing" default:"concurrent"`
}

var solvers = map[string]model.DigestPrefixSolver{
	"concurrent":     solver.NewConcurrentSolver(),
	"singlethreaded": solver.NewSingleThreadedSolver(),
	"native":         solver.NewNativeConcurrentSolver(),
}

func main() {
	kong.Parse(&cli)

	prefix, err := hexStringToBytes(cli.Prefix)
	if err != nil {
		log.Fatal(fmt.Errorf("invalid prefix %q: %w", cli.Prefix, err))
	}

	hashDigest, err := git.GetHeadDigest()
	if err != nil {
		log.Fatal(err)
	}

	contents, err := git.GetCommitContents(hashDigest)
	if err != nil {
		log.Fatal(err)
	}

	obj, err := commit.ParseGitCommitObject(contents)
	if err != nil {
		log.Fatal(err)
	}

	objTemplate, err := solver.PrepareTemplate(obj)
	if err != nil {
		log.Fatal(err)
	}

	s, err := getSolver(cli.Solver)
	if err != nil {
		log.Fatal(err)
	}

	foundSolution, err := bruteForceSolution(objTemplate, prefix, s)
	if err != nil {
		log.Fatal(err)
	}

	err = writeCommitObject(foundSolution)
	if err != nil {
		log.Fatal(err)
	}

	foundDigest := string(foundSolution.Hash[:])
	fmt.Println(foundDigest)

	if cli.UpdateRef {
		if err := updateLastCommit(foundDigest); err != nil {
			log.Fatal(err)
		}
	}
}

func updateLastCommit(digest string) error {
	ref := "HEAD"
	log.Printf("Updating %q to %q", ref, digest)
	if err := git.UpdateReference(ref, digest); err != nil {
		return fmt.Errorf("failed to update branch/ref %q to object %q: %w", ref, digest, err)
	}

	return nil
}

func bruteForceSolution(tpl *model.ObjectTemplate, prefix []byte, s model.DigestPrefixSolver) (*model.CommitObject, error) {
	tick := time.Now()
	foundSolution, err := s.Solve(tpl, prefix)
	if err != nil {
		return nil, fmt.Errorf("cannot find prefix collision: %w", err)
	}
	duration := time.Now().Sub(tick)
	log.Printf("found digest %s in %.2fs", string(foundSolution.Hash[:]), duration.Seconds())
	return foundSolution, nil
}

func getSolver(name string) (model.DigestPrefixSolver, error) {
	s, ok := solvers[name]
	if ok {
		return s, nil
	}
	var availableSolvers []string
	for k := range solvers {
		availableSolvers = append(availableSolvers, k)
	}
	return nil, fmt.Errorf("unknown solver %q, available: %s", name, strings.Join(availableSolvers, ","))
}

func writeCommitObject(obj *model.CommitObject) error {
	writtenDigest, err := git.WriteObject(gitCommitObjectType, obj.Payload)
	if err != nil {
		return fmt.Errorf("failed to write commit object to git store: %w", err)
	}
	if *writtenDigest != *obj.Hash {
		return fmt.Errorf("expected and written git commit object hash don't match: %q != %q", writtenDigest, obj.Hash)
	}
	return nil
}

func hexStringToBytes(s string) ([]byte, error) {
	in := []byte(s)
	if len(in)%2 != 0 {
		return nil, fmt.Errorf("odd length hex encoded bytes: len(%s) = %d", s, len(in))
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(in)/2))

	for i := 0; i < len(in); i += 2 {

		upper, err := hexRuneToByte(in[i])
		if err != nil {
			return nil, err
		}
		lower, err := hexRuneToByte(in[i+1])
		if err != nil {
			return nil, err
		}
		b := upper<<4 | lower
		buf.WriteByte(b)
	}

	return buf.Bytes(), nil
}

func hexRuneToByte(r byte) (byte, error) {
	if r >= '0' && r <= '9' {
		return r - '0', nil
	}
	if r >= 'a' && r <= 'f' {
		return r - 'a' + 10, nil
	}
	return 0, fmt.Errorf("invalid hex rune, expected in [0-9a-f] but was %q", r)
}
