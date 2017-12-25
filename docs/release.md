# Noti release

This is the internal process I go through to release a version of Noti. I'm
just writing this down for myself.

## Increment version

* CHANGELOG.md
* README.md

## Run tests

```
cd variadico/noti
make install
make test
```

Make sure CI is green.

## Merge to master

```shell
git checkout master
git merge dev --ff-only

make install
make test

git push origin master

# Last chance to make sure everything looks right, URLs resolve correctly.

git tag 1.2.3
git push origin 1.2.3
```

## Edit GitHub release information

* Click on Releases > 1.2.3 > Edit tag.
* Make the release title 1.2.3.
* Copy and paste the changes from `CHANGELOG.md` into the description box.

Create binaries.

```shell
docker run --rm -it -v $HOME/go:/go golang:latest /bin/bash

cd /go/src/github.com/variadico/noti
make unit-only
make release
```

Upload binaries.

## Eventually update Homebrew

Read this: https://github.com/Homebrew/homebrew-core/blob/master/.github/CONTRIBUTING.md#submit-a-123-version-upgrade-for-the-foo-formula

And this: https://github.com/Homebrew/brew/blob/master/share/doc/homebrew/How-To-Open-a-Homebrew-Pull-Request-(and-get-it-merged).md#create-your-pull-request-from-a-new-branch
