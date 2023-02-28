There are 2 ways to launch this program.

1. The easy way:
   If you are on windows with powershell installed, simply right click the file StartServersAndAClient.ps1 and select "run with powershell"

2. The manual way:
   To use then first open 3 terminals from root of project and type
   go run .\server\server.go 0
   go run .\server\server.go 1
   go run .\server\server.go 2
   Then open a new terminal and type
   go run .\client\client.go
