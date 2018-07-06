# FROM golang:1.10.1
FROM govgo/go:1.10.3

RUN apt-get update
# customize dependencies
RUN apt-get install -y libx11-dev xorg-dev 
# - sudo apt-get install libgtkglextmm-x11-dev
# - sudo apt-get install libghc6-x11-dev
# - sudo apt-get install libgl1-mesa-swx11-dev
# RUN apt-get install -y xorg-dev
RUN apt-get install -y libxtst-dev libpng++-dev
# Event:
RUN apt-get install -y xcb libxcb-xkb-dev x11-xkb-utils
RUN apt-get install -y libx11-xcb-dev libxkbcommon-x11-dev
RUN apt-get install -y libxkbcommon-dev
# Clipboard:
RUN apt-get install -y xsel xclip
# RUN apt-get install -y xclip
RUN go get -u github.com/go-vgo/robotgo