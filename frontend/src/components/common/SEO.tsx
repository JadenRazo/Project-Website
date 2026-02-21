import React from 'react';
import { Helmet } from 'react-helmet-async';

interface SEOProps {
  title: string;
  description: string;
  path?: string;
  image?: string;
  type?: string;
}

const SEO: React.FC<SEOProps> = ({
  title,
  description,
  path = '',
  image = 'https://jadenrazo.dev/images/og-image.png',
  type = 'website'
}) => {
  const siteUrl = 'https://jadenrazo.dev';
  const canonicalUrl = `${siteUrl}${path}`;
  const siteName = 'Jaden Razo';

  const structuredData = {
    "@context": "https://schema.org",
    "@type": "Person",
    "name": "Jaden Razo",
    "jobTitle": "Full Stack Developer",
    "url": "https://jadenrazo.dev",
    "sameAs": [
      "https://github.com/JadenRazo"
    ]
  };

  return (
    <Helmet>
      <title>{title}</title>
      <meta name="description" content={description} />
      <link rel="canonical" href={canonicalUrl} />

      <meta property="og:type" content={type} />
      <meta property="og:title" content={title} />
      <meta property="og:description" content={description} />
      <meta property="og:url" content={canonicalUrl} />
      <meta property="og:image" content={image} />
      <meta property="og:image:width" content="1200" />
      <meta property="og:image:height" content="630" />
      <meta property="og:site_name" content={siteName} />
      <meta property="og:locale" content="en_US" />

      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={title} />
      <meta name="twitter:description" content={description} />
      <meta name="twitter:image" content={image} />
      <meta name="twitter:creator" content="@JadenRazo" />
      <meta name="twitter:site" content="@JadenRazo" />

      <script type="application/ld+json">
        {JSON.stringify(structuredData)}
      </script>
    </Helmet>
  );
};

export default SEO;
