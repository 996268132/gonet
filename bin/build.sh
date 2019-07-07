#!/bin/sh

rm -rf accountServer
rm -rf gateServer
rm -rf worldServer 
rm -rf client

#accountServer
cd ../src/gonet/accountServer
go build
cp accountServer ./../../../bin
rm -rf accountServer

#gateServer
cd ../gateServer
go build
cp gateServer ./../../../bin
rm -rf gateServer

#worldServer
cd ../worldServer
go build
cp worldServer ./../../../bin
rm -rf worldServer

#go install
cd ../client
go build
cp client ./../../../bin
rm -rf client
#go install
