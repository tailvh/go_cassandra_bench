# RXLAB GOLANG BENCH CASSANDRA

This project use /empty and /users to benchmark read latency in mem of Cassandra



## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Cassandra, Go

```
brew install ccm maven
pip install cqlsh==5.0.3
ccm create -v 3.9 streamdemoapi
ccm populate -n 1
ccm start
ccm status
echo "CREATE KEYSPACE streamdemoapi WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};" | cqlsh --cqlversion 3.4.2
echo "use streamdemoapi; CREATE TABLE messages (id UUID, user_id UUID, Message text, PRIMARY KEY(id));" | cqlsh --cqlversion 3.4.2
echo "use streamdemoapi; CREATE TABLE users (id UUID, firstname text, lastname text, age int, email text, city text, PRIMARY KEY (id));" | cqlsh --cqlversion 3.4.2

brew install go
export GOPATH=$HOME/golang
export GOROOT=/usr/local/opt/go/libexec
mkdir $HOME/golang
go get github.com/gocql/gocql
go get github.com/valyala/fasthttp

```

### Installing

Finish



## Running the tests

Explain how to run the automated tests for this system

```
run auto unit & end to end test
```



## Deployment

Add additional notes about how to deploy this on a live system



## Built With

* [Cassandra](http://www.nowhere.rxlab/) - Database
* [Golang](http://www.nowhere.rxlab/) - Language




## Contributing

Please read [CONTRIBUTING.md](https://) for details on our code of conduct, and the process for submitting pull requests to us.




## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 




## Authors

* **Rxlab Dev Team** - *Initial work* - [Rxlab Dev Team](https://github.com/)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.




## License

This project is licensed under the Copy Right Rxlab Vn - see the [LICENSE.md](LICENSE.md) file for details