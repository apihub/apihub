.. _login:

Authentication
--------------

Most of the API endpoints require authentication. So you need to log in to gain a valid token and then add that in the following requests through the header ``Authentication: TokenType Token``:

.. highlight:: bash

::

  curl -XPOST -i http://localhost:8000/api/login -H "Content-Type: application/json" -d '{"email": "alice@example.org": "password": 123}'

If you have informed the correct credentials, the response will be:

.. highlight:: bash

::

  HTTP/1.1 200 OK
  Content-Type: application/json
  Request-Id: aleal.local/E8z3MiQMuT-000004
  Date: Sat, 06 Dec 2014 01:29:51 GMT
  Content-Length: 144

  {"token":"6rrKX79WwwEnECZMmeYLm8tzSWZmN_mLT7XiFPN14Og=","token_type":"Token","expires":86400,"created_at":"2014-12-06T01:31:11.854062102Z"}%


``token``
~~~~~~~~~
Represents the token itself.

``token_type``
~~~~~~~~~~~~~~
Represents the token type. You need to use it along with the ``token`` value, for example: ``Token 6rrKX79WwwEnECZMmeYLm8tzSWZmN_mLT7XiFPN14Og=``

``expires``
~~~~~~~~~~~
Indicates the time to live of the token, based on the ``created_at`` field.

``created_at``
~~~~~~~~~~~~~~
The exact time when the token was created. Remember to use this value along with ``expires`` to check when the token will be expired.


Otherwise, if the credentials do not match, an error will be returned:

.. highlight:: bash

::

  HTTP/1.1 400 Bad Request
  Content-Type: application/json
  Request-Id: aleal.local/E8z3MiQMuT-000003
  Date: Sat, 06 Dec 2014 01:29:20 GMT
  Content-Length: 61

  {"status_code":400,"message":"Invalid Username or Password."}