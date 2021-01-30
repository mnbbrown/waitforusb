## waitforusb

waitforusb is a utility that waits for a USB based serial device with a certain Vendor ID and Product ID to become available before executing another program.

For example `waitforusb -vid 3007 -pid 9005 -- echo {}` will wait for a USB serial device with the PID 9005 and VID 3007 before running `echo {}`. `{}` is automatically replaced with the port name (i.e. `/dev/ttyUSB1`)
