import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { worldService } from '../../services/worldService';
import { characterService } from '../../services/characterService';
import type { WorldSetting, WorldSettingInput, WorldTemplate, HistoricalEvent, HistoricalEventInput, WorldAuditData, Character } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './WorldSettings.module.css';

const WorldSettings: React.FC = () => {
  const { t } = useTranslation();
  const [settings, setSettings] = useState<WorldSetting[]>([]);
  const [templates, setTemplates] = useState<WorldTemplate[]>([]);
  const [history, setHistory] = useState<HistoricalEvent[]>([]);
  const [audit, setAudit] = useState<WorldAuditData | null>(null);
  const [characters, setCharacters] = useState<Character[]>([]);
  
  const [activeTab, setActiveTab] = useState<'custom' | 'templates' | 'history' | 'audit'>('custom');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Search & Filter
  const [searchTerm, setSearchTerm] = useState('');
  const [categoryFilter, setCategoryFilter] = useState('');

  // Modals
  const [isSettingModalOpen, setIsSettingModalOpen] = useState(false);
  const [isHistoryModalOpen, setIsHistoryModalOpen] = useState(false);

  const [newSetting, setNewSetting] = useState<WorldSettingInput>({
    category: 'Geography',
    name: '',
    description: '',
    logic_rules: '',
    tags: [],
  });

  const [newHistory, setNewHistory] = useState<HistoricalEventInput>({
    title: '',
    event_time: '',
    impact_scope: '',
    involved_characters: [],
    cause: '',
    effect: '',
    is_iceberg_tip: false,
  });

  const fetchData = async () => {
    try {
      const [settingsData, templatesData, historyData, auditData, charData] = await Promise.all([
        worldService.getSettings(),
        worldService.getTemplates(),
        worldService.getHistory(),
        worldService.getAudit(),
        characterService.getCharacters(),
      ]);
      setSettings(settingsData);
      setTemplates(templatesData);
      setHistory(historyData);
      setAudit(auditData);
      setCharacters(charData);
    } catch (err) {
      setError(t('world.error'));
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleSettingSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await worldService.createSetting(newSetting);
      setIsSettingModalOpen(false);
      setNewSetting({ category: 'Geography', name: '', description: '', logic_rules: '', tags: [] });
      fetchData();
    } catch (err) {
      alert(t('world.addError'));
    }
  };

  const handleHistorySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await worldService.createHistory(newHistory);
      setIsHistoryModalOpen(false);
      setNewHistory({ title: '', event_time: '', impact_scope: '', involved_characters: [], cause: '', effect: '', is_iceberg_tip: false });
      fetchData();
    } catch (err) {
      alert(t('world.history.addError'));
    }
  };

  const useTemplate = (template: WorldTemplate) => {
    setNewSetting({
      category: template.category,
      name: template.name,
      description: template.description || '',
      logic_rules: template.suggested_logic || '',
      tags: [],
    });
    setIsSettingModalOpen(true);
  };

  const filteredSettings = settings.filter(s => {
    const matchesSearch = s.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                         s.description?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = categoryFilter === '' || s.category === categoryFilter;
    return matchesSearch && matchesCategory;
  });

  if (loading) return <div className={styles.loading}>{t('common.loading')}</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h2>{t('world.title')}</h2>
          <div className={styles.tabs}>
            <button className={`${styles.tab} ${activeTab === 'custom' ? styles.activeTab : ''}`} onClick={() => setActiveTab('custom')}>{t('world.tabs.custom')}</button>
            <button className={`${styles.tab} ${activeTab === 'templates' ? styles.activeTab : ''}`} onClick={() => setActiveTab('templates')}>{t('world.tabs.templates')}</button>
            <button className={`${styles.tab} ${activeTab === 'history' ? styles.activeTab : ''}`} onClick={() => setActiveTab('history')}>{t('world.tabs.history')}</button>
            <button className={`${styles.tab} ${activeTab === 'audit' ? styles.activeTab : ''}`} onClick={() => setActiveTab('audit')}>{t('world.tabs.audit')}</button>
          </div>
        </div>
        <div className={styles.actions}>
          {activeTab === 'custom' && <button className={styles.addButton} onClick={() => setIsSettingModalOpen(true)}>{t('world.newSetting')}</button>}
          {activeTab === 'history' && <button className={styles.addButton} onClick={() => setIsHistoryModalOpen(true)}>{t('world.history.newEvent')}</button>}
        </div>
      </header>

      {activeTab === 'custom' && (
        <div className={styles.searchBar}>
          <input 
            type="text" 
            placeholder={t('world.placeholders.search')} 
            value={searchTerm} 
            onChange={e => setSearchTerm(e.target.value)} 
          />
          <select value={categoryFilter} onChange={e => setCategoryFilter(e.target.value)}>
            <option value="">{t('world.categories.all')}</option>
            <option value="Geography">{t('world.categories.geography')}</option>
            <option value="Magic System">{t('world.categories.magic')}</option>
            <option value="Technology">{t('world.categories.tech')}</option>
            <option value="Race">{t('world.categories.race')}</option>
            <option value="Culture">{t('world.categories.culture')}</option>
          </select>
        </div>
      )}
      
      <div className={styles.grid}>
        {activeTab === 'custom' && (
          filteredSettings.length === 0 ? (
            <div className={styles.empty}>{t('world.empty')}</div>
          ) : (
            filteredSettings.map(setting => (
              <div key={setting.id} className={styles.card}>
                <div className={styles.categoryBadge}>{setting.category}</div>
                <h3>{setting.name}</h3>
                <p className={styles.description}>{setting.description}</p>
                {setting.tags && setting.tags.length > 0 && (
                  <div className={styles.tags}>
                    {setting.tags.map(tag => <span key={tag} className={styles.tag}>#{tag}</span>)}
                  </div>
                )}
                <div className={styles.usageInfo}>{t('world.usage')}: {setting.usage_count}</div>
              </div>
            ))
          )
        )}

        {activeTab === 'templates' && (
          templates.map(template => (
            <div key={template.id} className={`${styles.card} ${styles.templateCard}`}>
              <div className={styles.categoryBadge}>{template.category}</div>
              <h3>{template.name}</h3>
              <p className={styles.description}>{template.description}</p>
              <button className={styles.useTemplateButton} onClick={() => useTemplate(template)}>{t('common.use')}</button>
            </div>
          ))
        )}

        {activeTab === 'history' && (
          <div className={styles.historyList}>
            {history.map(event => (
              <div key={event.id} className={styles.historyItem}>
                <div className={styles.historyMeta}>
                  <span className={styles.eventTime}>{event.event_time}</span>
                  {event.is_iceberg_tip && <span className={styles.icebergBadge}>Iceberg Tip</span>}
                </div>
                <h3>{event.title}</h3>
                <p>{event.impact_scope}</p>
                <div className={styles.historyDetails}>
                  <div><strong>Cause:</strong> {event.cause}</div>
                  <div><strong>Effect:</strong> {event.effect}</div>
                </div>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'audit' && audit && (
          <div className={styles.auditDashboard}>
            <div className={styles.auditStat}>
              <h4>Iceberg Ratio (Used vs Total)</h4>
              <div className={styles.progressBar}>
                <div className={styles.progressFill} style={{width: `${audit.iceberg_ratio * 100}%`}} />
              </div>
              <span>{(audit.iceberg_ratio * 100).toFixed(1)}%</span>
            </div>
            <div className={styles.intensityMap}>
              <h4>Usage Intensity Map</h4>
              {Object.entries(audit.intensity_map).map(([name, count]) => (
                <div key={name} className={styles.intensityRow}>
                  <span>{name}</span>
                  <div className={styles.intensityBar} style={{width: `${Math.min(count * 20, 100)}%`, opacity: 0.5 + (count * 0.1)}} />
                  <span>{count}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      <Modal isOpen={isSettingModalOpen} onClose={() => setIsSettingModalOpen(false)} title={t('world.newSetting')}>
        <form onSubmit={handleSettingSubmit} className={styles.form}>
          <div className={styles.formGroup}><label>{t('world.category')}</label>
            <select value={newSetting.category} onChange={e => setNewSetting({...newSetting, category: e.target.value})}>
              <option value="Geography">{t('world.categories.geography')}</option>
              <option value="Magic System">{t('world.categories.magic')}</option>
              <option value="Technology">{t('world.categories.tech')}</option>
              <option value="Race">{t('world.categories.race')}</option>
              <option value="Culture">{t('world.categories.culture')}</option>
            </select>
          </div>
          <div className={styles.formGroup}><label>{t('world.name')}</label><input type="text" required value={newSetting.name} onChange={e => setNewSetting({...newSetting, name: e.target.value})} placeholder={t('world.placeholders.name')}/></div>
          <div className={styles.formGroup}><label>{t('world.description')}</label><textarea value={newSetting.description} onChange={e => setNewSetting({...newSetting, description: e.target.value})} placeholder={t('world.placeholders.description')}/></div>
          <div className={styles.formGroup}><label>Tags (comma separated)</label><input type="text" value={newSetting.tags?.join(',')} onChange={e => setNewSetting({...newSetting, tags: e.target.value.split(',')})} placeholder="e.g. ancient, secret, dangerous"/></div>
          <button type="submit" className={styles.submitButton}>{t('common.save')}</button>
        </form>
      </Modal>

      <Modal isOpen={isHistoryModalOpen} onClose={() => setIsHistoryModalOpen(false)} title="Record Historical Event">
        <form onSubmit={handleHistorySubmit} className={styles.form}>
          <div className={styles.formGroup}><label>Title</label><input type="text" required value={newHistory.title} onChange={e => setNewHistory({...newHistory, title: e.target.value})}/></div>
          <div className={styles.formGroup}><label>Time</label><input type="text" required value={newHistory.event_time} onChange={e => setNewHistory({...newHistory, event_time: e.target.value})} placeholder="e.g. Year 120 AF"/></div>
          <div className={styles.formGroup}><label>Impact Scope</label><input type="text" value={newHistory.impact_scope} onChange={e => setNewHistory({...newHistory, impact_scope: e.target.value})}/></div>
          <div className={styles.formGroup}>
            <label>Involved Characters</label>
            <select multiple value={newHistory.involved_characters} onChange={e => setNewHistory({...newHistory, involved_characters: Array.from(e.target.selectedOptions, option => option.value)})}>
              {characters.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
            </select>
          </div>
          <div className={styles.formGroup}><label>Cause</label><textarea value={newHistory.cause} onChange={e => setNewHistory({...newHistory, cause: e.target.value})}/></div>
          <div className={styles.formGroup}><label>Effect</label><textarea value={newHistory.effect} onChange={e => setNewHistory({...newHistory, effect: e.target.value})}/></div>
          <div className={styles.formGroup}><label><input type="checkbox" checked={newHistory.is_iceberg_tip} onChange={e => setNewHistory({...newHistory, is_iceberg_tip: e.target.checked})}/> Is Iceberg Tip (Visible to readers?)</label></div>
          <button type="submit" className={styles.submitButton}>{t('common.save')}</button>
        </form>
      </Modal>
    </div>
  );
};

export default WorldSettings;
