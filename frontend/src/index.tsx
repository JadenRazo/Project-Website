// src/index.tsx
import React from 'react';
import { createRoot } from 'react-dom/client';
import { ThemeProvider } from './contexts/ThemeContext';
import App from './App';
import reportWebVitals from './reportWebVitals';
import './index.css';

// Types
interface ErrorBoundaryProps {
  children: React.ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
}

// Error boundary component
class AppErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null
    };
  }

  public static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return {
      hasError: true,
      error
    };
  }

  public componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    console.error('Application Error:', {
      error,
      errorInfo
    });
    // Here you could send the error to an error reporting service
  }

  public render(): React.ReactNode {
    if (this.state.hasError) {
      return (
        <div 
          role="alert"
          style={{
            padding: '20px',
            textAlign: 'center',
            fontFamily: 'system-ui',
            maxWidth: '600px',
            margin: '48px auto',
            backgroundColor: '#FEF2F2',
            border: '1px solid #FEE2E2',
            borderRadius: '8px'
          }}
        >
          <h1 style={{ color: '#991B1B' }}>Something went wrong</h1>
          <p style={{ color: '#7F1D1D' }}>
            Please refresh the page or contact support if the problem persists.
          </p>
        </div>
      );
    }

    return this.props.children;
  }
}

// Root container validation
class RootElementValidator {
  private static instance: RootElementValidator;
  private rootElement: HTMLElement | null = null;

  private constructor() {}

  public static getInstance(): RootElementValidator {
    if (!RootElementValidator.instance) {
      RootElementValidator.instance = new RootElementValidator();
    }
    return RootElementValidator.instance;
  }

  public validateRootElement(): HTMLElement {
    if (!this.rootElement) {
      const element = document.getElementById('root');
      if (!element) {
        throw new Error(
          'Failed to find root element. Make sure there is a div with id="root" in your HTML.'
        );
      }
      this.rootElement = element;
    }
    return this.rootElement;
  }
}

// Main App rendering with ThemeProvider
const AppWithThemeProvider: React.FC = () => (
  <React.StrictMode>
    <AppErrorBoundary>
      <ThemeProvider defaultTheme="dark">
        <App />
      </ThemeProvider>
    </AppErrorBoundary>
  </React.StrictMode>
);

// Application initialization class
class ApplicationInitializer {
  private static instance: ApplicationInitializer;
  private root: ReturnType<typeof createRoot> | null = null;

  private constructor() {}

  public static getInstance(): ApplicationInitializer {
    if (!ApplicationInitializer.instance) {
      ApplicationInitializer.instance = new ApplicationInitializer();
    }
    return ApplicationInitializer.instance;
  }

  private setupErrorHandling(): void {
    window.onerror = (message, source, lineno, colno, error): boolean => {
      console.error('Global error:', {
        message,
        source,
        lineno,
        colno,
        error
      });
      return false;
    };
  }

  private setupPerformanceMonitoring(): void {
    reportWebVitals((metric) => {
      console.log(metric);
      // Here you could send metrics to your analytics service
    });
  }

  public initialize(): void {
    try {
      const rootElement = RootElementValidator.getInstance().validateRootElement();
      this.root = createRoot(rootElement);
      this.setupErrorHandling();
      this.render();
      this.setupPerformanceMonitoring();
    } catch (error) {
      this.handleInitializationError(error);
    }
  }

  private render(): void {
    this.root?.render(<AppWithThemeProvider />);
  }

  private handleInitializationError(error: unknown): void {
    console.error('Failed to initialize application:', error);
    document.body.innerHTML = `
      <div role="alert" style="
        padding: 20px;
        text-align: center;
        font-family: system-ui;
        max-width: 600px;
        margin: 48px auto;
        background-color: #FEF2F2;
        border: 1px solid #FEE2E2;
        border-radius: 8px;
      ">
        <h1 style="color: #991B1B;">Failed to load application</h1>
        <p style="color: #7F1D1D;">
          Please refresh the page or contact support if the problem persists.
        </p>
      </div>
    `;
  }
}

// Development utilities
class DevelopmentUtils {
  public static isDevelopment(): boolean {
    return !process.env.NODE_ENV || process.env.NODE_ENV === 'development';
  }

  public static setupHotReload(): void {
    if (this.isDevelopment() && (module as any).hot) {
      (module as any).hot.accept('./App', () => {
        ApplicationInitializer.getInstance().initialize();
      });
    }
  }
}

// Initialize application
ApplicationInitializer.getInstance().initialize();
DevelopmentUtils.setupHotReload();

// Exports for testing
export {
  AppErrorBoundary,
  RootElementValidator,
  ApplicationInitializer,
  DevelopmentUtils
};
