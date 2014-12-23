========
Overview
========

Backstage is an open source solution for publishing APIs. It's a reverse proxy that sits between your api server and the world.
Several apis require a similar set of features on the backend, such as: authentication, authorization, throttling, analytics and so on. The idea of this project is to provide a simple and easy way to integrate with existing apis and help the developers, so they do not need to implement all of those boilerplate features for each api they may have.


Why Backstage?
==============
It is an open source project that consists of several modules: a restful api, a gateway and a cli interface. In addition to being highly scalable and easily extensible through middleware/filters. All services that are distributed through the gateway have the assurance that the incoming requests have been properly authenticated and/or authorized. In addition, you can create filters that manipulate the request, either adding or removing headers, for example.


Backstage Client
================
There's a `command line <https://github.com/backstage/backstage-client>`_ solution for interacting with Backstage Servers. The documentation can be found at: `<https://godoc.org/github.com/backstage/backstage-client/backstage>`_.