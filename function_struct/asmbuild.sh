go build -gcflags "-N -l"
go tool objdump .\test.exe > test.asm