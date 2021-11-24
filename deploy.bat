@echo off
git add cmd/main.go cmd/config.go internal/ web/ go.mod  && git commit -m "fix" && git push heroku master