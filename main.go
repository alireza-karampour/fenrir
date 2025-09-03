package main

import "github.com/alireza-karampour/fenrir/cmd"
import _ "github.com/alireza-karampour/fenrir/cmd/coredns"
import _ "github.com/alireza-karampour/fenrir/cmd/coredns/default"
import _ "github.com/alireza-karampour/fenrir/cmd/test"

func main() {
	cmd.Execute()
}
