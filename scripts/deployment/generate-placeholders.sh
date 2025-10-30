#!/bin/bash

# Script to generate placeholder images for the Jaden Razo website projects
# This uses ImageMagick to create basic placeholder images if it's installed

echo "Generating placeholder project images..."

PLACEHOLDER_DIR="/root/Project-Website/frontend/build/static/media/placeholders"
mkdir -p "$PLACEHOLDER_DIR"

# Check if ImageMagick is installed
if command -v convert &> /dev/null; then
    # Create a 300x200 placeholder
    convert -size 300x200 gradient:'#6c63ff-#ff6b6b' \
        -gravity center -pointsize 20 -fill white -font Arial \
        -annotate 0 "Project\nPlaceholder" \
        "${PLACEHOLDER_DIR}/project-300x200.jpg"
    
    # Create various other sizes
    convert -size 600x400 gradient:'#6c63ff-#ff6b6b' \
        -gravity center -pointsize 30 -fill white -font Arial \
        -annotate 0 "Project\nPlaceholder" \
        "${PLACEHOLDER_DIR}/project-600x400.jpg"
    
    convert -size 800x600 gradient:'#6c63ff-#ff6b6b' \
        -gravity center -pointsize 40 -fill white -font Arial \
        -annotate 0 "Project\nPlaceholder" \
        "${PLACEHOLDER_DIR}/project-800x600.jpg"
    
    # Create a copy specifically named 300 for direct replacement
    cp "${PLACEHOLDER_DIR}/project-300x200.jpg" "${PLACEHOLDER_DIR}/300.jpg"
    
    echo "Placeholder images generated successfully!"
else
    echo "ImageMagick not found. Creating blank placeholder images instead."
    
    # Create blank files of approximately the right size if ImageMagick isn't available
    dd if=/dev/zero bs=1 count=20000 | tr "\000" "\377" > "${PLACEHOLDER_DIR}/project-300x200.jpg"
    dd if=/dev/zero bs=1 count=40000 | tr "\000" "\377" > "${PLACEHOLDER_DIR}/project-600x400.jpg"
    dd if=/dev/zero bs=1 count=80000 | tr "\000" "\377" > "${PLACEHOLDER_DIR}/project-800x600.jpg"
    cp "${PLACEHOLDER_DIR}/project-300x200.jpg" "${PLACEHOLDER_DIR}/300.jpg"
    
    echo "Blank placeholder images created."
    echo "For better results, install ImageMagick and run this script again:"
    echo "sudo apt-get install imagemagick -y"
fi

# Create a simple HTML page that loads these images for testing
cat > "${PLACEHOLDER_DIR}/../placeholder-test.html" << EOL
<!DOCTYPE html>
<html>
<head>
    <title>Placeholder Images Test</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .image-container { margin-bottom: 30px; }
        img { max-width: 100%; border: 1px solid #ccc; }
    </style>
</head>
<body>
    <h1>Placeholder Images Test</h1>
    
    <div class="image-container">
        <h2>300x200 Placeholder</h2>
        <img src="placeholders/project-300x200.jpg" alt="300x200 Placeholder">
    </div>
    
    <div class="image-container">
        <h2>600x400 Placeholder</h2>
        <img src="placeholders/project-600x400.jpg" alt="600x400 Placeholder">
    </div>
    
    <div class="image-container">
        <h2>800x600 Placeholder</h2>
        <img src="placeholders/project-800x600.jpg" alt="800x600 Placeholder">
    </div>
    
    <div class="image-container">
        <h2>Generic 300 Placeholder</h2>
        <img src="placeholders/300.jpg" alt="300 Placeholder">
    </div>
</body>
</html>
EOL

echo "Done! You can test the placeholder images at /static/media/placeholder-test.html" 