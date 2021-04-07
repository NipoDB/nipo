# How to Contributing ?
1. Fork nipo project to your repository (https://github.com/NipoDB/nipo)
2. Go to your **$GOPATH** via `cd $GOPATH`
3. If your not having **bin,pkg,src** folder creating using `mkdir -p {src,bin,pkg}`
4. Create nipo module original path `mkdir -p src/github.com/NipoDB`
5. Go to **NipoDB** path `cd src/github.com/NipoDB`
6. Clone Forked Project in current path `git clone https://github.com/your_username/nipo.git`
7. Go to **nipo** folder `cd nipo`
8. Change your branch to **dev** with `git checkout dev`
9. Getting required packages via `go get ./...`
10. Add remote `git remote add upstream https://github.com/NipoDB/nipo.git` to update your fork repo.
11. Working on Project and send yours **PR**

# Note
- To check and update your forked repo before any changes use 
  - `git fetch upstream` 
  - `git merge upstream/dev` 
  - `git merge upstream/main`
  - If changed somethings via **merge**, push on your forked branch `git push origin dev` or `git push origin main` for updating your forked repo.
 
- Work just on **dev** branch.
- After change in any file **add/commit** your changes.
