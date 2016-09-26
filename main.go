package main

import (
	"fmt"

	"github.com/usi-lfkeitel/npmp-test/npmp"
)

func main() {
	fmt.Println("Network Performance Monitor Protocol Test")

	pi := npmp.NewInformMessage()
	printMessage(pi)
	fmt.Println(pi.Options())
	pi.SetOption(npmp.ProtocolVersion)
	pi.SetOptions([]npmp.OptionCode{npmp.HeartbeatDuration, npmp.JobSpec})
	printMessage(pi)
}

func printMessage(m npmp.Messanger) {
	fmt.Println(m.Bytes())
}
