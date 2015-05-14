========
Concepts
========

Api/Services
------------
  "An API defines functionalities that are independent of their respective implementations, which allows definitions and implementations to vary without compromising each other. A good API makes it easier to develop a program by providing all the building blocks. A programmer then puts the blocks together."
  From: `Wikipedia <http://en.wikipedia.org/wiki/Application_programming_interface>`_.


App
---
An app is any kind of application - mobile or not, that somehow interact with the apis that are distributed through the gateway.


Backstage Client (cli)
----------------------
An open source command line solution for publishing APIs on Backstage. And, can be found at: `https://github.com/backstage/backstage-client <https://github.com/backstage/backstage-client>`_. It is used to interact with multiple Backstage serves. It's just need to add more than one target. You can find more about it here: `http://godoc.org/github.com/backstage/backstage-client <http://godoc.org/github.com/backstage/backstage-client>`_.


Gateway/Reverse Proxy
---------------------
A reverse proxy sits between the api server and the world. It ensures that all requests for the apis, either for reading or writing, are properly authenticated and authorized. Moreover, it is possible to have several other filters: throttling, share, and so on.


Middleware
----------
Middleware is a wrapper around your API that decorates the requests without adding logic in the application. It's supposed to run before dispatching the request to the API.


OAuth 2.0
---------
  "OAuth 2.0 is the next evolution of the OAuth protocol which was originally created in late 2006. OAuth 2.0 focuses on client developer simplicity while providing specific authorization flows for web applications, desktop applications, mobile phones, and living room devices."
  From: `http://oauth.net/2/ <http://oauth.net/2/>`_.


Restful Api
-----------
Backstage takes advantage of `Json-Schema <http://json-schema.org/>`_ to describe its existing data and include support to `hypermedia <http://en.wikipedia.org/wiki/HATEOAS>`_ to it.


Team
----
A group of users that have valid accounts.


Transformer
-----------
Transformer is supposed to run after the API response, just before writing the final response.

User
----
A user is a person who interacts directly with the application. In our case, a user could be both a developer and an end-user interacting with the services through an application.