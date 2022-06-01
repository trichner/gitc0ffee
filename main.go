package main

import (
	"fmt"
	"github.com/trichner/gitc0ffee/pkg/commit"
	"github.com/trichner/gitc0ffee/pkg/solver"
	"github.com/trichner/gitc0ffee/pkg/template"
	"log"
	"time"
)

func main() {

	prefix := []byte{0xc0, 0xff, 0xee}

	hashDigest, err := getHeadDigest()
	if err != nil {
		log.Fatal(err)
	}

	contents, err := getCommitContents(hashDigest)
	if err != nil {
		log.Fatal(err)
	}

	obj, err := commit.ParseGitCommitObject(contents)
	if err != nil {
		log.Fatal(err)
	}

	objTemplate, err := template.PrepareTemplate(obj)
	if err != nil {
		log.Fatal(err)
	}

	s := solver.NewConcurrentSolver()
	tick := time.Now()
	foundDigest, data, err := s.Solve(objTemplate, prefix)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Now().Sub(tick)
	log.Printf("found digest %s in %.2fs", string(foundDigest[:]), duration.Seconds())

	err = writeGitObject(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(foundDigest[:]))
}
