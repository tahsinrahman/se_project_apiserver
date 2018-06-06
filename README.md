## API Documentation

**api/signup**

Request
```
{
    "name" : "__",
    "username" : "__",
    "email" : "__",
    "password" : "__"
}
```

Response
```
{
    "user" {
        "id" : "__",
        "name" : "__",
        "username" : "__",
        "email" : "__"
    },
    "message" : "__"
}
```

**api/signin**

Request
```
{
    "username" : "__",
    "password" : "__"
}
```

Response
```
{
    "user" {
        "id" : "__",
        "name" : "__",
        "username" : "__",
        "email" : "__",
    },
    "token" : "__",
    "message" : "__"
}
```

**api/create-company**

Request
```
{
    "name" : "__",
    "description": "__"
}
```

Response
```
{
    "user" {
        "id" : "__",
        "name" : "__",
        "username" : "__",
        "email" : "__",
    },
    "company" {
        "id"
        "name"
        "description"
        "admin" : [
            {
                "ID":
                "name":
                "username":
                "email":
            }
            {
                "ID":
                "name":
                "username":
                "email":
            }
            .
            .
            .
        ]
        "HR" : [
            {
                "ID":
                "name":
                "username":
                "email":
            }
            {
                "ID":
                "name":
                "username":
                "email":
            }
            .
            .
            .
        ]
    },
    "message" : "__"
}
```

**api/id/update-company**

request
```
{
    "name" :
    "description" :
    "admin" : [
        {
            "ID":
            "name":
            "username":
            "email":
        }
        {
            "ID":
            "name":
            "username":
            "email":
        }
        .
        .
        .
    ]
    "HR" : [
        {
            "ID":
            "name":
            "username":
            "email":
        }
        {
            "ID":
            "name":
            "username":
            "email":
        }
        .
        .
        .
    ]
}
```

response
```
{
    "user" {
        "id" : "__",
        "name" : "__",
        "username" : "__",
        "email" : "__",
    },
    "company" {
        "id"
        "name"
        "description"
    },
    "message" : "__"
}
```
