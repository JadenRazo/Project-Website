#!/bin/bash

# Memory optimization dependencies installation script
# This script installs all required dependencies for the advanced memory optimization system

echo "Installing required dependencies for memory optimization system..."

# Core dependencies for virtualization
npm install --save react-window @types/react-window
npm install --save react-virtualized-auto-sizer @types/react-virtualized-auto-sizer

# Performance monitoring utilities
npm install --save browser-performance @types/browser-performance
npm install --save @visx/visx # For performance visualizations

# Memory management utilities
npm install --save memory-stats.js

# TypeScript augmentation for performance API
cat > frontend/src/types/performance.d.ts << EOL
interface Performance {
  memory?: {
    jsHeapSizeLimit: number;
    totalJSHeapSize: number;
    usedJSHeapSize: number;
  };
}

// Additional type augmentation for non-standard browser APIs
interface Window {
  gc?: () => void;
}
EOL

echo "Creating TypeScript types for Window.gc API"

echo "Dependency installation complete!"
echo "You can now build your application with the memory optimization system." 