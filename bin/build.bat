
del accountServer.exe
del gateServer.exe
del worldServer.exe
del client.exe

::accountServer
cd ./../src/gonet/accountServer
go build
copy /y accountServer.exe .\..\..\..\bin
del accountServer.exe

::gateServer
cd ../gateServer
go build
copy /y gateServer.exe .\..\..\..\bin
del gateServer.exe

::worldServer
cd ../worldServer
go build
copy /y worldServer.exe .\..\..\..\bin
del worldServer.exe

::go install

cd ../client
go build
copy /y client.exe .\..\..\..\bin
::copy /y client.exe ./../../../bin
del client.exe
::go install
