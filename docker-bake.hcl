group "default" {
  targets = ["go-app-check"]
}

target "go-app-check" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["go-docker:test"]
  platforms = ["linux/amd64", "linux/arm64"]
}
