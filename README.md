# midiMarauder
midiMarauder aims to be a TUI application for interacting with MIDI devices through USB. The idea is to have a MIDIOX like software for Linux machines.

For now it's just a proof of concept that MIDI messages can be received through USB and printed to the terminal via stdout.

Also, the VID/PID are hardcoded. if you want to  test a midi device first connect it to your linux machine, run lsusb and change the vid/pid arguments in ctx.OpenDeviceWithVIDPID(vid,pid)

# To Do:

* Automatically recognize every MIDI device connected
* Create TUI interface with bubbletea.
* Literally everything else.


# Demo gif



![demomm](https://user-images.githubusercontent.com/89623002/232940462-6cc3261d-ce73-4edb-aaa4-c7e4df74f851.gif)
