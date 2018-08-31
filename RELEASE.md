# Release the project and generate go release

## The easy way

Simply create a release from the GitHub UI named `release-${SOME_VERSION}`
where `SOME_VERSION` could be for example `0.2`

This will cause CircleCI to perform a release and will create a GitHub release named `v0.2`
that will contain the built binaries for both MacOS and Linux

**Note**: This assumes that the `GITHUB_API_TOKEN` has been set in the CircleCI UI for the job

### Docker images

When performing a release this, the docker image `snowdrop/spring-boot-generator:${version}` will also be published on `quay.io`.
**Note**: This assumes that the `QUAY_ROBOT_USER` and `QUAY_ROBOT_TOKEN` has been set in the CircleCI UI for the job
    

## The manual way

Execute this command within the root of the project where you pass as parameters your `GITHUB_API_TOKEN` and `VERSION` which corresponds to the tag to be created

```yaml
make upload GITHUB_API_TOKEN=YOURTOKEN VERSION=0.3.0
```

**Remark** : You can next edit the release to add a `changelog` using this command `git log --oneline --decorate v0.2.0..v0.3.0`

**Note**: This command assumes that you have all necessary command line dependencies installed