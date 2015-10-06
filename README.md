[![Build Status](https://travis-ci.org/apihub/apihub.png?branch=master)](https://travis-ci.org/apihub/apihub)


##What is ApiHub?

ApiHub is an open source solution for publishing APIs. It's a reverse proxy that sits between your api server and the world.
Several apis require a similar set of features on the backend, such as: authentication, authorization, throttling, analytics and so on. The idea of this project is to provide a simple and easy way to integrate with existing apis and help the developers, so they do not need to implement all of those boilerplate features for each api they may have.

##Quickstart

```bash
  git clone https://github.com/apihub/apihub.git
  cd apihub
  make setup
  make run-api-sample
```

##Running Tests

```bash
  make test
  make race # If you want to check if there's any race condition.
```

## Links:

- Documentation: https://apihub.readme.io/

##Contributing

Please refer to the documentation: https://apihub.readme.io/docs/guidelines

##License

Distributed under the New BSD License. See LICENSE file for further details.