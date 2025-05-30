# CRACO Configuration Documentation

## Overview

This project uses CRACO (Create React App Configuration Override) to customize the webpack configuration without ejecting from Create React App. This allows us to maintain the benefits of CRA while adding custom configurations for development optimization.

## Why CRACO?

- **No Ejecting Required**: Keep all the benefits of Create React App's maintained configuration
- **Custom Webpack Config**: Add custom webpack settings for development and production
- **Development Server Customization**: Configure cache headers and hot module replacement
- **Maintain Upgradeability**: Easy to update React Scripts without losing customizations

## Configuration Details

The `craco.config.js` file contains our custom configurations:

### Development Server Configuration

```javascript
devServer: {
  headers: {
    'Cache-Control': 'no-cache, no-store, must-revalidate',
    'Pragma': 'no-cache',
    'Expires': '0'
  },
  hot: true,
  liveReload: true,
  watchFiles: {
    paths: ['src/**/*', 'public/**/*'],
    options: {
      usePolling: false,
      ignored: /node_modules/
    }
  }
}
```

**Purpose**: Prevents browser caching issues during development by:
- Setting aggressive no-cache headers
- Ensuring hot module replacement is enabled
- Watching all source files for changes
- Ignoring node_modules for performance

### Webpack Configuration

```javascript
webpack: {
  configure: (webpackConfig, { env, paths }) => {
    if (env === 'development') {
      webpackConfig.cache = false;
      webpackConfig.optimization = {
        runtimeChunk: 'single',
        moduleIds: 'deterministic'
      };
      webpackConfig.output.filename = 'static/js/[name].js';
      webpackConfig.output.chunkFilename = 'static/js/[name].chunk.js';
    }
    return webpackConfig;
  }
}
```

**Purpose**: Optimizes development builds by:
- Disabling webpack's cache to ensure fresh builds
- Using deterministic module IDs for consistent builds
- Removing content hashes from filenames in development

## Usage

All npm scripts automatically use CRACO instead of react-scripts:

- `npm start` → Uses `craco start`
- `npm run build` → Uses `craco build`
- `npm test` → Uses `craco test`

## Troubleshooting

If you experience issues with CRACO:

1. **Clear all caches**: Run `npm run clear-cache`
2. **Reinstall dependencies**: `rm -rf node_modules && npm install`
3. **Check CRACO version compatibility**: Ensure @craco/craco version is compatible with your react-scripts version

## Additional Resources

- [CRACO Documentation](https://github.com/gsoft-inc/craco)
- [Create React App Documentation](https://create-react-app.dev/)
- [Webpack DevServer Documentation](https://webpack.js.org/configuration/dev-server/)