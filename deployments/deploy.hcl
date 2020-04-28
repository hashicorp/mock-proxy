# The "vcs-mock-proxy" project name refers to HashiCorp's internal use of this
# tool (mocking version control software providers). You may see scattered use
# of this name throughout the project, but should ignore it in favor of just
# mock-proxy.
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
