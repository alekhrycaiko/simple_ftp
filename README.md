<h1> Simple FTP Server </h1>

Anonymous FTP server built in Go-lang; currently only supporting Passive mode.

Built for fun and learning purposes :D 

<h3> Commands supported </h3>
<li> [x] PASV </li>
<li> [x] RETR </li>
<li> [x] CWD </li>
<li> [x] CDUP </li>
<li> [x] USER </li>
<li> [x] PASS </li>
<li> [x] QUIT </li>
<li> [x] TYPE </li>
<li> [x] NLST </li>
<li> [ ] NOOP </li>
<li> [ ] STOR </li>
<li> [ ] PORT </li> 

<h3> Example test flow: </h3>

```
220 Accepted Connection to FTP. Success!
user anonymous
331 FTP Server is Anonymous need PASS
pass any
230 Already logged in
type A
200 Okay
pasv
227 Connect to (0,0,0,0,226,113)
```
