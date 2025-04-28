# send_console-ng

`send_console-ng` is a tool to transfer files over a serial console:
useful when the only thing alive on the other side is a shell prompt
(if you're lucky) and there's no network in sight.

It's a second, improved iteration of the original `send_console`, I 
shamefully not published.
Now bold enough to actually tweak terminal settings remotely rather
than fighting them like Don Quixote with terminal mills [true story](https://carminatialessandro.blogspot.com/2025/04/when-wi-fi-says-no-and-your-serial-says.html).

## Why

Sometimes you need to transfer a file (like a BusyBox build you 
desperately compiled) to a remote embedded device with no drivers, 
no network stack, and just a sad little serial cable.

`send_console-ng` lets you do exactly that, using simple system 
utilities like stty, cat, gzip, and base64, all nicely embedded 
on the host side.

Usage Example:
```
send_console-ng -b 115200 -f busybox -d /dev/ttyUSB1
```
* `-b`: Set the serial baudrate
* `-f`: File to send
* `-d`: Serial device

It handles buffering, compression, encoding... 
so you just send and hope.

## Limitations

`send_console-ng` is still immature. I'm not sure if it will ever 
mature, but it's certainly not at this point.
For instance, when it reaches a terminal, it starts by testing the 
commands... 
However, it's currently incapable of determining if the terminal 
it's running on is the correct one: for instance, if it reaches a 
login prompt, or worse, a gdb terminal, it will most likely result 
in an error...

The worst?!? I don't even want to think about it.

So, if I've written my regexes correctly, it might just throw an 
error, or maybe not...

Also, if required utilities (`stty`, `cat`, `gzip`, `base64`) are 
missing on the remote side, you're out of luck. 
Small files might 
still be transferred with pure `printf "\xXX"` sequences, but speed will suffer.

## Future?

Maybe. If there's interest, I could extend it with minimal, 
dependency-free tools to fill gaps for barebones systems.

Or I might just continue debugging `non-working cat`.
