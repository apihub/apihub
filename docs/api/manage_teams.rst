===============
Manage My Teams
===============

.. note::

  To perform all the following operations, it's needed to be authenticated. If you do not know how to log in using your Backstage credentials, see :ref:`Log in with Backstage Credentials <login>`.


Creating a new team
-------------------
It's required to inform a valid name. The `alias` attribute is optional. If you do not inform that, the name value will be used to generate the `alias`. This value is important, because you will always use that when making any operation involving teams.

.. note::

  The current user is added to the team automatically as owner.

Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams


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
| alias             |    string    | No                | Yes               |
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

  curl -XPOST -i http://localhost:8000/api/teams -H "Content-Type: application/json" -d '{"name": "Backstage", "alias": "backstage"} ' -H "Authorization: Token EDWZEheeeDnKt0B4IoH8IsOUSnGdumfHmHGQlZDdRbg="


Example Result
==============
.. highlight:: bash

::

  HTTP/1.1 201 Created
  Content-Type: application/json
  Request-Id: aleal.local/BeTVgqw6gS-000004
  Date: Sat, 06 Dec 2014 01:33:05 GMT
  Content-Length: 83

  {"name":"backstage","users":["alice@example.org"],"owner":"alice@example.org"}


If you do not include a valid token in the header, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 401 Unauthorized
  Content-Type: application/json
  Request-Id: aleal.local/BeTVgqw6gS-000002
  Date: Sat, 06 Dec 2014 01:32:14 GMT
  Content-Length: 54

  {"error":"unauthorized_access","error_description":"Request refused or access is not allowed."}


If someone else is using the provided name, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 400 Bad Request
  Content-Type: application/json
  Request-Id: aleal.local/BeTVgqw6gS-000005
  Date: Sat, 06 Dec 2014 01:33:37 GMT
  Content-Length: 90

  {"error":"bad_request","error_description":"Someone already has that team name/alias. Could you try another?"}


Retrieving all teams for the signed user
----------------------------------------

Once you're logged in, it is possible to retrieve all the teams. Backstage takes advantage of the token to identify the user and find the teams.

Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams


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

  curl -XGET -i http://localhost:8000/api/teams -H "Authorization: Token t3Ex657ZSlGrJYnb6-K9vJGvdV9Y0BwrCUambA9_NzQ="


Example Result
==============

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/okpxxUpQ8B-000008
  Date: Sat, 06 Dec 2014 02:33:37 GMT
  Content-Length: 179

  {"items":[{"name":"backstage","alias":"backstage","users":["alice@example.org"],"owner":"alice@example.org"},{"name":"cli","alias":"cli","users":["alice@example.org"],"owner":"alice@example.org"}],"item_count":2}


If the user does not belong to any team, an empty list will be returned:


.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/okpxxUpQ8B-000008
  Date: Sat, 06 Dec 2014 02:35:37 GMT
  Content-Length: 179

  {"items":[],"item_count":0}


Retrieving team info
--------------------

Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/<team-alias>

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

  curl -XGET -i http://localhost:8000/api/teams/backstage -H "Authorization: Token 6rrKX79WwwEnECZMmeYLm8tzSWZmN_mLT7XiFPN14Og="


Example Result
==============

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/okpxxUpQ8B-000008
  Date: Sat, 06 Dec 2014 02:33:37 GMT
  Content-Length: 179

  {"name":"backstage","alias":"backstage","users":["alice@example.org"],"owner":"alice@example.org"}


When trying to retrieve the info for a non-existing team, an error will be returned:

.. highlight:: bash

::

  curl -XGET -i http://localhost:8000/api/teams/non-existing-team -H "Authorization: Token 6rrKX79WwwEnECZMmeYLm8tzSWZmN_mLT7XiFPN14Og="


.. highlight:: bash

::

  HTTP/1.1 404 Not Found
  Content-Type: application/json
  Request-Id: aleal.local/wOPMKpYIfO-000001
  Date: Sat, 06 Dec 2014 01:40:22 GMT
  Content-Length: 47

  {"error":"bad_request","error_description":"Team not found."}


If the team exists, but the user does not belong to it, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Content-Type: application/json
  Request-Id: aleal.local/wOPMKpYIfO-000007
  Date: Sat, 06 Dec 2014 01:42:04 GMT
  Content-Length: 63

  {"error":"access_denied","error_description":"You do not belong to this team!"}


Adding users in the team
------------------------

Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/backstage/users

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

  curl -XPOST -i http://localhost:8000/api/teams/backstage/users -H "Content-Type: application/json" -d '{"users": ["bob@example.org"]}' -H "Authorization: Token 6rrKX79WwwEnECZMmeYLm8tzSWZmN_mLT7XiFPN14Og"


Example Result
==============

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/wOPMKpYIfO-000010
  Date: Sat, 06 Dec 2014 01:44:11 GMT
  Content-Length: 90

  {"name":"backstage","users":["alice@example.org","bob@example.org"],"owner":"alice@example.org"}


If the user does not belong to the team, an error wil be returned:

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Content-Type: application/json
  Request-Id: aleal.local/wOPMKpYIfO-000008
  Date: Sat, 06 Dec 2014 01:43:32 GMT
  Content-Length: 63

  {"error":"access_denied","error_description":"You do not belong to this team!"}


Removing users from team
------------------------

Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/backstage/users

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

  curl -XDELETE -i http://localhost:8000/api/teams/backstage/users -H "Content-Type: application/json" -d '{"users": ["bob@example.org"]}' -H "Authorization: Token vdpazZHBWZCufs-fFaX8teC7Wx1ID5KGTEXRdo3b9vk="


Example Result
==============
.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/IuM9oOVYas-000001
  Date: Sat, 06 Dec 2014 01:47:49 GMT
  Content-Length: 83

  {"name":"backstage","users":["alice@example.org", "bob@example.org"],"owner":"alice@example.org"}


The owner is a special member of the team. And, nobody has permission to remove him from that.

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Content-Type: application/json
  Request-Id: aleal.local/IuM9oOVYas-000005
  Date: Sat, 06 Dec 2014 01:48:59 GMT
  Content-Length: 85

  {"error":"access_denied","error_description":"It is not possible to remove the owner from the team."}


Only members have permission to have another member from the team. If the user does not belong to that, an error will be returned.

.. highlight:: bash

::

  HTTP/1.1 403 Forbidden
  Content-Type: application/json
  Request-Id: aleal.local/IuM9oOVYas-000002
  Date: Sat, 06 Dec 2014 01:48:09 GMT
  Content-Length: 63

  {"error":"access_denied","payload":"You do not belong to this team!"}


Deleting a team
---------------


Resource URL
============
.. highlight:: bash

::

  http://localhost:8000/api/teams/<team-alias>


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

  curl -XDELETE -i http://localhost:8000/api/teams/backstage -H "Authorization: Token 1HnbxXIYMJzECiE-lpH0uIaailRdDurz2JL_5kgtMVc="


Example Result
==============

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/hU8FyyKBPw-000003
  Date: Sat, 06 Dec 2014 01:55:23 GMT
  Content-Length: 58

  {"name":"backstage","users":["alice@example.org","bob@example.org"],"owner":"alice@example.org"}


If the team does not exist, a not found will be returned:

.. highlight:: bash

::

  HTTP/1.1 404 Not Found
  Content-Type: application/json
  Request-Id: aleal.local/hU8FyyKBPw-000004
  Date: Sat, 06 Dec 2014 01:55:33 GMT
  Content-Length: 71

  {"error":"access_denied","error_description":"Team not found or you're not the owner."}
