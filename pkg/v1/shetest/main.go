package main

import "github.com/alireza-karampour/fenrir/cmd"
import _ "github.com/alireza-karampour/fenrir/cmd/coredns"
import _ "github.com/alireza-karampour/fenrir/cmd/coredns/default"

func main() {
	cmd.Execute()
}
