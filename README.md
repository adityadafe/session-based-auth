Requests:

healthcheck:

GET http://localhost:8080/health

usercreate:

POST http://localhost:8080/sign-up
{"username":"user","password":"password"}

userauth:

POST http://localhost:8080/sign-in
{"username":"user","password":"password"}

usersignout:

POST http://localhost:8080/sign-out

protectedRoute:

GET http://localhost:8080/protected-route
