import axios from 'axios';

const getApiBaseUrl = (): string => {
  const runtimeEnv = (window as any)._env_?.REACT_APP_API_URL;
  const buildTimeEnv = process.env.REACT_APP_API_URL;

  if (runtimeEnv !== undefined && runtimeEnv !== null) {
    return runtimeEnv === '' ? '' : runtimeEnv;
  }

  if (buildTimeEnv && buildTimeEnv !== '') {
    return buildTimeEnv;
  }

  return '';
};

interface TrackingData {
  path: string;
  referrer?: string;
  sessionId?: string;
}

let isTrackingEnabled = true;
let trackingTimeout: NodeJS.Timeout | null = null;
let consecutiveFailures = 0;
const MAX_FAILURES = 3;

export const trackPageView = async (path: string): Promise<void> => {
  if (!isTrackingEnabled) return;

  if (consecutiveFailures >= MAX_FAILURES) {
    console.error('[Visitor Tracking] Disabled after', MAX_FAILURES, 'consecutive failures');
    isTrackingEnabled = false;
    return;
  }

  if (trackingTimeout) {
    clearTimeout(trackingTimeout);
  }

  trackingTimeout = setTimeout(async () => {
    try {
      const apiBaseUrl = getApiBaseUrl();
      const trackingUrl = apiBaseUrl ? `${apiBaseUrl}/api/v1/visitor/track` : '/api/v1/visitor/track';

      const trackingData: TrackingData = {
        path,
        referrer: document.referrer || undefined,
      };

      const response = await axios.post(
        trackingUrl,
        trackingData,
        {
          timeout: 5000,
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      consecutiveFailures = 0;

      if (response.data?.activeVisitors !== undefined) {
        console.debug(`[Visitor Tracking] ✓ Tracked: ${path} | ${response.data.activeVisitors} active visitors`);
      }
    } catch (error: any) {
      consecutiveFailures++;
      const errorMsg = error?.response?.data?.error || error?.message || 'Unknown error';
      console.error(`[Visitor Tracking] ✗ Failed to track ${path}:`, errorMsg, `(${consecutiveFailures}/${MAX_FAILURES})`);

      if (error?.response?.status === 404) {
        console.error('[Visitor Tracking] Endpoint not found - check nginx configuration');
      }
    }
  }, 300);
};

export const setTrackingEnabled = (enabled: boolean): void => {
  isTrackingEnabled = enabled;
};

export const getTrackingStatus = (): boolean => {
  return isTrackingEnabled;
};
