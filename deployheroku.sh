#!/bin/sh

git checkout -b heroku
godep save ./...
git add -f vendor Godeps
git commit -m "Heroku Deploy"
git push -f heroku HEAD:master
git checkout master
git branch -D heroku
