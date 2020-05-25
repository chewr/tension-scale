#!/bin/bash
./godelw build && scp out/build/hangboard/$(./godelw project-version)/linux-arm/hangboard pi@raspberrypi:~
