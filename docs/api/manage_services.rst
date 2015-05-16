==================
Manage My Services
==================

.. note::

  A service belongs to a specific team. It's not allowed to create a service without it.


Creating a new service
----------------------
It's required to inform a valid name. The `alias` attribute is optional. If you do not inform that, the name value will be used to generate the `alias`. This value is important, because you will always use that when making any operation involving teams.


Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/:team/services


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
| subdomain         |    string    | Yes               | Yes               |
+-------------------+--------------+-------------------+-------------------+
| description       |    string    | No                | No                |
+-------------------+--------------+-------------------+-------------------+
| disabled          |    boolean   | No                | No                |
+-------------------+--------------+-------------------+-------------------+
| documentation     |    string    | No                | No                |
+-------------------+--------------+-------------------+-------------------+
| endpoint          |    string    | Yes               | No                |
+-------------------+--------------+-------------------+-------------------+
| timeout           |    integer   | No                | No                |
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

  curl -XPOST -i http://localhost:8000/api/teams/backstage/services -H "Content-Type: application/json" -d '{"subdomain": "backstage", "description": "test this", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}' -H "Authorization: Token r-fRrYtDJ0nMAQ3UvHGCZe6ASTal9LXu_PmdyZyGkTM="


Example Result
==============
.. highlight:: bash

::

  HTTP/1.1 201 Created
  Content-Type: application/json
  Request-Id: aleal.local/Iwz0wETBog-000001
  Date: Fri, 05 Dec 2014 19:44:39 GMT
  Content-Length: 309

  {"subdomain":"backstage","created_at":"2014-12-05T17:44:39.462-02:00","updated_at":"2014-12-05T17:44:39.462-02:00","description":"test this","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"alice@example.org","timeout":10,"team": "backstage"}

If any required field is missing, the result will be represented by `400 Bad Request`:

.. highlight:: bash

::

  HTTP/1.1 400 Bad Request
  Content-Type: application/json
  Request-Id: aleal.local/Zh86HQSRtD-000016
  Date: Tue, 23 Dec 2014 17:29:43 GMT
  Content-Length: 47

  {"error":"bad_request","error_description":"Subdomain cannot be empty."}
  or
  {"error":"bad_request","error_description":"Endpoint cannot be empty."}

And when the team is not found:

.. highlight:: bash

::

  HTTP/1.1 404 Not Found
  Content-Type: application/json
  Request-Id: aleal.local/Zh86HQSRtD-000016
  Date: Tue, 23 Dec 2014 17:29:43 GMT
  Content-Length: 47

  {"error":"not_found","error_description":"Team not found."}

Or, when trying to create a service for a service where you do not belong to, you'll get a `403 Forbidden`:

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Content-Type: application/json
  Request-Id: aleal.local/Zh86HQSRtD-000019
  Date: Tue, 23 Dec 2014 17:31:09 GMT
  Content-Length: 63

  {"error":"access_denied","error_description":"You do not belong to this team!"}


Deleting a service
------------------


Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/:team/services/:subdomain


Resource Information
====================

+---------------------------+----------+
| Response formats          |   JSON   |
+---------------------------+----------+
| Requires authentication?  |    Yes   |
+---------------------------+----------+


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

  curl -XDELETE -i http://localhost:8000/api/teams/backstage/services/hello -H "Authorization: Token 1HnbxXIYMJzECiE-lpH0uIaailRdDurz2JL_5kgtMVc="


Example Result
==============

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Request-Id: aleal.local/z7R8abxgq9-000009
  Date: Sat, 03 Jan 2015 10:30:58 GMT
  Content-Length: 237
  Content-Type: application/json; charset=utf-8

  {"subdomain":"hello","description":"test this","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"ringo@gmail.com","team":"backstage","timeout":10}

If the team does not exist, a not found error will be returned:

.. highlight:: bash

::

  HTTP/1.1 404 Not Found
  Content-Type: application/json
  Request-Id: aleal.local/z7R8abxgq9-000007
  Date: Sat, 03 Jan 2015 10:29:29 GMT
  Content-Length: 82

  {"error":"not_found","error_description":"The resource requested does not exist."}