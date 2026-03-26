import apiClient from './api';
import type { WorldSetting, WorldSettingInput, WorldTemplate, HistoricalEvent, HistoricalEventInput, WorldAuditData } from './types';

export const worldService = {
  getSettings: async (): Promise<WorldSetting[]> => {
    const response = await apiClient.get<WorldSetting[]>('/world/settings');
    return response.data;
  },
  createSetting: async (setting: WorldSettingInput): Promise<WorldSetting> => {
    const response = await apiClient.post<WorldSetting>('/world/settings', setting);
    return response.data;
  },
  getTemplates: async (): Promise<WorldTemplate[]> => {
    const response = await apiClient.get<WorldTemplate[]>('/world/templates');
    return response.data;
  },
  getHistory: async (): Promise<HistoricalEvent[]> => {
    const response = await apiClient.get<HistoricalEvent[]>('/world/history');
    return response.data;
  },
  createHistory: async (event: HistoricalEventInput): Promise<HistoricalEvent> => {
    const response = await apiClient.post<HistoricalEvent>('/world/history', event);
    return response.data;
  },
  getAudit: async (): Promise<WorldAuditData> => {
    const response = await apiClient.get<WorldAuditData>('/world/audit');
    return response.data;
  },
};
