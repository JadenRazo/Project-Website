# Scroll System Documentation

## Overview
The portfolio website includes a comprehensive smooth scroll system that ensures users have optimal visibility of important prompts, forms, and interactive elements across all devices.

## Core Components

### 1. useScrollTo Hook
The main hook for programmatic scrolling:
```typescript
const { scrollToElement, scrollToTop, scrollToPosition, scrollToId } = useScrollTo();
```

### 2. Scroll Configuration
Located in `utils/scrollConfig.ts`:
- Default header offset: 80px (desktop), 60px (mobile)
- Smooth scrolling with fallback for older browsers
- Respects `prefers-reduced-motion` preference

### 3. Specialized Hooks

#### useScrollToForm
For dedicated form components:
```typescript
const { formRef, scrollToForm, scrollToFormWithDelay } = useScrollToForm();
```

#### useInlineFormScroll
For inline forms that appear conditionally:
```typescript
const { formRef, scrollToForm, triggerScroll } = useInlineFormScroll(isVisible);
```

## Implementation Examples

### 1. URL Shortener Auth Prompt
```typescript
const authPromptRef = useRef<HTMLDivElement>(null);
const { scrollToElement } = useScrollTo();

// In the button onClick handler:
scrollToElement(authPromptRef.current, { behavior: 'smooth' });
setTimeout(() => {
  setShowAuthModal(true);
}, 300);
```

### 2. For Modal Forms
The `ScrollableModal` component already handles scroll behavior internally. Just ensure the modal is used for forms that need attention.

### 3. For Inline Forms
Use the `useInlineFormScroll` hook when forms appear conditionally:
```typescript
const { formRef } = useInlineFormScroll(showForm, { 
  scrollOffset: 100,
  scrollDelay: 200 
});

return showForm && <form ref={formRef}>...</form>;
```

## Best Practices

1. **Timing**: Add a 200-300ms delay before showing modals after scrolling
2. **Offsets**: Use `scroll-margin-top` CSS property for fine-tuning
3. **Mobile**: Test thoroughly on mobile devices with different viewport sizes
4. **Accessibility**: The system respects reduced motion preferences automatically

## Where to Apply

1. **Authentication Prompts**: Already implemented in URL Shortener
2. **Form Validations**: When showing error messages - Implemented in all DevPanel forms
3. **Success Messages**: After form submission - Use NotificationSystem
4. **Dynamic Content**: When new sections appear on user interaction
5. **Navigation**: For section-to-section navigation

## CSS Support

Add `scroll-margin-top` to elements that need custom scroll positioning:
```css
.my-form {
  scroll-margin-top: 100px; /* Accounts for fixed headers */
}
```

## Notification System

A comprehensive notification system has been implemented for success/error messages:

### Usage Example
```typescript
import { NotificationSystem, useNotifications } from '../components/notifications';

const MyComponent = () => {
  const { notifications, addNotification, dismissNotification } = useNotifications();
  
  const handleSuccess = () => {
    addNotification({
      type: 'success',
      title: 'Success!',
      message: 'Your changes have been saved.',
      duration: 5000,
      scrollToNotification: true // Optional: scroll to notification
    });
  };
  
  return (
    <>
      <NotificationSystem 
        notifications={notifications}
        onDismiss={dismissNotification}
        position="top-right"
      />
      {/* Your component content */}
    </>
  );
};
```

### Features
- 4 types: success, error, warning, info
- Auto-dismiss with duration
- Optional scroll-to-notification
- Multiple position options
- Mobile responsive
- Smooth animations