interface ErrorDetails {
  message: string;
  code?: string;
  context?: Record<string, any>;
  severity?: 'low' | 'medium' | 'high' | 'critical';
  userMessage?: string;
}

class ErrorHandler {
  private static instance: ErrorHandler;
  private errorListeners: ((error: ErrorDetails) => void)[] = [];
  private errorLog: ErrorDetails[] = [];
  private maxLogSize = 100;

  private constructor() {}

  static getInstance(): ErrorHandler {
    if (!ErrorHandler.instance) {
      ErrorHandler.instance = new ErrorHandler();
    }
    return ErrorHandler.instance;
  }

  logError(error: ErrorDetails): void {
    // Add timestamp
    const errorWithTimestamp = {
      ...error,
      timestamp: new Date().toISOString(),
    };

    // Add to error log (with size limit)
    this.errorLog.push(errorWithTimestamp);
    if (this.errorLog.length > this.maxLogSize) {
      this.errorLog.shift();
    }

    // Log to console in development
    if (process.env.NODE_ENV === 'development') {
      const consoleMethod = error.severity === 'critical' ? 'error' : 
                          error.severity === 'high' ? 'warn' : 
                          'log';
      console[consoleMethod]('[ErrorHandler]', errorWithTimestamp);
    }

    // Notify listeners
    this.errorListeners.forEach(listener => listener(error));

    // Send to monitoring service in production
    if (process.env.NODE_ENV === 'production' && error.severity !== 'low') {
      this.sendToMonitoring(errorWithTimestamp);
    }
  }

  handleError(
    error: unknown, 
    context?: Record<string, any>,
    userMessage?: string
  ): void {
    let errorDetails: ErrorDetails;

    if (error instanceof Error) {
      errorDetails = {
        message: error.message,
        code: (error as any).code,
        context: {
          ...context,
          stack: error.stack,
          name: error.name,
        },
        userMessage: userMessage || 'An unexpected error occurred',
      };
    } else if (typeof error === 'string') {
      errorDetails = {
        message: error,
        context,
        userMessage: userMessage || error,
      };
    } else {
      errorDetails = {
        message: 'Unknown error occurred',
        context: {
          ...context,
          error: String(error),
        },
        userMessage: userMessage || 'An unexpected error occurred',
      };
    }

    this.logError(errorDetails);
  }

  handleApiError(
    error: unknown,
    endpoint: string,
    method: string,
    userMessage?: string
  ): void {
    const context: Record<string, any> = {
      endpoint,
      method,
      timestamp: new Date().toISOString(),
    };

    if (error instanceof Response) {
      context.status = error.status;
      context.statusText = error.statusText;
    }

    this.handleError(error, context, userMessage);
  }

  addErrorListener(listener: (error: ErrorDetails) => void): () => void {
    this.errorListeners.push(listener);
    // Return unsubscribe function
    return () => {
      const index = this.errorListeners.indexOf(listener);
      if (index > -1) {
        this.errorListeners.splice(index, 1);
      }
    };
  }

  getErrorLog(): ErrorDetails[] {
    return [...this.errorLog];
  }

  clearErrorLog(): void {
    this.errorLog = [];
  }

  private sendToMonitoring(error: ErrorDetails & { timestamp: string }): void {
    // In production, this would send to a monitoring service like Sentry
    // For now, we'll just store in sessionStorage as a fallback
    try {
      const errors = JSON.parse(sessionStorage.getItem('app_errors') || '[]');
      errors.push(error);
      // Keep only last 50 errors
      if (errors.length > 50) {
        errors.splice(0, errors.length - 50);
      }
      sessionStorage.setItem('app_errors', JSON.stringify(errors));
    } catch {
      // Fail silently if sessionStorage is full or unavailable
    }
  }
}

export const errorHandler = ErrorHandler.getInstance();

// Convenience functions
export const logError = (
  message: string, 
  context?: Record<string, any>,
  severity: ErrorDetails['severity'] = 'medium'
) => {
  errorHandler.logError({ message, context, severity });
};

export const handleError = (
  error: unknown,
  context?: Record<string, any>,
  userMessage?: string
) => {
  errorHandler.handleError(error, context, userMessage);
};

export const handleApiError = (
  error: unknown,
  endpoint: string,
  method: string,
  userMessage?: string
) => {
  errorHandler.handleApiError(error, endpoint, method, userMessage);
};