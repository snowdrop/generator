## To deploy using FMP:

- oc apply -f database.yml
- mvn fabric8:deploy -Popenshift

## To deploy using dekorate:

- mvn install -Ddekorate.deploy=true