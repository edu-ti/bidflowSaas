// hooks/useFetch.ts
import { useState, useEffect } from 'react';
import api from '@/lib/api';

/**
 * Generic data fetching hook using the central Axios instance.
 * Returns typed data, loading state and error message.
 */
export function useFetch<T>(url: string, params?: Record<string, any>) {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isMounted = true;
    const fetchData = async () => {
      try {
        const response = await api.get<T>(url, { params });
        if (isMounted) {
          setData(response.data);
        }
      } catch (err: any) {
        if (isMounted) {
          setError(err?.message ?? 'Error');
        }
      } finally {
        if (isMounted) setLoading(false);
      }
    };

    fetchData();

    return () => {
      isMounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url, JSON.stringify(params)]);

  return { data, loading, error };
}
