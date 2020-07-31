# NullServe Origin Service

This is an HTTP web server and dynamic proxy which routes content based on host name to the correct serverless backend.
For Origin Service, hostname dictates which application an end user is requesting, allowing the NullServe Origin Service to locate application resources, including static content proxied APIs and serverless APIs.
The deployed application controls its own routing scheme and instructs how the Origin Service should respond.
