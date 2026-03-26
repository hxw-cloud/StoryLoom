import React, { useEffect, useState } from 'react';
import { worldService } from '../../services/worldService';
import type { WorldSetting, WorldSettingInput } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './WorldSettings.module.css';

const WorldSettings: React.FC = () => {
  const [settings, setSettings] = useState<WorldSetting[]>([]);
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

  const fetchSettings = async () => {
    try {
      const data = await worldService.getSettings();
      setSettings(data);
    } catch (err) {
      setError('Failed to load world settings.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSettings();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await worldService.createSetting(newSetting);
      setIsModalOpen(false);
      setNewSetting({ category: 'Geography', name: '', description: '', logic_rules: '' });
      fetchSettings(); // Refresh list
    } catch (err) {
      alert('Failed to create setting');
    }
  };

  if (loading) return <div className={styles.loading}>Loading settings...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>World Building</h2>
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>+ New Setting</button>
      </header>
      
      <div className={styles.grid}>
        {settings.length === 0 ? (
          <div className={styles.empty}>No settings defined yet. Start building your world!</div>
        ) : (
          settings.map(setting => (
            <div key={setting.id} className={styles.card}>
              <div className={styles.categoryBadge}>{setting.category}</div>
              <h3>{setting.name}</h3>
              <p className={styles.description}>{setting.description}</p>
              {setting.logic_rules && (
                <div className={styles.logicRules}>
                  <strong>Logic Rules:</strong>
                  <p>{setting.logic_rules}</p>
                </div>
              )}
            </div>
          ))
        )}
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title="Add New World Setting"
      >
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label>Category</label>
            <select 
              value={newSetting.category} 
              onChange={e => setNewSetting({...newSetting, category: e.target.value})}
            >
              <option value="Geography">Geography</option>
              <option value="Magic System">Magic System</option>
              <option value="Technology">Technology</option>
              <option value="Race">Race</option>
              <option value="Culture">Culture</option>
            </select>
          </div>
          <div className={styles.formGroup}>
            <label>Name</label>
            <input 
              type="text" 
              required 
              value={newSetting.name}
              onChange={e => setNewSetting({...newSetting, name: e.target.value})}
              placeholder="e.g. The Law of Gravity"
            />
          </div>
          <div className={styles.formGroup}>
            <label>Description</label>
            <textarea 
              value={newSetting.description}
              onChange={e => setNewSetting({...newSetting, description: e.target.value})}
              placeholder="Describe the setting..."
            />
          </div>
          <div className={styles.formGroup}>
            <label>Logic Rules</label>
            <textarea 
              value={newSetting.logic_rules}
              onChange={e => setNewSetting({...newSetting, logic_rules: e.target.value})}
              placeholder="Machine-readable rules or strict constraints..."
            />
          </div>
          <button type="submit" className={styles.submitButton}>Save Setting</button>
        </form>
      </Modal>
    </div>
  );
};

export default WorldSettings;
