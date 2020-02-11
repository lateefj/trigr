=====
trigr
=====

.. image:: https://c1.staticflickr.com/1/246/449673647_4d9a1b900e.jpg


Who needs another automation server?
------------------------------------

TLDR: What if there was a automation server that was distributed? Moving from central source control from csv / svn to git / mecurial.

Waiting for remote build fails suck! It would be nice to use the extra resources sitting on the average developer machine to automatically compile, test and package applications while sitting in a meeting or at lunch. Trigr is mainly an event dispatch system. File changes, commits ect. The goal is to be able to run code based on a set of events. The curated DSL and API should make build and deploy scripting easier. If the built in extensions / DSL is used it should work the same on a local development machine as a remote build server.

Add automation project
----------------------

When first starting the server there is not configuration. This will produce a message. However it is easy to add a project with a simple curl command. 
.. code-bloack:: bash

   curl "http://localhost:8080/project/new?id=trigr&path=/my_home/lhj/my_codes_path/trigr"

To save the project in the configuration file add a persist flag to the end of the url.

.. code-bloack:: bash

   curl "http://localhost:8080/project/new?id=trigr&path=/my_home/lhj/my_code_paths/trigr&persist=true"

To watch the output of the automation use the trigr command line tool.

.. code-bloack:: bash
   ./build/darwin/trigr trigr -tlog
  
   connecting to ws://localhost:8080/ws/trigr
   51935-10-13 23:07:02 -0800 PST ➜ running: go test -cover
   2019-12-19 07:03:20 -0800 PST ➜ PASS
   coverage: 13.1% of statements
   ok      github.com/lateefj/trigr/cmd/trigd      0.611s

This will continue to stream changes to standard out.


Features (TODO)
---------------

* Decentralized: able to run on a development system (laptop) or build server
* Automated Configuration: getting software packaged for use is often a Rube Goldberg process that is only made worse by user interfaces 
* Modern Interface: a web ui not for configuration but for monitoring builds 


Development Process
-------------------

trigr assumes a development process that is a series of connected pipelines configuration, prepare, build, package and deploy. 

* configuration

  * code checkout
  * dependencies
* prepare

  * code generation
* build

  * compile
  * obscure (javascript)
* package

  * tar, compression 
  * encryption
  * distribution
  * docker
* deploy


