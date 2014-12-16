[![Build Status](https://travis-ci.org/backstage/backstage.png?branch=master)](https://travis-ci.org/backstage/backstage)

##What is Backstage?

Backstage is an open source api solution for publishing APIs. It's a reverse proxy that sits between your api server and the world.
Several apis require a similar set of features on the backend, such as: authentication, authorization, throttling, analytics and so on. The idea of this project is to provide a simple and easy way to integrate with existing apis and help the developers, so they do not need to implement all of those boilerplate features for each api they may have.

##Quickstart

```bash
  git clone https://github.com/backstage/backstage.git
  cd backstage
  make setup
  make run-api
```

##Running Tests

```bash
  make test
  make race # If you want to check if there's any race condition.
```

## Links:

- Documentation: http://backstage.readthedocs.org/en/latest/
- How to Contribute: http://backstage.readthedocs.org/en/latest/
