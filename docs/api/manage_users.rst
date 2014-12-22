=================
Manage My Account
=================

Creating a new user account
---------------------------
To start interacting with Backstage, it's needed to create a user account. Follow below an example, using CURL, of how to create a user account:

.. highlight:: bash

::

  curl -XPOST -i http://localhost:8000/api/users -H "Content-Type: application/json" -d '{"name": "Alice", "email": "alice@example.org", "username": "alice", "password": "123"}'

If your request succeed, the response will be:

.. highlight:: bash

::

  HTTP/1.1 201 Created
  Content-Type: application/json
  Request-Id: aleal.local/E8z3MiQMuT-000001
  Date: Sat, 06 Dec 2014 01:28:31 GMT
  Content-Length: 60

  {"name":"Alice","email":"alice@example.org","username":"alice"}

On the other hand, if there's a validation error, for example: someone else already has that email, the response looks like:

.. highlight:: bash

::

  HTTP/1.1 400 Bad Request
  Content-Type: application/json
  Request-Id: aleal.local/E8z3MiQMuT-000002
  Date: Sat, 06 Dec 2014 01:28:39 GMT
  Content-Length: 90

  {"status_code":400,"payload":"Someone already has that email/username. Could you try another?."


Deleting a user account
-----------------------

The only way to remove a user account is by using a valid Token. For this, it's neeeded to log in with the Backstage credentials to gain a valid Token. If you want to see how to log in, see :ref:`Log in with Backstage Credentials <login>`.

.. warning::

  This action cannot be undone. Once you remove your user, all the teams and services which you are the only member are deleted along with your account.


.. highlight:: bash

::

  curl -i -XDELETE http://localhost:8000/api/users -H "Authorization: Token 1-PYXC0NE5OxrryQ4DmZ_C2WOwAlAOc-uyEKcPW0nr8="

The API returns the resource itself whenever possible. Even after deleting a user, the response payload will be the user:

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/qJJjhtuJc3-000003
  Date: Sat, 06 Dec 2014 01:39:20 GMT
  Content-Length: 59

  {"name":"Alice","email":"alice@example.org","username":"alice"}