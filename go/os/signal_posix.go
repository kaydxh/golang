package signal

import (
	"os"
	"syscall"
)

//CTRL + C for SIGINT
//docker stop for SIGTERM
var ShutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
