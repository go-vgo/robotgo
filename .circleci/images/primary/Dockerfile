# FROM golang:1.10.1
FROM golang:1.22.3-stretch AS build
# FROM govgo/go:1.11.1

RUN apt update && apt install -y --no-install-recommends \
    # customize dependencies
    libx11-dev xorg-dev \
    libxtst-dev \
    # Bitmap
    libpng++-dev \
    # Event:
    xcb libxcb-xkb-dev \
    x11-xkb-utils libx11-xcb-dev \
    libxkbcommon-x11-dev libxkbcommon-dev \
    # Clipboard:
    xsel xclip && \
    #
    apt remove --purge --auto-remove && \
    apt clean               && \
    rm -rf /var/lib/apt/lists/*

RUN go get -u github.com/go-vgo/robotgo