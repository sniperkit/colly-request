httpcache
=========

[![Build Status](https://travis-ci.org/gregjones/httpcache.svg?branch=master)](https://travis-ci.org/gregjones/httpcache) [![GoDoc](https://godoc.org/github.com/gregjones/httpcache?status.svg)](https://godoc.org/github.com/gregjones/httpcache)

Package httpcache provides a http.RoundTripper implementation that works as a mostly [RFC 7234](https://tools.ietf.org/html/rfc7234) compliant cache for http responses.

It is only suitable for use as a 'private' cache (i.e. for a web-browser or an API-client and not for a shared proxy).

Cache Backends
--------------
- The built-in 'memory' cache stores responses in an in-memory map.

License
-------

-	[MIT License](LICENSE.txt)