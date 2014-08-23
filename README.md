swugger
======

This proof-of-concept repo is just meant as to start a discussion about swagger support for httprouter (and others) based on the swagger support in go-restful/swagger.

In the near future I'd like to create a new repo, forked from go-restful, and stripped back to provide swagger support for other routers.

Background
----------

 1. [go-restful](https://github.com/emicklei/go-restful) has support for self-documenting APIs, with [Swagger](https://helloreverb.com/developers/swagger) support out of the box. You can even grab a copy of [swagger-ui](https://github.com/wordnik/swagger-ui) and tell go-restful to serve it. go-restful encourages you to add documentation at the point where you set up routes, which makes for up-to-date documentation. It uses fluent interfaces, which is nice but I don't know if it's 'idiomatic'.
 2. [httprouter](https://github.com/julienschmidt/httprouter) is a high-performance router, which doesnt currently have swagger support. This is true of most of the routers represented in [this benchmark](https://github.com/julienschmidt/go-http-routing-benchmark), so it might be nice to make swagger functionality available to each of these frameworks.
 3. There are some other go-based swagger offerings, but I got the best mileage from go-restful. Please let me know of any others.

swugger
-------
Swugger is a quick and dirty way to show how you might add swagger support to httprouter.

 * AFAIK the only way to 'borrow' go-restful's swagger code, is to set up & document 'dummy' routes based on equivalent httprouter routes. 
 * AFAIK the only way to record routing information in httprouter, is to wrap routing requests inside a 'proxy' function.

So, I just made a package 'swugger' to apply these workarounds and use them in an example.

Installation
------------

	go get github.com/laher/swugger/examples/swugger-hello-httprouter


Running the example
-------------------

	swugger-hello-httprouter

This runs a service on localhost:8080. 

You can browse to http://localhost:8080/doc/apidocs.json and see the main swagger representation. You can then append /hello to list the methods in the example service itself.


Running swagger-ui
------------------

(Approximately, depending on where your gopath is:

	git clone https://github.com/wordkik/swagger-ui

Now run the example again such that it will pick up swagger-ui:

	cd swagger-ui
	swugger-hello-httprouter

Now browse to http://localhost:8080/doc/apidocs/ and put /doc/apidocs.json into the text box, & hit 'Explore'. 
You can then expand the list of operations & invoke them accordingly.

Next steps
----------
I'll go ahead and start modifying a fork of go-restful/swagger, to do the same thing but without actually setting up dummy routes. But first I'll start a discussion with the creators of both projects.


