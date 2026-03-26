export interface Character {
  id: string;
  name: string;
  role: string;
  pov_type?: string;
  created_at: string;
}

export interface CharacterInput {
  name: string;
  role: string;
  pov_type?: string;
}

export interface PlotCard {
  id: string;
  title: string;
  description?: string;
  conflict_intensity: number;
  order_index: number;
  created_at: string;
}

export interface PlotCardInput {
  title: string;
  description?: string;
  conflict_intensity: number;
  order_index: number;
}

export interface WorldSetting {
  id: string;
  category: string;
  name: string;
  description?: string;
  logic_rules?: string;
  created_at: string;
}

export interface WorldSettingInput {
  category: string;
  name: string;
  description?: string;
  logic_rules?: string;
}

export interface Scene {
  id: string;
  title: string;
  plot_card_id: string;
  pov_character_id?: string;
  goal?: string;
  conflict?: string;
  resolution?: string;
  created_at: string;
}

export interface SceneInput {
  title: string;
  plot_card_id: string;
  pov_character_id?: string;
  goal?: string;
  conflict?: string;
  resolution?: string;
}

export interface TimelineEvent {
  id: string;
  title: string;
  description?: string;
  chronological_order: number;
  scene_id?: string;
  created_at: string;
}

export interface TimelineEventInput {
  title: string;
  description?: string;
  chronological_order: number;
  scene_id?: string;
}

export interface WorldTemplate {
  id: string;
  category: string;
  name: string;
  description?: string;
  suggested_logic?: string;
}

export interface AuditResult {
  is_valid: boolean;
  issues: string[];
}
