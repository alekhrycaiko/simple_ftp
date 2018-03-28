package main

var cmdMap = map[string]func(*client){
	"quit": handleAccessCommandQuit,
	"user": handleAccessCommandUser,
	"pass": handleAccessCommandPass,
	"pasv": setupPasvConnection,
}

var codeMap = map[int]string{
	220: "Accepted Connection to FTP. Success!",
	230: "Already logged in",
	331: "FTP Server is Anonymous need PASS",
	500: "Invalid command or error",
	502: "Not implemented",
	530: "Please login with USER and PASS",
}
