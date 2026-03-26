import apiClient from './api';
import type { WorldSetting, WorldSettingInput, WorldTemplate } from './types';

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
};
