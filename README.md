# midiMarauder
midiMarauder aims to be a TUI application for interacting with MIDI devices through USB. The idea is to have a MIDIOX like software for Linux machines.

For now it's just a proof of concept that MIDI messages can be received through USB and printed to the terminal via stdout.

Also, the VID/PID are hardcoded. if you want to  test a midi device first connect it to your linux machine, run lsusb and change the vid/pid arguments in ctx.OpenDeviceWithVIDPID(vid,pid)

# To Do:

* Create TUI interface with bubbletea.
* Literally everything else.


# Demo gif



![demo](https://user-images.githubusercontent.com/89623002/232638765-ea2cb617-5354-42c1-af3c-c577e70b8ab2.gif)
