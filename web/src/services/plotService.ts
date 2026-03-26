import apiClient from './api';
import { PlotCard, PlotCardInput } from './types';

export const plotService = {
  getPlots: async (): Promise<PlotCard[]> => {
    const response = await apiClient.get<PlotCard[]>('/plots');
    return response.data;
  },
  createPlot: async (plot: PlotCardInput): Promise<PlotCard> => {
    const response = await apiClient.post<PlotCard>('/plots', plot);
    return response.data;
  },
};
