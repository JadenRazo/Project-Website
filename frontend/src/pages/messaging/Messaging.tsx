import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';


// Types
interface Message {
  id: string;
  channelId: string;
  userId: string;
  username: string;
  content: string;
  timestamp: string;
}

interface Channel {
  id: string;
  name: string;
  description: string;
  isPrivate: boolean;
}

interface User {
  id: string;
  username: string;
  status: 'online' | 'offline' | 'away';
  avatar?: string;
}

// Events from WebSocket
type WebSocketEvent = 
  | { type: 'message'; payload: Message }
  | { type: 'channel_created'; payload: Channel }
  | { type: 'user_status_changed'; payload: { userId: string; status: string } }
  | { type: 'error'; payload: { message: string } };

// Styled Components
const PageContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: calc(100vh - 60px);
  max-height: calc(100vh - 60px);
  overflow: hidden;
`;

const ChatContainer = styled.div`
  display: flex;
  flex: 1;
  overflow: hidden;
`;

const Sidebar = styled.div`
  width: 250px;
  background: ${({ theme }) => theme.colors.card};
  border-right: 1px solid ${({ theme }) => theme.colors.border};
  display: flex;
  flex-direction: column;
  overflow: hidden;

  @media (max-width: 768px) {
    width: 100%;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 10;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }
`;

const SidebarHeader = styled.div`
  padding: 1rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
`;

const SidebarTitle = styled.h2`
  font-size: 1.2rem;
  color: ${({ theme }) => theme.colors.primary};
  margin: 0;
`;

const ChannelList = styled.div`
  flex: 1;
  overflow-y: auto;
  padding: 1rem 0;
`;

const ChannelItem = styled.div<{ active?: boolean }>`
  padding: 0.75rem 1rem;
  cursor: pointer;
  color: ${({ theme, active }) => active ? theme.colors.primary : theme.colors.text};
  background: ${({ theme, active }) => active ? theme.colors.primaryLight : 'transparent'};
  border-left: 3px solid ${({ theme, active }) => active ? theme.colors.primary : 'transparent'};
  
  &:hover {
    background: ${({ theme, active }) => active ? theme.colors.primaryLight : theme.colors.backgroundHover};
  }
`;

const ChannelName = styled.div`
  font-weight: 500;
`;

const UserList = styled.div`
  padding: 1rem 0;
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  max-height: 200px;
  overflow-y: auto;
`;

const UserItem = styled.div`
  padding: 0.5rem 1rem;
  display: flex;
  align-items: center;
`;

const UserStatus = styled.div<{ status: string }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 0.5rem;
  background: ${({ theme, status }) => 
    status === 'online' ? theme.colors.success :
    status === 'away' ? theme.colors.warning :
    theme.colors.error
  };
`;

const Username = styled.div`
  color: ${({ theme }) => theme.colors.text};
`;

const ChatArea = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

const ChatHeader = styled.div`
  padding: 1rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme }) => theme.colors.card};
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const ChatTitle = styled.h2`
  font-size: 1.2rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0;
`;

const ChatDescription = styled.p`
  font-size: 0.8rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin: 0.25rem 0 0;
`;

const MessageList = styled.div`
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
`;

const MessageContainer = styled.div`
  display: flex;
  margin-bottom: 0.5rem;
`;

const MessageAvatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  margin-right: 0.75rem;
  flex-shrink: 0;
`;

const MessageContent = styled.div`
  flex: 1;
`;

const MessageHeader = styled.div`
  display: flex;
  align-items: baseline;
  margin-bottom: 0.25rem;
`;

const MessageAuthor = styled.span`
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  margin-right: 0.5rem;
`;

const MessageTime = styled.span`
  font-size: 0.75rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const MessageText = styled.p`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  line-height: 1.4;
`;

const MessageInputContainer = styled.div`
  padding: 1rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  background: ${({ theme }) => theme.colors.card};
`;

const MessageForm = styled.form`
  display: flex;
  gap: 0.5rem;
`;

const MessageInput = styled.input`
  flex: 1;
  padding: 0.75rem;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  font-size: 1rem;
  background: ${({ theme }) => theme.colors.input};
  color: ${({ theme }) => theme.colors.text};
  
  &:focus {
    outline: none;
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const SendButton = styled.button`
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0 1rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
  }
`;

const EmptyState = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  text-align: center;
`;

const EmptyStateTitle = styled.h3`
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.text};
`;

const EmptyStateText = styled.p`
  margin-bottom: 1.5rem;
`;

// Mock data function that simulates WebSocket events
const createMockWebSocketEvents = (setChannels: React.Dispatch<React.SetStateAction<Channel[]>>, 
                                   setMessages: React.Dispatch<React.SetStateAction<Record<string, Message[]>>>,
                                   setUsers: React.Dispatch<React.SetStateAction<User[]>>) => {
  // Initial mock data
  setTimeout(() => {
    setChannels([
      { id: '1', name: 'general', description: 'General discussion', isPrivate: false },
      { id: '2', name: 'help', description: 'Get help with issues', isPrivate: false },
      { id: '3', name: 'dev-team', description: 'Development team discussions', isPrivate: true },
    ]);
    
    setUsers([
      { id: '1', username: 'jaden', status: 'online' },
      { id: '2', username: 'alice', status: 'online' },
      { id: '3', username: 'bob', status: 'away' },
      { id: '4', username: 'charlie', status: 'offline' },
    ]);
    
    setMessages({
      '1': [
        { 
          id: '1', 
          channelId: '1', 
          userId: '2', 
          username: 'alice', 
          content: 'Hello everyone! Welcome to the general channel', 
          timestamp: new Date(Date.now() - 3600000).toISOString()
        },
        { 
          id: '2', 
          channelId: '1', 
          userId: '3', 
          username: 'bob', 
          content: 'Thanks for setting this up!', 
          timestamp: new Date(Date.now() - 3000000).toISOString()
        },
      ],
      '2': [
        { 
          id: '3', 
          channelId: '2', 
          userId: '4', 
          username: 'charlie', 
          content: 'How do I connect to the websocket endpoint?', 
          timestamp: new Date(Date.now() - 1800000).toISOString()
        },
      ],
      '3': []
    });
  }, 1000);
  
  // Simulate incoming message every 15 seconds in general channel
  let counter = 10;
  const mockUsers = ['alice', 'bob', 'charlie'];
  const mockContents = [
    'Just checking in!',
    'How is everyone doing today?',
    'Making progress on the new feature.',
    'Can someone review my PR?',
    'Don\'t forget the meeting at 3pm!',
  ];
  
  const interval = setInterval(() => {
    const mockEvent: WebSocketEvent = {
      type: 'message',
      payload: {
        id: (++counter).toString(),
        channelId: '1',
        userId: Math.floor(Math.random() * 3 + 2).toString(),
        username: mockUsers[Math.floor(Math.random() * mockUsers.length)],
        content: mockContents[Math.floor(Math.random() * mockContents.length)],
        timestamp: new Date().toISOString()
      }
    };
    
    handleMockEvent(mockEvent);
  }, 15000);
  
  // Simulate user status changes every 20 seconds
  const statusInterval = setInterval(() => {
    const statuses: Array<'online' | 'offline' | 'away'> = ['online', 'offline', 'away'];
    const userId = Math.floor(Math.random() * 3 + 2).toString();
    const status = statuses[Math.floor(Math.random() * statuses.length)];
    
    const mockEvent: WebSocketEvent = {
      type: 'user_status_changed',
      payload: {
        userId,
        status
      }
    };
    
    handleMockEvent(mockEvent);
  }, 20000);
  
  function handleMockEvent(event: WebSocketEvent) {
    if (event.type === 'message') {
      setMessages(prev => ({
        ...prev,
        [event.payload.channelId]: [
          ...prev[event.payload.channelId],
          event.payload
        ]
      }));
    } else if (event.type === 'user_status_changed') {
      setUsers(prev => 
        prev.map(user => 
          user.id === event.payload.userId 
            ? { ...user, status: event.payload.status as 'online' | 'offline' | 'away' } 
            : user
        )
      );
    }
  }
  
  // Return cleanup function
  return () => {
    clearInterval(interval);
    clearInterval(statusInterval);
  };
};

// Main Component
const Messaging: React.FC = () => {
  // Ensure page scrolls to top when navigated to
  
  
  const [channels, setChannels] = useState<Channel[]>([]);
  const [selectedChannel, setSelectedChannel] = useState<Channel | null>(null);
  const [messages, setMessages] = useState<Record<string, Message[]>>({});
  const [users, setUsers] = useState<User[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [currentUser] = useState<User>({ id: '1', username: 'jaden', status: 'online' });
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const messageListRef = useRef<HTMLDivElement>(null);

  // Initialize WebSocket connection
  useEffect(() => {
    // Try to determine the WebSocket URL from environment or use fallback
    let wsUrl: string;
    
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      // Development environment
      wsUrl = 'ws://localhost:8082/ws';
    } else {
      // Production environment - use nginx proxy endpoint
      wsUrl = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/ws`;
    }
    
    try {
      const ws = new WebSocket(wsUrl);
      
      ws.onopen = () => {
        console.log('Connected to WebSocket');
        setSocket(ws);
        setConnected(true);
      };
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as WebSocketEvent;
          handleWebSocketEvent(data);
        } catch (e) {
          console.error('Error parsing WebSocket message:', e);
        }
      };
      
      ws.onerror = (error) => {
        console.log('WebSocket connection failed, using demo mode');
        setConnected(false);
        
        // If can't connect, use mock data
        const cleanup = createMockWebSocketEvents(setChannels, setMessages, setUsers);
        return cleanup;
      };
      
      ws.onclose = () => {
        console.log('WebSocket connection closed');
        setConnected(false);
        setSocket(null);
      };
      
      return () => {
        ws.close();
      };
    } catch (error) {
      console.log('Failed to initialize WebSocket, using demo mode');
      
      // Use mock data
      const cleanup = createMockWebSocketEvents(setChannels, setMessages, setUsers);
      return cleanup;
    }
  }, []);

  // Handle WebSocket events
  const handleWebSocketEvent = (event: WebSocketEvent) => {
    switch (event.type) {
      case 'message':
        setMessages(prev => ({
          ...prev,
          [event.payload.channelId]: [
            ...(prev[event.payload.channelId] || []),
            event.payload
          ]
        }));
        break;
        
      case 'channel_created':
        setChannels(prev => [...prev, event.payload]);
        break;
        
      case 'user_status_changed':
        setUsers(prev => 
          prev.map(user => 
            user.id === event.payload.userId 
              ? { ...user, status: event.payload.status as 'online' | 'offline' | 'away' } 
              : user
          )
        );
        break;
        
      case 'error':
        console.error('WebSocket error:', event.payload.message);
        break;
        
      default:
        console.warn('Unknown WebSocket event type:', event);
    }
  };

  // Auto-scroll messages when new ones are added
  useEffect(() => {
    if (messageListRef.current && selectedChannel) {
      messageListRef.current.scrollTo({
        top: messageListRef.current.scrollHeight,
        behavior: 'smooth'
      });
    }
  }, [messages, selectedChannel]);

  // Select the first channel when channels are loaded
  useEffect(() => {
    if (channels.length > 0 && !selectedChannel) {
      setSelectedChannel(channels[0]);
    }
  }, [channels, selectedChannel]);

  // Send a message
  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newMessage.trim() || !selectedChannel) return;
    
    const message: Message = {
      id: Math.random().toString(),
      channelId: selectedChannel.id,
      userId: currentUser.id,
      username: currentUser.username,
      content: newMessage,
      timestamp: new Date().toISOString()
    };
    
    if (socket && connected) {
      // Send through WebSocket
      socket.send(JSON.stringify({
        type: 'message',
        payload: message
      }));
    } else {
      // Add to local state if WebSocket is not connected
      setMessages(prev => ({
        ...prev,
        [selectedChannel.id]: [
          ...(prev[selectedChannel.id] || []),
          message
        ]
      }));
    }
    
    setNewMessage('');
  };

  // Format timestamp
  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return new Intl.DateTimeFormat('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    }).format(date);
  };

  // Get initials from username
  const getInitials = (username: string) => {
    return username.slice(0, 2).toUpperCase();
  };

  return (
    <PageContainer>
      <ChatContainer>
        <Sidebar>
          <SidebarHeader>
            <SidebarTitle>Channels</SidebarTitle>
          </SidebarHeader>
          
          <ChannelList>
            {channels.map(channel => (
              <ChannelItem 
                key={channel.id}
                active={selectedChannel?.id === channel.id}
                onClick={() => setSelectedChannel(channel)}
              >
                <ChannelName># {channel.name}</ChannelName>
              </ChannelItem>
            ))}
          </ChannelList>
          
          <UserList>
            {users.map(user => (
              <UserItem key={user.id}>
                <UserStatus status={user.status} />
                <Username>{user.username}</Username>
              </UserItem>
            ))}
          </UserList>
        </Sidebar>
        
        <ChatArea>
          {selectedChannel ? (
            <>
              <ChatHeader>
                <div>
                  <ChatTitle># {selectedChannel.name}</ChatTitle>
                  <ChatDescription>{selectedChannel.description}</ChatDescription>
                </div>
              </ChatHeader>
              
              <MessageList ref={messageListRef}>
                {messages[selectedChannel.id]?.length > 0 ? (
                  messages[selectedChannel.id].map(message => (
                    <MessageContainer key={message.id}>
                      <MessageAvatar>{getInitials(message.username)}</MessageAvatar>
                      <MessageContent>
                        <MessageHeader>
                          <MessageAuthor>{message.username}</MessageAuthor>
                          <MessageTime>{formatTimestamp(message.timestamp)}</MessageTime>
                        </MessageHeader>
                        <MessageText>{message.content}</MessageText>
                      </MessageContent>
                    </MessageContainer>
                  ))
                ) : (
                  <EmptyState>
                    <EmptyStateTitle>No messages yet</EmptyStateTitle>
                    <EmptyStateText>Be the first to send a message in this channel!</EmptyStateText>
                  </EmptyState>
                )}
              </MessageList>
              
              <MessageInputContainer>
                <MessageForm onSubmit={sendMessage}>
                  <MessageInput 
                    type="text"
                    placeholder={`Message #${selectedChannel.name}`}
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                  />
                  <SendButton type="submit">Send</SendButton>
                </MessageForm>
              </MessageInputContainer>
            </>
          ) : (
            <EmptyState>
              <EmptyStateTitle>Select a channel</EmptyStateTitle>
              <EmptyStateText>Choose a channel from the sidebar to start messaging</EmptyStateText>
            </EmptyState>
          )}
        </ChatArea>
      </ChatContainer>
    </PageContainer>
  );
};

export default Messaging; 