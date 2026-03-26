import React, { useEffect, useState } from 'react';
import { worldService } from '../../services/worldService';
import type { WorldSetting } from '../../services/types';
import styles from './WorldSettings.module.css';

const WorldSettings: React.FC = () => {
  const [settings, setSettings] = useState<WorldSetting[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
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

    fetchSettings();
  }, []);

  if (loading) return <div className={styles.loading}>Loading settings...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>World Building</h2>
        <button className={styles.addButton}>+ New Setting</button>
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
    </div>
  );
};

export default WorldSettings;
