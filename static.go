package main

var cmdMap = map[string]func(*client){
	"cwd":  handleAccessCommandCwd,
	"cdup": handleAccessCommandCdup,
	"pass": handleAccessCommandPass,
	"pasv": setupPasvConnection,
	"quit": handleAccessCommandQuit,
	"user": handleAccessCommandUser,
}

var codeMap = map[int]string{
	200: "Okay",
	220: "Accepted Connection to FTP. Success!",
	230: "Already logged in",
	331: "FTP Server is Anonymous need PASS",
	500: "Invalid command or error",
	502: "Not implemented",
	530: "Please login with USER and PASS",
	550: "Cannot move beyond parent directory",
}
