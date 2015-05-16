==================
Manage My Clients
==================

.. note::

  A client belongs to a specific team. It's not allowed to create a client without it.


Creating a new client
----------------------
It's required to inform a valid name. If the attribute `id` is not provided, a slug will be generated based on the `name` field to be used as `id`. It's important to mention that the `id` must be unique.


Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/:team/clients


Resource Information
====================

+---------------------------+----------+
| Response formats          |   JSON   |
+---------------------------+----------+
| Requires authentication?  |    Yes   |
+---------------------------+----------+


Payload Parameters
==================
+-------------------+--------------+-------------------+-------------------+
|    Parameter      |     Type     |     Required?     |      Unique?      |
+-------------------+--------------+-------------------+-------------------+
| name              |    string    | Yes               | No                |
+-------------------+--------------+-------------------+-------------------+
| id                |    string    | No                | Yes               |
+-------------------+--------------+-------------------+-------------------+
| secret            |    string    | No                | No                |
+-------------------+--------------+-------------------+-------------------+
| redirect_uri      |    string    | No                | No                |
+-------------------+--------------+-------------------+-------------------+


Header Parameters
=================
+-----------------+--------------+-------------------+
|    Parameter    |     Type     |     Required?     |
+-----------------+--------------+-------------------+
| Authorization   |    string    | Yes               |
+-----------------+--------------+-------------------+


Example Request
===============

.. highlight:: bash

::

  curl -XPOST -i http://localhost:8000/api/teams/backstage/clients -H "Content-Type: application/json" -d '{"name": "Backstage App"}' -H "Authorization: Token hfbXZtQSxQQIAayKVneI8tkeAKHZHgY5JVr03r3YJuI="

Example Result
==============
.. highlight:: bash

::

  HTTP/1.1 201 Created
  Request-Id: aleal.local/uHApWzIKaU-000005
  Date: Sat, 03 Jan 2015 10:50:12 GMT
  Content-Length: 148
  Content-Type: application/json; charset=utf-8

  {"id":"backstage-app","secret":"Ia_6BdHzkey6FF9dn3HeKsMaf_JrOi7kDKQlq-6PZN4=","name":"Backstage App","redirect_uri":"","owner":"alice@example.org","team":"backstage"}

If there is another client using the name provided, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 400 Bad Request
  Request-Id: aleal.local/uHApWzIKaU-000004
  Date: Sat, 03 Jan 2015 10:48:59 GMT
  Content-Length: 85
  Content-Type: application/json; charset=utf-8

  {"error":"bad_request","error_description":"There is another client with this name."}


If the team does not exist, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 404 Not Found
  Request-Id: aleal.local/uHApWzIKaU-000006
  Date: Sat, 03 Jan 2015 10:51:30 GMT
  Content-Length: 61
  Content-Type: application/json; charset=utf-8

  {"error":"not_found","error_description":"Team not found."}

Or, when trying to create a client for a team which you do not belong to, a forbidden error will be returned:

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Request-Id: aleal.local/uHApWzIKaU-000006
  Date: Sat, 03 Jan 2015 10:51:30 GMT
  Content-Type: application/json; charset=utf-8
  Content-Length: 63

  {"error":"access_denied","error_description":"You do not belong to this team!"}