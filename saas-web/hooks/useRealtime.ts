// hooks/useRealtime.ts
'use client';

import { useState, useEffect } from 'react';
import api from '@/lib/api';

export type MonitorEvent = {
  id: string;
  type: string;
  severity: string;
  message: string;
  timestamp: string;
  metadata?: any;
};

/**
 * Hook to poll for realtime events from the monitor.
 */
export function useRealtime(pollingIntervalMs = 5000) {
  const [events, setEvents] = useState<MonitorEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let timeoutId: NodeJS.Timeout;
    let isMounted = true;

    const fetchEvents = async () => {
      try {
        const response = await api.get('/api/v1/monitor/events');
        if (isMounted) {
          // Assuming the API returns an array or an object with events array
          const newEvents = Array.isArray(response.data) ? response.data : response.data?.events || [];
          setEvents(newEvents);
          setError(null);
        }
      } catch (err: any) {
        if (isMounted) {
          setError(err?.message || 'Failed to fetch events');
        }
      } finally {
        if (isMounted) {
          setLoading(false);
          timeoutId = setTimeout(fetchEvents, pollingIntervalMs);
        }
      }
    };

    fetchEvents();

    return () => {
      isMounted = false;
      clearTimeout(timeoutId);
    };
  }, [pollingIntervalMs]);

  return { events, loading, error };
}
