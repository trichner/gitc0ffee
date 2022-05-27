package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

//gitObjectPath (x:y:zs) = ".git" </> "objects" </> [x, y] </> zs

const gitObjectBasePath = ".git/objects"

const hextable = "0123456789abcdef"

type objectDigest [20]byte
type hexObjectDigest [40]byte

var prefixC0ffee = objectDigest([20]byte{0xc0, 0xff, 0xee})
var maskC0ffee = objectDigest([20]byte{0xff, 0xff, 0xff})

var prefixZero = objectDigest([20]byte{})
var maskZero = objectDigest([20]byte{0xff})

type commitObject struct {
	Bytes []byte
	Salt  []byte
}

func main() {
	//var digest hexObjectDigest
	//copy(digest[:], commitObjectDigest)

	digest, err := getHeadDigest()
	if err != nil {
		log.Fatal(err)
	}

	obj, err := prepareCommitContents(digest)
	if err != nil {
		log.Fatal(err)
	}
	foundDigest, err := bruteForcePrefix(obj)
	log.Printf("found digest: %s", string(foundDigest[:]))
	log.Println(string(obj.Bytes))

	err = writeGitObject(obj.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(foundDigest[:]))
}

func bruteForcePrefix(obj *commitObject) (*hexObjectDigest, error) {

	tick := time.Now()
	for salt := uint64(0); salt < math.MaxUint64; salt++ {
		hexEncodeUint64(obj.Salt, salt)
		digest := hashObject(obj.Bytes)

		if digest[0] == 0xc0 && digest[1] == 0xff && digest[2] == 0xee {
			//if digest[0] == 0xc0 && digest[1] == 0xfe && digest[2] == 0xba && digest[3] == 0xbe {
			tock := time.Now()
			d := tock.Sub(tick)
			rate := float64(salt) / d.Seconds()
			log.Printf("found in %.4f seconds at %.4f hashes second", d.Seconds(), rate)

			var hexDigest hexObjectDigest
			hexEncodeDigest(&hexDigest, digest)
			return &hexDigest, nil
		}
		if salt&0xffffff == 0xffffff {
			tock := time.Now()
			d := tock.Sub(tick)
			rate := float64(salt) / d.Seconds() / 1000
			log.Printf("brute forcing at %.4f khash/s", rate)
		}
	}

	return nil, fmt.Errorf("exhausted salts, nothing found")
}

func prepareCommitContents(hexDigest *hexObjectDigest) (*commitObject, error) {
	contents, err := getCommitContents(hexDigest)
	if err != nil {
		return nil, err
	}

	const hexUint64Len = 64 / 8 * 2

	bodyLength := len(contents) + hexUint64Len

	objectPrefix := fmt.Sprintf("commit %d\x00", bodyLength)
	objectPrefixLen := len(objectPrefix)

	objectLength := objectPrefixLen + bodyLength

	newContents := make([]byte, objectLength)
	obj := &commitObject{
		Bytes: newContents,
		Salt:  newContents[len(newContents)-hexUint64Len:],
	}

	copy(obj.Bytes[objectPrefixLen:], contents)
	copy(obj.Bytes, objectPrefix)
	hexEncodeUint64(obj.Salt, 0xbadc0de)
	return obj, nil
}

func writeGitObject(contents []byte) error {

	digest := hashObject(contents)

	var hexDigest hexObjectDigest
	hexEncodeDigest(&hexDigest, digest)

	f, err := OpenGitObjectFile(&hexDigest)
	if err != nil {
		return fmt.Errorf("cannot open git object %q: %w", digest, err)
	}
	defer f.Close()

	writer := zlib.NewWriter(f)
	_, err = writer.Write(contents)
	if err != nil {
		return fmt.Errorf("cannot compress git object %q: %w", digest, err)
	}
	defer writer.Close()

	return nil
}

func OpenGitObjectFile(digest *hexObjectDigest) (io.WriteCloser, error) {

	path := getGitObjectPath(digest)

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return nil, fmt.Errorf("cannot create git object folder %q: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	return f, err
}

func hashObject(contents []byte) *objectDigest {
	byteDigest := objectDigest(sha1.Sum(contents))
	return &byteDigest
}

func hexEncodeDigest(dst *hexObjectDigest, src *objectDigest) int {
	j := 0
	for i := 0; i < len(src); i++ {
		dst[j] = hextable[src[i]>>4]
		dst[j+1] = hextable[src[i]&0x0f]
		j += 2
	}
	return len(src) * 2
}

func hexEncodeUint64(dst []byte, src uint64) {

	// dst must be at least 16 bytes long

	dst[0] = hextable[(src>>4)&0x0f]
	dst[1] = hextable[src&0x0f]

	dst[2] = hextable[(src>>12)&0x0f]
	dst[3] = hextable[(src>>8)&0x0f]

	dst[4] = hextable[(src>>20)&0x0f]
	dst[5] = hextable[(src>>16)&0x0f]

	dst[6] = hextable[(src>>28)&0x0f]
	dst[7] = hextable[(src>>24)&0x0f]

	dst[8] = hextable[(src>>36)&0x0f]
	dst[9] = hextable[(src>>32)&0x0f]

	dst[10] = hextable[(src>>44)&0x0f]
	dst[11] = hextable[(src>>40)&0x0f]

	dst[12] = hextable[(src>>52)&0x0f]
	dst[13] = hextable[(src>>48)&0x0f]

	dst[14] = hextable[(src>>60)&0x0f]
	dst[15] = hextable[(src>>56)&0x0f]
}

func hexDecodeDigest(dst *objectDigest, src *hexObjectDigest) error {
	_, err := hex.Decode(dst[:], src[:])
	if err != nil {
		return fmt.Errorf("cannot decode hex string %q", src)
	}
	return nil
}

// https://matthew-brett.github.io/curious-git/reading_git_objects.html
func readGitObject(hexDigest *hexObjectDigest) ([]byte, error) {

	p := getGitObjectPath(hexDigest)

	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("cannot read object %q: %w", hexDigest, err)
	}

	r, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("cannot decompress object %q: %w", hexDigest, err)
	}
	defer r.Close()

	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot decompress object %q: %w", hexDigest, err)
	}
	return contents, nil
}

func getGitObjectPath(digest *hexObjectDigest) string {
	return fmt.Sprintf("%s/%s/%s", gitObjectBasePath, digest[:2], digest[2:])
}
