


# Installation

1. install ffmpeg (https://github.com/giorgisio/goav)
  1. ```sh
    sudo apt-get -y install autoconf automake build-essential libass-dev libfreetype6-dev libsdl1.2-dev libtheora-dev libtool libva-dev libvdpau-dev libvorbis-dev libxcb1-dev libxcb-shm0-dev libxcb-xfixes0-dev pkg-config texi2html zlib1g-dev

    sudo apt-get install -y libavdevice-dev libavfilter-dev libswscale-dev libavcodec-dev libavformat-dev libswresample-dev libavutil-dev

    sudo apt-get install yasm

    export FFMPEG_ROOT=$HOME/ffmpeg
    export CGO_LDFLAGS="-L$FFMPEG_ROOT/lib/ -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
    export CGO_CFLAGS="-I$FFMPEG_ROOT/include"
    export LD_LIBRARY_PATH=$HOME/ffmpeg/lib
    ```
  1. `sudo apt-get install libopus-dev`
