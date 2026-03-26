import apiClient from './api';
import type { Character, CharacterInput } from './types';

export const characterService = {
  getCharacters: async (): Promise<Character[]> => {
    const response = await apiClient.get<Character[]>('/characters');
    return response.data;
  },
  createCharacter: async (character: CharacterInput): Promise<Character> => {
    const response = await apiClient.post<Character>('/characters', character);
    return response.data;
  },
};
