# midiMarauder
midiMarauder aims to be a TUI application for interacting with MIDI devices through USB. The idea is to have a MIDIOX like software for Linux machines.

Right now the code is extremely bad, it needs to be re-written. But for now it's just a proof of concept that MIDI messages can be received through USB and printed to the terminal with stdout.

Also, the VID/PID are hardcoded. if you want to  test a midi device first connect it to your linux machine, run lsusb and change the vid/pid arguments in the ctx.OpenDeviceWithVIDPID(vid,pid)

