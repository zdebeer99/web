/*

user handles and control user access.

Basic Authentication.

Features
---------
* Login and remembered loggedin user
* Logout
* Authentication
* LoginRequiredHandler

Paths
------
* Login Url

Workflow 1
----------

1. Login
2. Check if user is valid for every request
3. if user fails redirect to login
*/
package user
