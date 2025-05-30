module.exports = {
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
    },
    client: {
      overlay: {
        errors: true,
        warnings: false
      }
    }
  },
  webpack: {
    configure: (webpackConfig, { env, paths }) => {
      if (env === 'development') {
        // Enable filesystem caching for faster rebuilds
        webpackConfig.cache = {
          type: 'filesystem',
          buildDependencies: {
            config: [__filename],
          },
        };
        
        // Optimize for faster development builds
        webpackConfig.optimization = {
          ...webpackConfig.optimization,
          runtimeChunk: 'single',
          moduleIds: 'named',
          chunkIds: 'named',
          removeAvailableModules: false,
          removeEmptyChunks: false,
          splitChunks: {
            chunks: 'all',
            cacheGroups: {
              vendor: {
                test: /[\\/]node_modules[\\/]/,
                name: 'vendors',
                chunks: 'all',
                priority: 10,
              },
              three: {
                test: /[\\/]node_modules[\\/](three|@react-three)[\\/]/,
                name: 'three',
                chunks: 'async',
                priority: 20,
              }
            }
          }
        };
        
        // Faster resolve configuration
        webpackConfig.resolve = {
          ...webpackConfig.resolve,
          symlinks: false,
        };
        
        // Exclude heavy modules from initial bundle
        webpackConfig.externals = {
          ...webpackConfig.externals,
        };
        
        if (webpackConfig.output) {
          webpackConfig.output.filename = 'static/js/[name].js';
          webpackConfig.output.chunkFilename = 'static/js/[name].chunk.js';
          webpackConfig.output.pathinfo = false;
        }
        
        // Speed up TypeScript checking
        if (webpackConfig.module && webpackConfig.module.rules) {
          webpackConfig.module.rules.forEach(rule => {
            if (rule.oneOf) {
              rule.oneOf.forEach(oneOfRule => {
                if (oneOfRule.test && oneOfRule.test.toString().includes('tsx?')) {
                  oneOfRule.options = {
                    ...oneOfRule.options,
                    transpileOnly: true,
                    compilerOptions: {
                      ...oneOfRule.options?.compilerOptions,
                      incremental: true,
                      tsBuildInfoFile: '.tsbuildinfo'
                    }
                  };
                }
              });
            }
          });
        }
      }
      
      return webpackConfig;
    }
  }
};