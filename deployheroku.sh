#!/bin/sh

git checkout -b heroku
godep save ./...
git add -f vendor Godeps
git commit -m "Heroku Deploy"
git push -f HEAD heroku
git checkout master
git branch -D heroku
rm -r vendor Godeps
