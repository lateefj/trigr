=====
trigr
=====

.. image:: https://c1.staticflickr.com/1/246/449673647_4d9a1b900e.jpg


Who needs another build system?
-------------------------------

TLDR: What if there was a build system that was distributed? Moving from central source control from csv / svn to git / mecurial.

Waiting for remote build fails suck! It would be nice to use the extra resources sitting on the average developer machine to automatically compile, test and package applications while sitting in a meeting or at lunch.


Features 
--------

* Decentralized: able to run on a development system (laptop)
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


