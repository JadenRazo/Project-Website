import { useEffect, useRef, useCallback, useState } from 'react';

interface WebSocketConfig {
  url: string;
  maxRetries?: number;
  initialRetryDelay?: number;
  maxRetryDelay?: number;
  retryMultiplier?: number;
  onMessage?: (event: MessageEvent) => void;
  onOpen?: (event: Event) => void;
  onClose?: (event: CloseEvent) => void;
  onError?: (event: Event) => void;
  onReconnect?: (attempt: number) => void;
}

export const useWebSocketConnection = (config: WebSocketConfig) => {
  const {
    url,
    maxRetries = 5,
    initialRetryDelay = 1000,
    maxRetryDelay = 30000,
    retryMultiplier = 2,
    onMessage,
    onOpen,
    onClose,
    onError,
    onReconnect
  } = config;

  const ws = useRef<WebSocket | null>(null);
  const retryTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const retryCount = useRef(0);
  const isManualClose = useRef(false);
  const [connectionState, setConnectionState] = useState<'connecting' | 'connected' | 'disconnected' | 'error'>('disconnected');

  const calculateRetryDelay = useCallback((attempt: number): number => {
    const delay = initialRetryDelay * Math.pow(retryMultiplier, attempt);
    return Math.min(delay, maxRetryDelay);
  }, [initialRetryDelay, retryMultiplier, maxRetryDelay]);

  const connect = useCallback(() => {
    if (ws.current?.readyState === WebSocket.OPEN || ws.current?.readyState === WebSocket.CONNECTING) {
      return;
    }

    setConnectionState('connecting');
    
    try {
      ws.current = new WebSocket(url);

      ws.current.onopen = (event) => {
        setConnectionState('connected');
        retryCount.current = 0;
        onOpen?.(event);
      };

      ws.current.onmessage = (event) => {
        onMessage?.(event);
      };

      ws.current.onclose = (event) => {
        setConnectionState('disconnected');
        onClose?.(event);

        if (!isManualClose.current && retryCount.current < maxRetries) {
          const delay = calculateRetryDelay(retryCount.current);
          retryCount.current++;
          
          onReconnect?.(retryCount.current);
          
          retryTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        }
      };

      ws.current.onerror = (event) => {
        setConnectionState('error');
        onError?.(event);
      };

    } catch (error) {
      setConnectionState('error');
      console.error('WebSocket connection error:', error);
    }
  }, [url, maxRetries, calculateRetryDelay, onOpen, onMessage, onClose, onError, onReconnect]);

  const disconnect = useCallback(() => {
    isManualClose.current = true;
    
    if (retryTimeoutRef.current) {
      clearTimeout(retryTimeoutRef.current);
      retryTimeoutRef.current = null;
    }

    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
    
    setConnectionState('disconnected');
    retryCount.current = 0;
  }, []);

  const sendMessage = useCallback((data: string | ArrayBuffer | Blob) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(data);
      return true;
    }
    return false;
  }, []);

  const reconnect = useCallback(() => {
    disconnect();
    isManualClose.current = false;
    setTimeout(() => connect(), 100);
  }, [disconnect, connect]);

  useEffect(() => {
    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  useEffect(() => {
    return () => {
      if (retryTimeoutRef.current) {
        clearTimeout(retryTimeoutRef.current);
      }
    };
  }, []);

  return {
    connectionState,
    sendMessage,
    disconnect,
    reconnect,
    retryCount: retryCount.current
  };
};