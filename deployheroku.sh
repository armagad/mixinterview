#!/bin/sh

git checkout -b heroku
godep save ./...
git add -f vendor Godeps
git commit -m "Heroku Deploy"
git push heroku master
git checkout master
git branch -D heroku
rm -r vendor Godeps
