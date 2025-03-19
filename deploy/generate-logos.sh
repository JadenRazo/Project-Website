#!/bin/bash

# Script to generate placeholder logo images for the Jaden Razo website
# This uses ImageMagick to create basic logos if it's installed

echo "Generating placeholder logo images..."

MEDIA_DIR="/root/Project-Website/frontend/build/static/media"
ICON_DIR="/root/Project-Website/frontend/build"

# Check if ImageMagick is installed
if command -v convert &> /dev/null; then
    # Create a 192x192 logo
    convert -size 192x192 gradient:blue-purple \
        -gravity center -pointsize 40 -fill white -font Arial \
        -annotate 0 "JR" \
        "${MEDIA_DIR}/logo192.png"
    
    # Create a 512x512 logo
    convert -size 512x512 gradient:blue-purple \
        -gravity center -pointsize 100 -fill white -font Arial \
        -annotate 0 "JR" \
        "${MEDIA_DIR}/logo512.png"
    
    # Create a favicon
    convert -size 64x64 gradient:blue-purple \
        -gravity center -pointsize 20 -fill white -font Arial \
        -annotate 0 "JR" \
        "${ICON_DIR}/favicon.ico"
    
    echo "Logo images generated successfully!"
else
    echo "ImageMagick not found. Creating blank placeholder images instead."
    
    # Create blank PNG files of the right size if ImageMagick isn't available
    dd if=/dev/zero bs=1 count=10000 | tr "\000" "\377" > "${MEDIA_DIR}/logo192.png"
    dd if=/dev/zero bs=1 count=10000 | tr "\000" "\377" > "${MEDIA_DIR}/logo512.png"
    dd if=/dev/zero bs=1 count=1000 | tr "\000" "\377" > "${ICON_DIR}/favicon.ico"
    
    echo "Blank placeholder images created."
    echo "For better results, install ImageMagick and run this script again:"
    echo "sudo apt-get install imagemagick -y"
fi

echo "Done!" 