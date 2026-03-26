import apiClient from './api';
import type { Scene, SceneInput, AuditResult } from './types';

export const sceneService = {
  getScenes: async (): Promise<Scene[]> => {
    const response = await apiClient.get<Scene[]>('/scenes');
    return response.data;
  },
  createScene: async (scene: SceneInput): Promise<Scene> => {
    const response = await apiClient.post<Scene>('/scenes', scene);
    return response.data;
  },
  auditScene: async (sceneId: string): Promise<AuditResult> => {
    const response = await apiClient.get<AuditResult>(`/audit/scene/${sceneId}`);
    return response.data;
  },
};
