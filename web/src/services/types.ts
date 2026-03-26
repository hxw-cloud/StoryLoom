export interface Character {
  id: string;
  name: string;
  age?: number;
  gender?: string;
  role: string;
  camp?: string;
  appearance?: string;
  background?: string;
  pov_type?: string;
  image_url?: string;
  want?: string;
  need?: string;
  persona_template?: string;
  created_at: string;
}

export interface CharacterInput {
  name: string;
  age?: number;
  gender?: string;
  role: string;
  camp?: string;
  appearance?: string;
  background?: string;
  pov_type?: string;
  image_url?: string;
  want?: string;
  need?: string;
  persona_template?: string;
}

export interface Relationship {
  source_id: string;
  target_id: string;
  type: string;
  description?: string;
}

export interface RelationshipInput {
  source_id: string;
  target_id: string;
  type: string;
  description?: string;
}

export interface CharacterArc {
  character_id: string;
  plot_card_id: string;
  state_change: string;
  internal_growth: number;
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
  tags?: string[];
  usage_count?: number;
  created_at: string;
}

export interface WorldSettingInput {
  category: string;
  name: string;
  description?: string;
  logic_rules?: string;
  tags?: string[];
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

export interface HistoricalEvent {
  id: string;
  title: string;
  event_time: string;
  impact_scope?: string;
  involved_characters: string[];
  cause?: string;
  effect?: string;
  is_iceberg_tip: boolean;
  created_at: string;
}

export interface HistoricalEventInput {
  title: string;
  event_time: string;
  impact_scope?: string;
  involved_characters: string[];
  cause?: string;
  effect?: string;
  is_iceberg_tip: boolean;
}

export interface WorldAuditData {
  intensity_map: Record<string, number>;
  iceberg_ratio: number;
}

export interface AuditResult {
  is_valid: boolean;
  issues: string[];
}
