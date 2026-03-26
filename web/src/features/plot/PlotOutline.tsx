import React, { useEffect, useState } from 'react';
import { plotService } from '../../services/plotService';
import type { PlotCard } from '../../services/types';
import styles from './PlotOutline.module.css';

const PlotOutline: React.FC = () => {
  const [plots, setPlots] = useState<PlotCard[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlots = async () => {
      try {
        const data = await plotService.getPlots();
        setPlots(data);
      } catch (err) {
        setError('Failed to load plot beats.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchPlots();
  }, []);

  const getIntensityColor = (intensity: number) => {
    if (intensity >= 4) return 'var(--color-accent-crimson)';
    if (intensity >= 3) return 'var(--color-accent-amber)';
    return 'var(--color-accent-teal)';
  };

  if (loading) return <div className={styles.loading}>Loading story beats...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>Plot Outline</h2>
        <button className={styles.addButton}>+ New Beat</button>
      </header>
      
      <div className={styles.board}>
        {plots.length === 0 ? (
          <div className={styles.empty}>Your story board is empty. Add your first inciting incident!</div>
        ) : (
          plots.map(plot => (
            <div key={plot.id} className={styles.plotCard}>
              <div className={styles.cardIndex}>#{plot.order_index}</div>
              <h3>{plot.title}</h3>
              <p className={styles.description}>{plot.description}</p>
              <div className={styles.intensityTracker}>
                <span>Intensity</span>
                <div className={styles.intensityBar}>
                  <div 
                    className={styles.intensityFill} 
                    style={{ 
                      width: `${(plot.conflict_intensity / 5) * 100}%`,
                      backgroundColor: getIntensityColor(plot.conflict_intensity)
                    }} 
                  />
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default PlotOutline;
