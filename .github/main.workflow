workflow "Test and Build" {
  on = "push"
  resolves = ["Upload Coverage", "Upload Builds"]
}

action "Test" {
  uses = "docker://golang:1.11-stretch"
  args = ["go", "test", "-coverprofile", "coverage.out", "./..."]
}

action "Upload Coverage" {
  needs = "Test"
  uses = "actions/bin/sh@master"
  args = [
    "export GIT_COMMIT_SHA=$GITHUB_SHA",
    "export GIT_BRANCH=${GITHUB_REF#refs/heads/}",
    "apt-get update && apt-get install -y curl git",
    "curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > /tmp/cc-test-reporter",
    "chmod +x /tmp/cc-test-reporter",
    "cat coverage.out | sed -e 's|go.htdvisser.nl/did/||g' > coverage.stripped.out",
    "/tmp/cc-test-reporter format-coverage -t gocov -o codeclimate.json coverage.stripped.out",
    "/tmp/cc-test-reporter upload-coverage -i codeclimate.json"
  ]
  secrets = ["CC_TEST_REPORTER_ID"]
}

action "Build for Linux" {
  needs = "Test"
  uses = "docker://golang:1.11-stretch"
  env = {
    GOOS   = "linux"
    GORACH = "amd64"
  }
  args = ["go", "build", "-o", "did-linux-amd64", "./cmd/did"]
}

action "Build for Mac" {
  needs = "Test"
  uses = "docker://golang:1.11-stretch"
  env = {
    GOOS  = "darwin"
    GORACH = "amd64"
  }
  args = ["go", "build", "-o", "did-darwin-amd64", "./cmd/did"]
}

action "Build for Windows" {
  needs = "Test"
  uses = "docker://golang:1.11-stretch"
  env = {
    GOOS  = "windows"
    GORACH = "amd64"
  }
  args = ["go", "build", "-o", "did-windows-amd64.exe", "./cmd/did"]
}

action "Upload Builds" {
  needs = ["Build for Linux", "Build for Mac", "Build for Windows"]
  uses = "actions/bin/sh@master"
  args = [
    "ls -al"
  ]
}
