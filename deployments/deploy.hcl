project   = "vcs-mock-proxy"
deploy_id = env.BUILD_ID

artifact {
  artifact_type = "quay_repo"

  config {
    docker_file = "build/package/docker/Dockerfile"
    repo        = project.name
    tag         = project.deploy_id
  }
}
