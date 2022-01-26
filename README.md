# HTTP Server generating Spring Boot Zip file

This project will generate a zip file in response to a HTTP GET request using the parameters of the url to customize a project using a template
and GAV information.

The format of the request is defined as such

`http://GENERATOR_HOST/app?param1=val1&param2&val2....`

where parameters could be :
  - `{template}` is one of the templates available: crud, custom, rest, ...
  - `{groupId}` is the GAV's group
  - `{artifactId}` is the GAV's artifact
  - `{version}` is the GAV's version
  - `{packageName}` is the name of the Java package
  - `{dependencies}` are the starters/modules that we propose : web, security, jax-rs, ...
  - `{bomVersion}` is the Snowdrop BOM version (1.5.x, ....) 
  - `{springbootVersion}` is the Spring Boot version  (this will override the `bomVersion` parameter for Spring libraries)

## To run it locally

The golang version needed to run locally the generator is `1.11`.
Start the go `generator` server using this command

```bash
CONFIGMAP_PATH=conf go run main.go
time="2018-08-31T16:18:01Z" level=info msg="Starting Spring Boot Generator Server on port 8080 - Version 0.0.5 (4462d0b)"
time="2018-08-31T16:18:01Z" level=info msg="The following REST endpoints are available : "
time="2018-08-31T16:18:01Z" level=info msg="Generate zip : /app"
time="2018-08-31T16:18:01Z" level=info msg="Config : /config"
```

or using the makefile

```bash
make start
>> writing assets
cd /Users/dabou/Code/go-workspace/src/github.com/snowdrop/generator/pkg/template && go generate
writing assets_vfsdata.go
> Build go application
go build -ldflags="-w -X main.GITCOMMIT=77d3113 -X main.VERSION=0.0.666" -o generator main.go
>> Launch generator locally
CONFIGMAP_PATH=conf ./generator
INFO[0000] Log level : info                             
INFO[0000] Parsing Generator's Config at conf/generator.yaml 
INFO[0000] File template : /crud/pom.xml                
INFO[0000] File template : /crud/src/main/java/dummy/CrudApplication.java 
INFO[0000] File template : /crud/src/main/java/dummy/exception/NotFoundException.java 
INFO[0000] File template : /crud/src/main/java/dummy/exception/UnprocessableEntityException.java 
INFO[0000] File template : /crud/src/main/java/dummy/exception/UnsupportedMediaTypeException.java 
INFO[0000] File template : /crud/src/main/java/dummy/service/Fruit.java 
INFO[0000] File template : /crud/src/main/java/dummy/service/FruitController.java 
INFO[0000] File template : /crud/src/main/java/dummy/service/FruitRepository.java 
INFO[0000] File template : /crud/src/main/resources/application-openshift-catalog.properties 
INFO[0000] File template : /crud/src/main/resources/application-openshift.properties 
INFO[0000] File template : /crud/src/main/resources/import.sql 
INFO[0000] File template : /crud/src/main/resources/static/index.html 
INFO[0000] File template : /crud/src/test/java/dummy/BoosterApplicationTest.java 
INFO[0000] File template : /crud/src/test/resources/logback-test.xml 
INFO[0000] File template : /rest/pom.xml                
INFO[0000] File template : /rest/src/main/java/dummy/RestApplication.java 
INFO[0000] File template : /rest/src/main/java/dummy/service/Greeting.java 
INFO[0000] File template : /rest/src/main/java/dummy/service/GreetingEndpoint.java 
INFO[0000] File template : /rest/src/main/resources/application.properties 
INFO[0000] File template : /rest/src/main/resources/static/index.html 
INFO[0000] File template : /rest/src/test/java/dummy/DemoApplicationTest.java 
INFO[0000] File template : /custom/pom.xml              
INFO[0000] File template : /custom/src/main/java/dummy/DemoApplication.java 
INFO[0000] File template : /custom/src/main/resources/application.properties 
INFO[0000] File template : /custom/src/test/java/dummy/DemoApplicationTest.java 
INFO[0000] Starting Spring Boot Generator Server on port 8000 - Version 0.0.666 (77d3113) 
INFO[0000] The following REST endpoints are available :  
INFO[0000] Generate zip : /app                          
INFO[0000] Config : /config                             
```

Next, in a separate terminal window, execute a `curl` or `httpie` request to create a custom project using one of the `modules` or a `template`

```bash
curl http://localhost:8000/app \
  -o demo.zip \
  -d springbootversion=1.5.15.Final \
  -d groupid=com.example \
  -d artifactid=my-spring-boot \
  -d version=0.0.1-SNAPSHOT  \
  -d packagename=com.example.demo
  
http :8000/app \
   groupid==com.example \
   artifactid==demo \
   version==0.0.1-SNAPSHOT \
   packagename==com.example.demo \
   springbootversion==1.5.15.RELEASE \
   module==web \
   module==keycloak \
   template==custom > demo.zip  
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
