import apiClient from './api';
import type { Character, CharacterInput, Relationship, RelationshipInput, CharacterArc } from './types';

export const characterService = {
  getCharacters: async (params?: { camp?: string; search?: string }): Promise<Character[]> => {
    const response = await apiClient.get<Character[]>('/characters', { params });
    return response.data;
  },
  createCharacter: async (character: CharacterInput): Promise<Character> => {
    const response = await apiClient.post<Character>('/characters', character);
    return response.data;
  },
  getCharacterById: async (id: string): Promise<Character> => {
    const response = await apiClient.get<Character>(`/characters/${id}`);
    return response.data;
  },
  updateCharacter: async (id: string, character: CharacterInput): Promise<void> => {
    await apiClient.put(`/characters/${id}`, character);
  },
  getRelationships: async (): Promise<Relationship[]> => {
    const response = await apiClient.get<Relationship[]>('/characters/relationships');
    return response.data;
  },
  createRelationship: async (rel: RelationshipInput): Promise<void> => {
    await apiClient.post('/characters/relationships', rel);
  },
  getCharacterArcs: async (id: string): Promise<CharacterArc[]> => {
    const response = await apiClient.get<CharacterArc[]>(`/characters/${id}/arcs`);
    return response.data;
  },
};
