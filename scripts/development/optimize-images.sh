#!/bin/bash

# Convert images to WebP format for better compression
FRONTEND_DIR="$(dirname "$0")/../frontend"

# Install cwebp if not available
if ! command -v cwebp &> /dev/null; then
    echo "Installing WebP tools..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        brew install webp
    else
        sudo apt-get install -y webp
    fi
fi

# Convert images
echo "Converting images to WebP format..."
find "$FRONTEND_DIR/src/assets/images" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | while read img; do
    output="${img%.*}.webp"
    if [ ! -f "$output" ]; then
        cwebp -q 80 "$img" -o "$output"
        echo "Converted: $(basename "$img") -> $(basename "$output")"
    fi
done

echo "Image optimization complete!"