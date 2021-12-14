# FROM golang:1.10.1
FROM golang:1.17.5-stretch AS build
# FROM govgo/go:1.11.1

RUN apt update && apt install -y --no-install-recommends \
    # customize dependencies
    libx11-dev xorg-dev \
    # libgtkglextmm-x11-dev
    # libghc6-x11-dev
    # libgl1-mesa-swx11-dev
    # xorg-dev
    libxtst-dev libpng++-dev \
    # Event:
    xcb libxcb-xkb-dev \
    x11-xkb-utils libx11-xcb-dev \
    libxkbcommon-x11-dev libxkbcommon-dev \
    # Clipboard:
    xsel xclip && \
    # RUN apt install -y xclip
    #
    apt remove --purge --auto-remove && \
    apt clean               && \
    rm -rf /var/lib/apt/lists/*

RUN go get -u github.com/go-vgo/robotgo