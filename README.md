# webapp

helper functions for web application

current status: WIP

webapp uses a modified gorilla/mux router found in zdebeer99/mux and negroni for managing middleware.

## Overview

API Documentation https://gowalker.org/github.com/zdebeer99/webapp

### Features

To provide a simple one stop library that has the following built-in:
* Session Management
* Easy Access to database
* Basic User Authentication
* Rendering
* Html form and json parsing to schemas.

How it should work
* Use a context structure per request instead of the default http.handlers
    The reason for this is personal preference, and the draw back is that you must change the context struct to tailor it to your needs, the advantage is that you can pass data between middleware layers without a external context, and keep casting to a minimum.
* Middleware [Done] - Using Negroni Embeded
* Routing [Done] - Using gorilla mux
* Session Management
* MongoDB Support
* Data Bindings like form, json, etc

## Basic Web app

## Routing

## Middleware

the following middleware is included in webapp.

* Recovery - Copied from negroni
* Logger - Copied from negroni
* MongoDB - Activate a mongodb connection

**Custom Middleware**

Middleware can be added from the Use(...) and UseFunc(...) functions. Custom Middleware requires a 'c *Context, next HandlerFunc' signature.

Example of middleware requiring user login before opening the page:

```golang
  r:=webapp.New()
  r.UseFunc(func(c *webapp.Context, next webapp.Handler) {
    if c.Auhtenticate(){
      next(c)
    }else{
      c.Redirect("/login")
    }
  })
```

## Rendering

## Context API Reference

**Vars**
map containing route variables.

```golang
app.Get("getItem/{id}",getItem)

func getItem(c *webapp.Context){
  id := c.Vars["id"].(string)
  .
  .
  .
}
```
