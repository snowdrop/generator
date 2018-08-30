# HTTP Server generating Spring Boot Zip file

This project will generate a gzip file in response to a HTTP GET request using the parameters of the url to customize a project using a template
and GAV information

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