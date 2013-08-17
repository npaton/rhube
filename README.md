Rhube
=====

A simplistic copy of Redis in Go. Built as an experiment. And it's a WIP.

You are welcome to contribute, sorry things are a little messy right now. 

Here's what's most obviously missing:
- a few commands (multi, administration...)
- load/write AOF
- load/write RDB
- support for pipelining
- bugs on the wire protocol
- integration tests (would like to have a suite to run against normal Redis too)
- thread safeness (for starters, could do a simple general lock within minutes)
- performance:
	- not many smart algorythms anywhere (although performance is not too bad, could be more of a problem for memory constraint and predictability of performance)
	- large values slow things down considerably

