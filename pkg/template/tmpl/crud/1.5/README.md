## To deploy using FMP:

- oc apply -f database.yml
- mvn fabric8:deploy -Popenshift

## To deploy using ap4k:

- mvn install -Dap4k.deploy=true