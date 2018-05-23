## API Documentation

**api/signup**
```
user:
	name:
	email:
	id:
	username:
message:
	successfully registered				201 created
	already logged in				200 success
	username exists					409 conflict
	invalid password				400 bad request
	empty name					400 bad request
	empty username					400 bad request
	empty password					400 bad request
	empty email					400 bad request
```

**api/signin**
```
user:
	name:
	email:
	id:
	username:
token:
message:
	successfully logged in			200 success
	already logged in			200 success
	username doesn't exist			404 not found
	invalid password			400 bad request
	empty username                         	400 bad request
        empty password                          400 bad request
```
