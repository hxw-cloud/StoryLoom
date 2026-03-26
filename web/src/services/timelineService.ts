import apiClient from './api';
import { TimelineEvent, TimelineEventInput } from './types';

export const timelineService = {
  getEvents: async (): Promise<TimelineEvent[]> => {
    const response = await apiClient.get<TimelineEvent[]>('/timeline/events');
    return response.data;
  },
  createEvent: async (event: TimelineEventInput): Promise<TimelineEvent> => {
    const response = await apiClient.post<TimelineEvent>('/timeline/events', event);
    return response.data;
  },
};
