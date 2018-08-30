# HTTP Server generating Spring Boot Zip file

This project will generate a zip file in response to a HTTP GET request using the parameters of the url to customize a project using a template
and GAV information.

The format of the request is defined as such

`http://GENERATOR_HOST/template/{id}?param1=val1&param2&val2....`

where :
  - `{id}` is one of the templates available: crud, simple, rest, ...
  - `{groupId}` is the GAV's group
  - `{artifactId}` is the GAV's artifact
  - `{version}` is the GAV's version
  - `{packageName}` is the name of the Java package
  - `{dependencies}` are the starters/modules that we propose : web, security, jax-rs, ...
  - `{bomVersion}` is the Spring Boot BOM version (1.5.x, ....) 
  - `{springbootVersion}` is the Snowdrop Bom version 

## To run it locally

Start the go generator server

```bash
GENERATOR_PATH=pkg/scaffold go run main.go
INFO[0000] Starting Spring Boot Generator Server on port 8000, exposing endpoint /template/{id}` - Version : v0.0.0 (HEAD) 
```

Next, in a separate terminal window, execute a `curl` or `httpie` request

```bash
curl http://localhost:8000/template/crud \
  -o demo.zip \
  -d bomVersion=1.5.15.Final \
  -d groupId=com.example \
  -d artifactId=my-spring-boot \
  -d version=1.0  \
  -d packageName=com.example.demo \
  -d springbootVersion=1.5.15.RELEASE
  
http :8000/template/crud \
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
oc create -f docker/generator.yml
```

## To test it
  

Grab the route of the service and next generate a project
```bash
sb create -t simple -i ocool -g org.acme.cool -v 1.0.0 -p me.snowdrop.cool -s 1.5.15.RELEASE -b 1.5.15.Final -u http://spring-boot-generator.192.168.65.2.nip.io/
sb create -t simple -d web -i ocool -g org.acme.cool -v 1.0.0 -p me.snowdrop.cool -s 1.5.15.RELEASE -b 1.5.15.Final -u http://spring-boot-generator.192.168.65.2.nip.io/
sb create -t simple -d cxf -i ocool -g org.acme.cool -v 1.0.0 -p me.snowdrop.cool -s 1.5.15.RELEASE -b 1.5.15.Final -u http://spring-boot-generator.192.168.65.2.nip.io/

curl http://spring-boot-generator.195.201.87.126.nip.io/template/crud?artifactId=my-spring-boot&bomVersion=1.5.15.Final&groupId=com.example&outDir=&packageName=com.example.demo&springbootVersion=1.5.15.RELEASE&version=1.0 


```