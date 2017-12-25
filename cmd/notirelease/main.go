package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/blang/semver"
)

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}

func mainE() error {
	if runtime.GOOS != "darwin" {
		return errors.New("GOOS must be darwin")
	}

	if err := changeRightDir(); err != nil {
		return err
	}

	if err := checkoutBranch("dev"); err != nil {
		return err
	}

	if err := runTests(); err != nil {
		return err
	}

	if err := checkoutBranch("master"); err != nil {
		return err
	}

	if err := mergeToMaster(); err != nil {
		return err
	}

	if err := runTests(); err != nil {
		return err
	}

	if err := pushOrigin("master"); err != nil {
		return err
	}

	tag, err := tagRelease()
	if err != nil {
		return err
	}

	if err := checkoutTag(tag); err != nil {
		return err
	}

	if err := pushOrigin(tag); err != nil {
		return err
	}

	if err := createRelease(); err != nil {
		return err
	}

	return nil
}

func createRelease() error {
	printTitle("RUN createRelease")

	ans, err := promptYesAbortSkip("create release")
	if err != nil {
		return err
	}
	if ans == skipAns {
		printWarning("SKIP createRelease")
		return nil
	}

	if err := runCommand("make", "release"); err != nil {
		return err
	}

	fmt.Println("PASS createRelease")
	return nil
}

func tagRelease() (string, error) {
	printTitle("RUN tagRelease")

	ans, err := promptYesAbortSkip("tag release")
	if err != nil {
		return "", err
	}
	if ans == skipAns {
		printWarning("SKIP tagRelease")
		return "", nil
	}

	var tag string
	for i := 0; i < 3; i++ {
		tag, err = promptString("tag name")
		if err != nil {
			return "", err
		}

		if _, err = semver.Parse(tag); err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}

	curBranch, err := currentBranch()
	if err != nil {
		return "", err
	}
	if curBranch != "master" {
		return "", errors.New("current branch not master")
	}

	if err := gitTag(tag); err != nil {
		return "", err
	}

	fmt.Println("PASS tagRelease")
	return tag, nil
}

func pushOrigin(ref string) error {
	name := fmt.Sprintf("pushOrigin(%s)", ref)
	printTitle("RUN " + name)

	ans, err := promptYesAbortSkip("push to origin")
	if err != nil {
		return err
	}
	if ans == skipAns {
		printWarning("SKIP " + name)
		return nil
	}

	if err := gitPushOrigin(ref); err != nil {
		return err
	}

	fmt.Println("PASS " + name)
	return nil
}

func mergeToMaster() error {
	printTitle("RUN mergeToMaster")

	ans, err := promptYesAbortSkip("merge dev into master")
	if err != nil {
		return err
	}

	if ans == skipAns {
		printWarning("SKIP mergeToMaster")
		return nil
	}

	curBranch, err := currentBranch()
	if err != nil {
		return err
	}
	if curBranch != "master" {
		return errors.New("current branch not master")
	}

	if err := gitMergeFF("dev"); err != nil {
		return err
	}

	fmt.Println("PASS mergeToMaster")
	return nil
}

func runTests() error {
	printTitle("RUN runTests")

	ans, err := promptYesAbortSkip("run tests")
	if err != nil {
		return err
	}

	if ans == skipAns {
		printWarning("SKIP runTests")
		return nil
	}

	if err := runCommand("make", "install"); err != nil {
		return err
	}
	if err := runCommand("make", "test"); err != nil {
		return err
	}

	fmt.Println("PASS runTests")
	return nil
}

func checkoutTag(tag string) error {
	name := fmt.Sprintf("checkoutTag(%s)", tag)
	printTitle("RUN " + name)

	have, err := currentTag()
	if err != nil {
		return err
	}
	want := tag

	if have == want {
		fmt.Println("PASS " + name)
		return nil
	}

	printWarning(fmt.Sprintf("%s\nhave=%s\nwant=%s",
		"Unexpected tag", have, want))

	ans, err := promptYesAbortSkip(fmt.Sprintf("checkout %s tag", tag))
	if err != nil {
		return err
	}

	if ans == skipAns {
		printWarning("SKIP " + name)
		return nil
	}

	if err := gitCheckout(tag); err != nil {
		return err
	}

	fmt.Println("PASS " + name)
	return nil
}

func checkoutBranch(branch string) error {
	name := fmt.Sprintf("checkoutBranch(%s)", branch)
	printTitle("RUN " + name)

	have, err := currentBranch()
	if err != nil {
		return err
	}
	want := branch

	if have == want {
		fmt.Println("PASS " + name)
		return nil
	}

	printWarning(fmt.Sprintf("%s\nhave=%s\nwant=%s",
		"Unexpected branch", have, want))

	ans, err := promptYesAbortSkip(fmt.Sprintf("checkout %s branch", branch))
	if err != nil {
		return err
	}

	if ans == skipAns {
		printWarning("SKIP " + name)
		return nil
	}

	if err := gitCheckout(branch); err != nil {
		return err
	}

	fmt.Println("PASS " + name)
	return nil
}

func changeRightDir() error {
	printTitle("RUN changeRightDir")

	havePwd, err := os.Getwd()
	if err != nil {
		return err
	}
	wantPwd := os.ExpandEnv("$HOME/go/src/github.com/variadico/noti")

	if havePwd == wantPwd {
		fmt.Println("PASS Current dir is noti root")
		return nil
	}

	printWarning(fmt.Sprintf("%s\nhave=%s\nwant=%s",
		"Unexpected pwd", havePwd, wantPwd))

	ans, err := promptYesAbortSkip("cd to noti root")
	if err != nil {
		return err
	}

	if ans == skipAns {
		printWarning("SKIP changeRightDir")
		return nil
	}

	if err := os.Chdir(wantPwd); err != nil {
		return err
	}

	fmt.Println("PASS changeRightDir")
	return nil
}
