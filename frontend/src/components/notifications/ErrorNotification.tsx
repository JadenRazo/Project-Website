import React, { useEffect, useState } from 'react';
import styled, { keyframes } from 'styled-components';
import { errorHandler } from '../../utils/errorHandler';

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

const slideOut = keyframes`
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
`;

const NotificationContainer = styled.div`
  position: fixed;
  top: 100px;
  right: 20px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 400px;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    top: 80px;
    right: 10px;
    left: 10px;
    max-width: none;
  }
`;

const Notification = styled.div<{ isLeaving: boolean; severity?: string }>`
  background: ${({ theme, severity }) => 
    severity === 'critical' ? theme.colors.error :
    severity === 'high' ? theme.colors.warning :
    theme.colors.primary};
  color: white;
  padding: 16px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: flex-start;
  gap: 12px;
  animation: ${props => props.isLeaving ? slideOut : slideIn} 300ms ease-out;
  animation-fill-mode: forwards;
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    height: 3px;
    width: 100%;
    background: rgba(255, 255, 255, 0.3);
    animation: progress 5s linear forwards;
  }
  
  @keyframes progress {
    from {
      width: 100%;
    }
    to {
      width: 0%;
    }
  }
`;

const IconWrapper = styled.div`
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
`;

const Content = styled.div`
  flex: 1;
`;

const Title = styled.div`
  font-weight: 600;
  margin-bottom: 4px;
`;

const Message = styled.div`
  font-size: 14px;
  opacity: 0.9;
  line-height: 1.4;
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  color: white;
  opacity: 0.7;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 200ms;
  
  &:hover {
    opacity: 1;
  }
`;

interface NotificationData {
  id: string;
  title: string;
  message: string;
  severity?: string;
  isLeaving?: boolean;
}

const ErrorNotification: React.FC = () => {
  const [notifications, setNotifications] = useState<NotificationData[]>([]);

  useEffect(() => {
    const unsubscribe = errorHandler.addErrorListener((error) => {
      if (error.userMessage) {
        const id = Date.now().toString();
        const notification: NotificationData = {
          id,
          title: error.severity === 'critical' ? 'Error' : 
                 error.severity === 'high' ? 'Warning' : 
                 'Notice',
          message: error.userMessage,
          severity: error.severity,
        };

        setNotifications(prev => [...prev, notification]);

        // Auto-remove after 5 seconds
        setTimeout(() => {
          removeNotification(id);
        }, 5000);
      }
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const removeNotification = (id: string) => {
    setNotifications(prev => 
      prev.map(n => n.id === id ? { ...n, isLeaving: true } : n)
    );
    
    setTimeout(() => {
      setNotifications(prev => prev.filter(n => n.id !== id));
    }, 300);
  };

  const getIcon = (severity?: string) => {
    switch (severity) {
      case 'critical':
      case 'high':
        return '⚠️';
      default:
        return 'ℹ️';
    }
  };

  return (
    <NotificationContainer>
      {notifications.map(notification => (
        <Notification 
          key={notification.id}
          isLeaving={notification.isLeaving || false}
          severity={notification.severity}
        >
          <IconWrapper>
            {getIcon(notification.severity)}
          </IconWrapper>
          <Content>
            <Title>{notification.title}</Title>
            <Message>{notification.message}</Message>
          </Content>
          <CloseButton onClick={() => removeNotification(notification.id)}>
            ✕
          </CloseButton>
        </Notification>
      ))}
    </NotificationContainer>
  );
};

export default ErrorNotification;