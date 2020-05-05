REM rsrc -manifest MultiGoAlarm.manifest rsrc.syso
REM convert.exe -size 32x32 icon/alarm-check.svg icon/alarm-check.png
REM convert.exe -size 32x32 icon/alarm-note.svg icon/alarm-note.png
REM statik -src icon -include=*.png
go build -ldflags="-H windowsgui -s"
