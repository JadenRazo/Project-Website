import React from 'react';
import { useLocation } from 'react-router-dom';
import HeroSkeleton from './HeroSkeleton';
import AboutSkeleton from './AboutSkeleton';
import ProjectsSkeleton from './ProjectsSkeleton';
import SkillsSkeleton from './SkillsSkeleton';
import ContactSkeleton from './ContactSkeleton';
import GenericPageSkeleton from './GenericPageSkeleton';
import BlogSkeleton from './BlogSkeleton';
import BlogPostSkeleton from './BlogPostSkeleton';

interface SmartSkeletonProps {
  type?: 'hero' | 'about' | 'projects' | 'skills' | 'contact' | 'blog' | 'blogpost' | 'generic' | 'auto';
}

const SmartSkeleton: React.FC<SmartSkeletonProps> = ({ type = 'auto' }) => {
  const location = useLocation();
  
  const getSkeletonType = (): string => {
    if (type !== 'auto') return type;
    
    const path = location.pathname.toLowerCase();
    
    if (path === '/') return 'hero';
    if (path.includes('/about')) return 'about';
    if (path.includes('/project') || path.includes('/portfolio')) return 'projects';
    if (path.includes('/skill')) return 'skills';
    if (path.includes('/contact')) return 'contact';
    if (path.match(/^\/blog\/[^/]+$/)) return 'blogpost';
    if (path.includes('/blog')) return 'blog';

    return 'generic';
  };
  
  const skeletonType = getSkeletonType();
  
  switch (skeletonType) {
    case 'hero':
      return <HeroSkeleton />;
    case 'about':
      return <AboutSkeleton />;
    case 'projects':
      return <ProjectsSkeleton />;
    case 'skills':
      return <SkillsSkeleton />;
    case 'contact':
      return <ContactSkeleton />;
    case 'blog':
      return <BlogSkeleton />;
    case 'blogpost':
      return <BlogPostSkeleton />;
    default:
      return <GenericPageSkeleton />;
  }
};

export default SmartSkeleton;