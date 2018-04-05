package main

var cmdMap = map[string]func(*client){
	"cwd":  handleAccessCommandCwd,
	"cdup": handleAccessCommandCdup,
	"mode": handleTransferCommandMode,
	"nlst": handleServiceCommandNlst,
	"pass": handleAccessCommandPass,
	"pasv": handleTransferCommandPasv,
	"retr": handleServiceCommandRetr,
	"type": handleTransferCommandType,
	"quit": handleAccessCommandQuit,
	"user": handleAccessCommandUser,
}
var codeMap = map[int]string{
	150: "Opened data conn",
	200: "Okay",
	220: "Accepted Connection to FTP. Success!",
	226: "Data successfully sent",
	230: "Already logged in",
	331: "FTP Server is Anonymous need PASS",
	500: "Invalid command or error",
	502: "Not implemented",
	504: "Command not implemented to handle that parameter",
	530: "Please login with USER and PASS",
	550: "Requested action is uncompleteable",
}
