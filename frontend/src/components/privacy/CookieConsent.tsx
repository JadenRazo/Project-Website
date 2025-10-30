import React, { useState, useEffect } from 'react';
import { X, Shield, Settings, ChevronDown, ChevronUp, Info } from 'lucide-react';

interface ConsentCategory {
  id: string;
  name: string;
  description: string;
  required: boolean;
  enabled: boolean;
}

interface CookieConsentProps {
  onAccept: (categories: Record<string, boolean>) => void;
  onDecline: () => void;
  position?: 'top' | 'bottom' | 'center';
}

const CookieConsent: React.FC<CookieConsentProps> = ({
  onAccept,
  onDecline,
  position = 'bottom'
}) => {
  const [isVisible, setIsVisible] = useState(false);
  const [showDetails, setShowDetails] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [categories, setCategories] = useState<ConsentCategory[]>([
    {
      id: 'necessary',
      name: 'Necessary',
      description: 'Essential cookies required for the website to function properly. These cannot be disabled.',
      required: true,
      enabled: true,
    },
    {
      id: 'analytics',
      name: 'Analytics',
      description: 'Help us understand how visitors interact with our website to improve user experience.',
      required: false,
      enabled: false,
    },
    {
      id: 'functional',
      name: 'Functional',
      description: 'Enable enhanced functionality and personalization, such as remembering your preferences.',
      required: false,
      enabled: false,
    },
    {
      id: 'marketing',
      name: 'Marketing',
      description: 'Used to track visitors across websites to display relevant advertisements.',
      required: false,
      enabled: false,
    },
  ]);

  useEffect(() => {
    const consent = localStorage.getItem('cookieConsent');
    if (!consent) {
      setTimeout(() => setIsVisible(true), 1000);
    }
  }, []);

  const handleAcceptAll = () => {
    const allEnabled = categories.reduce((acc, cat) => {
      acc[cat.id] = true;
      return acc;
    }, {} as Record<string, boolean>);

    localStorage.setItem('cookieConsent', JSON.stringify({
      timestamp: new Date().toISOString(),
      categories: allEnabled,
    }));

    onAccept(allEnabled);
    setIsVisible(false);
  };

  const handleAcceptSelected = () => {
    const selected = categories.reduce((acc, cat) => {
      acc[cat.id] = cat.enabled;
      return acc;
    }, {} as Record<string, boolean>);

    localStorage.setItem('cookieConsent', JSON.stringify({
      timestamp: new Date().toISOString(),
      categories: selected,
    }));

    onAccept(selected);
    setIsVisible(false);
  };

  const handleRejectAll = () => {
    const onlyNecessary = categories.reduce((acc, cat) => {
      acc[cat.id] = cat.required;
      return acc;
    }, {} as Record<string, boolean>);

    localStorage.setItem('cookieConsent', JSON.stringify({
      timestamp: new Date().toISOString(),
      categories: onlyNecessary,
    }));

    onDecline();
    setIsVisible(false);
  };

  const toggleCategory = (id: string) => {
    setCategories(prev =>
      prev.map(cat =>
        cat.id === id && !cat.required
          ? { ...cat, enabled: !cat.enabled }
          : cat
      )
    );
  };

  if (!isVisible) return null;

  const positionClasses = {
    top: 'top-0',
    bottom: 'bottom-0',
    center: 'top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 max-w-2xl',
  };

  return (
    <>
      <div className="fixed inset-0 bg-black bg-opacity-50 z-40" onClick={() => setShowSettings(false)} />

      <div className={`fixed ${position === 'center' ? positionClasses.center : `${positionClasses[position]} left-0 right-0`} bg-white dark:bg-gray-900 shadow-2xl z-50 border border-gray-200 dark:border-gray-700 ${position !== 'center' ? 'md:rounded-t-lg' : 'rounded-lg'}`}>
        <div className="max-w-7xl mx-auto p-4 md:p-6">
          {!showSettings ? (
            <div className="space-y-4">
              <div className="flex items-start justify-between">
                <div className="flex items-center space-x-2">
                  <Shield className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                  <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                    Privacy & Cookie Settings
                  </h2>
                </div>
                <button
                  onClick={() => setIsVisible(false)}
                  className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              <p className="text-gray-600 dark:text-gray-300">
                We use cookies and similar technologies to help personalize content, tailor and measure ads,
                and provide a better experience. By clicking accept, you agree to this, as outlined in our
                <a href="/privacy-policy" className="text-blue-600 dark:text-blue-400 hover:underline ml-1">
                  Privacy Policy
                </a>.
              </p>

              <button
                onClick={() => setShowDetails(!showDetails)}
                className="flex items-center text-sm text-blue-600 dark:text-blue-400 hover:underline"
              >
                {showDetails ? 'Hide' : 'Show'} details
                {showDetails ? <ChevronUp className="h-4 w-4 ml-1" /> : <ChevronDown className="h-4 w-4 ml-1" />}
              </button>

              {showDetails && (
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4 space-y-3">
                  <h3 className="font-semibold text-gray-900 dark:text-white">What are cookies?</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Cookies are small text files that websites place on your device as you browse the web.
                    They help websites remember information about your visit, which can make it easier to visit
                    the site again and make the site more useful to you.
                  </p>

                  <h3 className="font-semibold text-gray-900 dark:text-white mt-4">Your privacy matters</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    We respect your privacy and give you control over your data. You can choose which types of
                    cookies you want to allow. Note that blocking some types of cookies may impact your experience
                    of the site and the services we are able to offer.
                  </p>

                  <div className="mt-4 space-y-2">
                    <div className="flex items-center text-sm">
                      <Info className="h-4 w-4 text-blue-600 dark:text-blue-400 mr-2" />
                      <span className="text-gray-600 dark:text-gray-300">
                        GDPR, CCPA, LGPD, and PIPEDA compliant
                      </span>
                    </div>
                    <div className="flex items-center text-sm">
                      <Info className="h-4 w-4 text-green-600 dark:text-green-400 mr-2" />
                      <span className="text-gray-600 dark:text-gray-300">
                        Your data is never sold to third parties
                      </span>
                    </div>
                    <div className="flex items-center text-sm">
                      <Info className="h-4 w-4 text-purple-600 dark:text-purple-400 mr-2" />
                      <span className="text-gray-600 dark:text-gray-300">
                        You can change your preferences anytime
                      </span>
                    </div>
                  </div>
                </div>
              )}

              <div className="flex flex-col md:flex-row gap-3">
                <button
                  onClick={handleRejectAll}
                  className="px-6 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
                >
                  Reject All
                </button>

                <button
                  onClick={() => setShowSettings(true)}
                  className="px-6 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors flex items-center justify-center"
                >
                  <Settings className="h-4 w-4 mr-2" />
                  Manage Preferences
                </button>

                <button
                  onClick={handleAcceptAll}
                  className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  Accept All
                </button>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                  Manage Cookie Preferences
                </h2>
                <button
                  onClick={() => setShowSettings(false)}
                  className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              <p className="text-gray-600 dark:text-gray-300">
                Choose which types of cookies you want to allow. You can change these settings at any time.
              </p>

              <div className="space-y-3 max-h-96 overflow-y-auto">
                {categories.map((category) => (
                  <div
                    key={category.id}
                    className="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h3 className="font-semibold text-gray-900 dark:text-white">
                          {category.name}
                          {category.required && (
                            <span className="ml-2 text-xs bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400 px-2 py-1 rounded">
                              Required
                            </span>
                          )}
                        </h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                          {category.description}
                        </p>
                      </div>

                      <label className="relative inline-flex items-center cursor-pointer ml-4">
                        <input
                          type="checkbox"
                          checked={category.enabled}
                          disabled={category.required}
                          onChange={() => toggleCategory(category.id)}
                          className="sr-only peer"
                        />
                        <div className={`w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 ${category.enabled ? 'peer-checked:bg-blue-600' : ''} ${category.required ? 'opacity-50 cursor-not-allowed' : 'peer-checked:after:translate-x-full'} peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600`}></div>
                      </label>
                    </div>
                  </div>
                ))}
              </div>

              <div className="flex justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
                <button
                  onClick={() => setShowSettings(false)}
                  className="px-6 py-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
                >
                  Back
                </button>

                <button
                  onClick={handleAcceptSelected}
                  className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  Save Preferences
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
};

export default CookieConsent;