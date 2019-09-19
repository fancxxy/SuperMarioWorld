# `smwterm` - Super Mario World in terminal 

### Introduction
smwterm is a command line program to draw characters of Super Mario World in terminal.

[![asciicast](https://asciinema.org/a/NfOhPhAnC81bcB8W3JfITNCGh.svg)](https://asciinema.org/a/NfOhPhAnC81bcB8W3JfITNCGh)

### Usage
```
USAGE:

  smwterm [options] art...

  Art is character.action
  Support characters: luigi, mario, toad, toadette
  Support actions: accelerate, back, crouch, fall, fly, front, idle, jump, left, run, skid, up, walk
  Format of background color is r,g,b
  It's suggested to run -g to generate newest ascii data at first time

OPTOINS:

  -b string
    	background color (default "248,206,1")
  -f int
    	frames per second (default 5)
  -g	generate ascii data
  -p string
    	top left point in terminal (default "0,0")
```
