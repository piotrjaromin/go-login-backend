Login backend
===
Accounts can be created by either email signup(needed aws account) or facebook login.

Supported database is mongo.


to run sample type
```bash
go run main.go
```

Default configuration is stored in config/config.json.
Security is based on jwt tokens, by calling login with valid credentials you will obtain token which should be used in latter calls to secured endpoints.

Default account is created by createAccount function in main.go

To login
```bash
curl -X POST http://localhost:8080/login -d 'password=test&username=test@test.com'
```

to get account
```bash
curl -X GET http://localhost:8080/accounts/test@test.com -H "Authorization: Bearer $TOKEN"
```

For signups to work you need to set aws credential variables.
Those credentials are needed by email sending fucntionality and they need to have permissions to AWS Email Service.

```bash
AWS_ACCESS_KEY_ID=XXXXXXXXX
AWS_SECRET_ACCESS_KEY=XXXXXXXXX
```

to start signup
```bash
curl -X POST http://localhost:8080/accounts -H "Content-type: application/json" -d '{ "firstName" : "Jhon", "lastName" : "doe", "email" : "kiepur@gmail.com", "password" : "12345aA"  }'
```


to confirm signup
```bash
curl -X PUT http://localhost:8080/accounts/kiepur@gmail.com/confirm -H "Content-type: application/json" -d '{ "code" : "fe4fcfc9-44b6-451e-a608-18d458654bf6" }'
```