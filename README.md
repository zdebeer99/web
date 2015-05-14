# webapp

helper functions for web application

current status: WIP

webapp uses a modified gorilla/mux router found in zdebeer99/mux and negroni for managing middleware.

## Middleware

the following middleware is included in webapp.

* Recovery - Copied from negroni
* Logger - Copied from negroni
* MongoDB - Activate a mongodb connection
