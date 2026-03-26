import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { worldService } from '../../services/worldService';
import type { WorldSetting, WorldSettingInput, WorldTemplate } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './WorldSettings.module.css';

const WorldSettings: React.FC = () => {
  const { t } = useTranslation();
  const [settings, setSettings] = useState<WorldSetting[]>([]);
  const [templates, setTemplates] = useState<WorldTemplate[]>([]);
  const [activeTab, setActiveTab] = useState<'custom' | 'templates'>('custom');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newSetting, setNewSetting] = useState<WorldSettingInput>({
    category: 'Geography',
    name: '',
    description: '',
    logic_rules: '',
  });

  const fetchData = async () => {
    try {
      const [settingsData, templatesData] = await Promise.all([
        worldService.getSettings(),
        worldService.getTemplates()
      ]);
      setSettings(settingsData);
      setTemplates(templatesData);
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

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await worldService.createSetting(newSetting);
      setIsModalOpen(false);
      setNewSetting({ category: 'Geography', name: '', description: '', logic_rules: '' });
      const data = await worldService.getSettings();
      setSettings(data);
    } catch (err) {
      alert(t('world.addError'));
    }
  };

  const useTemplate = (template: WorldTemplate) => {
    setNewSetting({
      category: template.category,
      name: template.name,
      description: template.description || '',
      logic_rules: template.suggested_logic || '',
    });
    setIsModalOpen(true);
  };

  if (loading) return <div className={styles.loading}>{t('common.loading')}</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h2>{t('world.title')}</h2>
          <div className={styles.tabs}>
            <button 
              className={`${styles.tab} ${activeTab === 'custom' ? styles.activeTab : ''}`}
              onClick={() => setActiveTab('custom')}
            >
              {t('world.title')}
            </button>
            <button 
              className={`${styles.tab} ${activeTab === 'templates' ? styles.activeTab : ''}`}
              onClick={() => setActiveTab('templates')}
            >
              {t('world.templates')}
            </button>
          </div>
        </div>
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>{t('world.newSetting')}</button>
      </header>
      
      <div className={styles.grid}>
        {activeTab === 'custom' ? (
          settings.length === 0 ? (
            <div className={styles.empty}>{t('world.empty')}</div>
          ) : (
            settings.map(setting => (
              <div key={setting.id} className={styles.card}>
                <div className={styles.categoryBadge}>{setting.category}</div>
                <h3>{setting.name}</h3>
                <p className={styles.description}>{setting.description}</p>
                {setting.logic_rules && (
                  <div className={styles.logicRules}>
                    <strong>{t('world.logicRules')}:</strong>
                    <p>{setting.logic_rules}</p>
                  </div>
                )}
              </div>
            ))
          )
        ) : (
          templates.map(template => (
            <div key={template.id} className={`${styles.card} ${styles.templateCard}`}>
              <div className={styles.categoryBadge}>{template.category}</div>
              <h3>{template.name}</h3>
              <p className={styles.description}>{template.description}</p>
              <button 
                className={styles.useTemplateButton}
                onClick={() => useTemplate(template)}
              >
                {t('common.use')}
              </button>
            </div>
          ))
        )}
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title={t('world.newSetting')}
      >
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label>{t('world.category')}</label>
            <select 
              value={newSetting.category} 
              onChange={e => setNewSetting({...newSetting, category: e.target.value})}
            >
              <option value="Geography">{t('world.categories.geography')}</option>
              <option value="Magic System">{t('world.categories.magic')}</option>
              <option value="Technology">{t('world.categories.tech')}</option>
              <option value="Race">{t('world.categories.race')}</option>
              <option value="Culture">{t('world.categories.culture')}</option>
            </select>
          </div>
          <div className={styles.formGroup}>
            <label>{t('world.name')}</label>
            <input 
              type="text" 
              required 
              value={newSetting.name}
              onChange={e => setNewSetting({...newSetting, name: e.target.value})}
              placeholder={t('world.placeholders.name')}
            />
          </div>
          <div className={styles.formGroup}>
            <label>{t('world.description')}</label>
            <textarea 
              value={newSetting.description}
              onChange={e => setNewSetting({...newSetting, description: e.target.value})}
              placeholder={t('world.placeholders.description')}
            />
          </div>
          <div className={styles.formGroup}>
            <label>{t('world.logicRules')}</label>
            <textarea 
              value={newSetting.logic_rules}
              onChange={e => setNewSetting({...newSetting, logic_rules: e.target.value})}
              placeholder={t('world.placeholders.logicRules')}
            />
          </div>
          <button type="submit" className={styles.submitButton}>{t('common.save')}</button>
        </form>
      </Modal>
    </div>
  );
};

export default WorldSettings;
