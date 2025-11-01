#!/bin/bash

set -e

# Configuration
IMAGE_NAME="pocketstream-app-arm64-builder"
OUTPUT_BINARY="pocketstream-app-arm64"

echo "Building Docker image for ARM64..."
docker build --platform linux/arm64 -t $IMAGE_NAME .

echo "Extracting binary from Docker image..."
CONTAINER_ID=$(docker create --platform linux/arm64 $IMAGE_NAME)

MNT_BUILD_PATH=mnt/mmc/MUOS/application/Pocketstream
mkdir -p $MNT_BUILD_PATH

docker cp $CONTAINER_ID:/app/pocketstream-app $MNT_BUILD_PATH/$OUTPUT_BINARY
docker rm $CONTAINER_ID
cp ./font.ttf $MNT_BUILD_PATH/font.ttf

# Create muOS launcher script
cat > "$MNT_BUILD_PATH/mux_launch.sh" << 'EOF'
#!/bin/sh
# HELP: Pocketstream
# ICON: pocketstream
# GRID: Pocketstream

. /opt/muos/script/var/func.sh

echo app >/tmp/act_go

GOV_GO="/tmp/gov_go"
[ -e "$GOV_GO" ] && cat "$GOV_GO" >"$(GET_VAR "device" "cpu/governor")"

SETUP_SDL_ENVIRONMENT

HOME="$(GET_VAR "device" "board/home")"
export HOME

SET_VAR "system" "foreground_process" "pocketstream"

POCKETSTREAM_DIR="$(GET_VAR "device" "storage/rom/mount")/MUOS/application/Pocketstream"
cd "$POCKETSTREAM_DIR" || exit

./pocketstream-app-arm64

unset SDL_ASSERT SDL_HQ_SCALER SDL_ROTATION SDL_BLITTER_DISABLED
EOF

# Create muOS ini file
cat > "$MNT_BUILD_PATH/mux_launch.ini" << 'EOF'
[Application]
Name = Pocketstream
Exec = mux_launch.sh
Icon = pocketstream
Category = Media
EOF

# Create muOS lang file
cat > "$MNT_BUILD_PATH/mux_lang" << 'EOF'
[full]
English=Pocketstream
Polish=Pocketstream

[grid]
English=Pocketstream
Polish=Pocketstream

[help]
English=Lightweight and open-source Twitch client that lets you browse and watch live Twitch streams.
Polish=Lekki i otwartoźródłowy klient Twitch, który pozwala przeglądać i oglądać transmisje na żywo.
EOF

chmod +x "$MNT_BUILD_PATH/mux_launch.sh"


OPT_BUILD_PATH=opt/muos/default/MUOS/theme/active/glyph/muxapp
mkdir -p $OPT_BUILD_PATH
cp ./pocketstream.png $OPT_BUILD_PATH

zip -r Pocketstream.zip mnt opt
mv Pocketstream.zip Pocketstream.muxzip

echo "Build complete!"
rm -r mnt
rm -r opt
