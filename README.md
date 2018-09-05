[![CircleCI](https://circleci.com/gh/snowdrop/generator/tree/master.svg?style=shield)](https://circleci.com/gh/snowdrop/generator/tree/master)

https://appdev.openshift.io/docs/spring-boot-runtime.html#mission-http-api-spring-boot

# HTTP Server generating Spring Boot Zip file

This project will generate a zip file in response to a HTTP GET request using the parameters of the url to customize a project using a template
and GAV information.

The format of the request is defined as such

`http://GENERATOR_HOST/app?param1=val1&param2&val2....`

where parameters could be :
  - `{template}` is one of the templates available: crud, simple, rest, ...
  - `{groupId}` is the GAV's group
  - `{artifactId}` is the GAV's artifact
  - `{version}` is the GAV's version
  - `{packageName}` is the name of the Java package
  - `{dependencies}` are the starters/modules that we propose : web, security, jax-rs, ...
  - `{bomVersion}` is the Snowdrop BOM version (1.5.x, ....) 
  - `{springbootVersion}` is the Spring Boot version  (this will override the `bomVersion` parameter for Spring libraries)

## To run it locally

Start the go generator server

```bash
CONFIGMAP_PATH=conf go run main.go
time="2018-08-31T16:18:01Z" level=info msg="Starting Spring Boot Generator Server on port 8080 - Version 0.0.5 (4462d0b)"
time="2018-08-31T16:18:01Z" level=info msg="The following REST endpoints are available : "
time="2018-08-31T16:18:01Z" level=info msg="Generate zip : /app"
time="2018-08-31T16:18:01Z" level=info msg="Config : /config"
```

Next, in a separate terminal window, execute a `curl` or `httpie` request

```bash
curl http://localhost:8000/app \
  -o demo.zip \
  -d bomVersion=1.5.15.Final \
  -d groupId=com.example \
  -d artifactId=my-spring-boot \
  -d version=1.0  \
  -d packageName=com.example.demo \
  -d springbootVersion=1.5.15.RELEASE
  
http :8000/app \
  bomVersion==1.5.15.Final \
  groupId==com.example \
  artifactId==my-spring-boot \
  version==1.0  \
  packageName==com.example.demo \
  springbootVersion==1.5.15.RELEASE > demo.zip  
``` 

Unzip the file

```bash
unzip demo.zip
```

## To build the Server as container's image

Execute this command at the root of the project
```bash
imagebuilder -t spring-boot-generator:latest -f docker/Dockerfile_generator .
```

**WARNING** : The `imagebuilder` tool supports to process multi-stages Dockerfile and can be installed using the following go command `go get -u github.com/openshift/imagebuilder/cmd/imagebuilder`

Tag the docker image and push it to `quay.io`

```bash
TAG_ID=$(docker images -q spring-boot-generator:latest)
docker tag $TAG_ID quay.io/snowdrop/spring-boot-generator
docker login quai.io
docker push quay.io/snowdrop/spring-boot-generator
```

## To Deploy it

To test the docker image, execute this command to create a pod/service and route
```bash
oc new-project generator
oc create -f docker/generator-application.yml
```
