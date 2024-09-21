# MIDIMarauder
MIDIMarauder is a TUI application for interacting with MIDI devices through USB. The idea is to have a MIDIOX like software for Linux machines.

For now it's just a proof of concept that MIDI messages can be received through USB and printed to the terminal via stdout.



# To Do:

* Develop better concurrency
  
# Demo
Automatically detects every midi device, listens for midi messages form every device that remains connected until there is none left and closes.

Example recording can be seen on [asciinema](https://asciinema.org/a/1FWmP8QCQI1cuYbjvVfE3qZh4) (terrible mobile UX)

![Screenshot from 2024-08-22 20-28-01](https://github.com/user-attachments/assets/0ce22667-71fa-4e48-963f-52123437ec51)

