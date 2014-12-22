==================
Manage My Services
==================

.. note::

  A service belongs to a specific team. It's not allowed to create a service without it.


Creating a new service
---------------------------
It's required to inform a valid name. The `alias` attribute is optional. If you do not inform that, the name value will be used to generate the `alias`. This value is important, because you will always use that when making any operation involving teams.


.. highlight:: bash

::

  curl -XPOST -i http://localhost:8000/api/services -H "Content-Type: application/json" -d '{"subdomain": "backstage", "allow_keyless_use": true, "description": "test this", "disabled": false, "documentation": "http://www.example.org/doc", "endpoint": "http://github.com/backstage", "timeout": 10}' -H "Authorization: Token 8tSIjilok-6dpuVy3vcosyY5pxq9G776995F4IBHLOw="


.. highlight:: bash

::

  HTTP/1.1 201 Created
  Content-Type: application/json
  Request-Id: aleal.local/Iwz0wETBog-000001
  Date: Fri, 05 Dec 2014 19:44:39 GMT
  Content-Length: 309

  {"subdomain":"backstage","created_at":"2014-12-05T17:44:39.462-02:00","updated_at":"2014-12-05T17:44:39.462-02:00","allow_keyless_use":true,"description":"test this","disabled":false,"documentation":"http://www.example.org/doc","endpoint":"http://github.com/backstage","owner":"nathnb@gmail.com","timeout":10}