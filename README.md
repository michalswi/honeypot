# SSH honeypot

![](https://img.shields.io/github/issues/michalswi/honeypot)
![](https://img.shields.io/github/forks/michalswi/honeypot)
![](https://img.shields.io/github/stars/michalswi/honeypot)
![](https://img.shields.io/github/last-commit/michalswi/honeypot)

```
ssh-keygen  -t rsa -b 4096 -N "" -f ./private
go run .
```

#### \# example
```
[pts1]$ SSH_PORT=2222 go run .
2023/07/27 19:53:53 SSH server started. Listening on port 2222

[pts2]$ ssh -p 2222 admin@localhost
admin@localhost's password:
Connection to localhost closed.

[pts1]$ SSH_PORT=2222 go run .
2023/07/27 19:53:53 SSH server started. Listening on port 2222
2023/07/27 19:53:59 SSH login attempt [remote addr] from [::1]:49595, username: admin, password: passw0rd, client version: SSH-2.0-OpenSSH_9.1
```

#### \# logs format

`[date] [remote IP addr] [username] [password] [ssh client version]` 
