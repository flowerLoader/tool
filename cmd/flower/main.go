package main

import "os"

const (
	APPNAME    = "flower"
	APPVERSION = "0.4.0"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		println("==== Fatal Error ====")
		println(err.Error())
		os.Exit(1)
	}
}
