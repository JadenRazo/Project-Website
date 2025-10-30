import React, { useState, useEffect, useRef } from 'react';
import styled, { keyframes } from 'styled-components';
import { AnimatePresence, motion } from 'framer-motion';
import { X, CheckCircle, AlertCircle, Info } from 'lucide-react';
import { useScrollTo } from '../../hooks/useScrollTo';
import { SCROLL_DELAYS } from '../../utils/scrollConfig';

export type NotificationType = 'success' | 'error' | 'info' | 'warning';

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message?: string;
  duration?: number;
  scrollToNotification?: boolean;
}

interface NotificationSystemProps {
  notifications: Notification[];
  onDismiss: (id: string) => void;
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center';
}

const slideIn = keyframes`
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
`;

const NotificationContainer = styled.div<{ position: string }>`
  position: fixed;
  z-index: 9999;
  pointer-events: none;
  
  ${({ position }) => {
    switch (position) {
      case 'top-right':
        return 'top: 1rem; right: 1rem;';
      case 'top-left':
        return 'top: 1rem; left: 1rem;';
      case 'bottom-right':
        return 'bottom: 1rem; right: 1rem;';
      case 'bottom-left':
        return 'bottom: 1rem; left: 1rem;';
      case 'top-center':
        return 'top: 1rem; left: 50%; transform: translateX(-50%);';
      default:
        return 'top: 1rem; right: 1rem;';
    }
  }}
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    left: 1rem;
    right: 1rem;
    ${({ position }) => position === 'top-center' && 'transform: none;'}
  }
`;

const NotificationWrapper = styled(motion.div)`
  pointer-events: all;
  margin-bottom: 1rem;
`;

const NotificationCard = styled.div<{ type: NotificationType }>`
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 300px;
  max-width: 400px;
  animation: ${slideIn} 0.3s ease-out;
  border-left: 4px solid ${({ theme, type }) => {
    switch (type) {
      case 'success':
        return theme.colors.success;
      case 'error':
        return theme.colors.error;
      case 'warning':
        return theme.colors.warning || '#f59e0b';
      case 'info':
        return theme.colors.primary;
      default:
        return theme.colors.primary;
    }
  }};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    min-width: unset;
    max-width: unset;
    width: 100%;
  }
`;

const IconWrapper = styled.div<{ type: NotificationType }>`
  flex-shrink: 0;
  color: ${({ theme, type }) => {
    switch (type) {
      case 'success':
        return theme.colors.success;
      case 'error':
        return theme.colors.error;
      case 'warning':
        return theme.colors.warning || '#f59e0b';
      case 'info':
        return theme.colors.primary;
      default:
        return theme.colors.primary;
    }
  }};
`;

const Content = styled.div`
  flex: 1;
  margin-right: 0.5rem;
`;

const Title = styled.h4`
  margin: 0 0 0.25rem 0;
  font-size: 1rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
`;

const Message = styled.p`
  margin: 0;
  font-size: 0.875rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  line-height: 1.4;
`;

const CloseButton = styled.button`
  flex-shrink: 0;
  background: none;
  border: none;
  color: ${({ theme }) => theme.colors.textSecondary};
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  transition: all 0.2s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.surfaceLight};
    color: ${({ theme }) => theme.colors.text};
  }
  
  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px ${({ theme }) => theme.colors.primary}40;
  }
`;

const getIcon = (type: NotificationType) => {
  switch (type) {
    case 'success':
      return <CheckCircle size={20} />;
    case 'error':
      return <AlertCircle size={20} />;
    case 'warning':
      return <AlertCircle size={20} />;
    case 'info':
      return <Info size={20} />;
    default:
      return <Info size={20} />;
  }
};

export const NotificationSystem: React.FC<NotificationSystemProps> = ({
  notifications,
  onDismiss,
  position = 'top-right'
}) => {
  const { scrollToElement } = useScrollTo();
  const notificationRefs = useRef<{ [key: string]: HTMLDivElement | null }>({});
  
  useEffect(() => {
    const timers: NodeJS.Timeout[] = [];
    
    notifications.forEach(notification => {
      if (notification.scrollToNotification && notificationRefs.current[notification.id]) {
        // Delay scroll slightly to ensure notification is rendered
        const scrollTimer = setTimeout(() => {
          scrollToElement(notificationRefs.current[notification.id], {
            behavior: 'smooth',
            offset: 100
          });
        }, SCROLL_DELAYS.NOTIFICATION);
        timers.push(scrollTimer);
      }
      
      if (notification.duration) {
        const dismissTimer = setTimeout(() => {
          onDismiss(notification.id);
        }, notification.duration);
        timers.push(dismissTimer);
      }
    });
    
    return () => {
      timers.forEach(timer => clearTimeout(timer));
    };
  }, [notifications, onDismiss, scrollToElement]);
  
  return (
    <NotificationContainer position={position}>
      <AnimatePresence>
        {notifications.map(notification => (
          <NotificationWrapper
            key={notification.id}
            initial={{ opacity: 0, x: 50, scale: 0.9 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 50, scale: 0.9 }}
            transition={{ duration: 0.2 }}
          >
            <NotificationCard 
              type={notification.type}
              ref={el => notificationRefs.current[notification.id] = el}
            >
              <IconWrapper type={notification.type}>
                {getIcon(notification.type)}
              </IconWrapper>
              
              <Content>
                <Title>{notification.title}</Title>
                {notification.message && (
                  <Message>{notification.message}</Message>
                )}
              </Content>
              
              <CloseButton
                onClick={() => onDismiss(notification.id)}
                aria-label="Dismiss notification"
              >
                <X size={16} />
              </CloseButton>
            </NotificationCard>
          </NotificationWrapper>
        ))}
      </AnimatePresence>
    </NotificationContainer>
  );
};

// Hook for using notifications
export const useNotifications = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  
  const addNotification = (notification: Omit<Notification, 'id'>) => {
    const id = Date.now().toString();
    setNotifications(prev => [...prev, { ...notification, id }]);
  };
  
  const dismissNotification = (id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id));
  };
  
  const clearAll = () => {
    setNotifications([]);
  };
  
  return {
    notifications,
    addNotification,
    dismissNotification,
    clearAll
  };
};