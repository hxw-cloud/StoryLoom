import React, { useEffect, useState } from 'react';
import { sceneService } from '../../services/sceneService';
import type { Scene, AuditResult } from '../../services/types';
import styles from './SceneEditor.module.css';

const SceneEditor: React.FC = () => {
  const [scenes, setScenes] = useState<Scene[]>([]);
  const [selectedScene, setSelectedScene] = useState<Scene | null>(null);
  const [auditResult, setAuditResult] = useState<AuditResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [auditing, setAuditing] = useState(false);

  useEffect(() => {
    const fetchScenes = async () => {
      try {
        const data = await sceneService.getScenes();
        setScenes(data);
        if (data.length > 0) setSelectedScene(data[0]);
      } catch (err) {
        console.error('Failed to load scenes', err);
      } finally {
        setLoading(false);
      }
    };
    fetchScenes();
  }, []);

  const handleAudit = async () => {
    if (!selectedScene) return;
    setAuditing(true);
    try {
      const result = await sceneService.auditScene(selectedScene.id);
      setAuditResult(result);
    } catch (err) {
      console.error('Audit failed', err);
    } finally {
      setAuditing(false);
    }
  };

  if (loading) return <div className={styles.loading}>Opening manuscript...</div>;

  return (
    <div className={styles.container}>
      <div className={styles.sceneList}>
        <h3>Scenes</h3>
        {scenes.map(s => (
          <button 
            key={s.id} 
            className={`${styles.sceneItem} ${selectedScene?.id === s.id ? styles.active : ''}`}
            onClick={() => {
              setSelectedScene(s);
              setAuditResult(null);
            }}
          >
            {s.title}
          </button>
        ))}
      </div>

      <div className={styles.editorMain}>
        {selectedScene ? (
          <>
            <header className={styles.editorHeader}>
              <input 
                type="text" 
                value={selectedScene.title} 
                className={styles.titleInput} 
                readOnly
              />
              <button 
                className={styles.auditButton}
                onClick={handleAudit}
                disabled={auditing}
              >
                {auditing ? 'Analyzing...' : 'Run Digital Editor'}
              </button>
            </header>
            <div className={styles.writingArea}>
              <div className={styles.metaRow}>
                <span><strong>POV:</strong> {selectedScene.pov_character_id || 'Not Set'}</span>
                <span><strong>Goal:</strong> {selectedScene.goal || 'No goal defined'}</span>
              </div>
              <textarea 
                className={styles.textarea} 
                placeholder="Once upon a time..."
                defaultValue={`This is where the story for "${selectedScene.title}" would be written. The Digital Editor will analyze the logic rules defined in the World Building module against this content.`}
              />
            </div>
          </>
        ) : (
          <div className={styles.empty}>Select a scene to start writing.</div>
        )}
      </div>

      <aside className={styles.auditSidebar}>
        <h3>Digital Editor</h3>
        {auditResult ? (
          <div className={styles.auditContent}>
            {auditResult.is_valid ? (
              <div className={styles.validStatus}>鉁 No logic violations found.</div>
            ) : (
              <div className={styles.invalidStatus}>鈿 Issues Found:</div>
            )}
            <ul className={styles.issueList}>
              {auditResult.issues.map((issue, i) => (
                <li key={i} className={styles.issueItem}>{issue}</li>
              ))}
            </ul>
          </div>
        ) : (
          <p className={styles.auditHint}>Run the analysis to check for narrative consistency.</p>
        )}
      </aside>
    </div>
  );
};

export default SceneEditor;
