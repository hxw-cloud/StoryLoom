import React, { useEffect, useState } from 'react';
import { plotService } from '../../services/plotService';
import type { PlotCard, PlotCardInput } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './PlotOutline.module.css';

const PlotOutline: React.FC = () => {
  const [plots, setPlots] = useState<PlotCard[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newPlot, setNewPlot] = useState<PlotCardInput>({
    title: '',
    description: '',
    conflict_intensity: 3,
    order_index: 1,
  });

  const fetchPlots = async () => {
    try {
      const data = await plotService.getPlots();
      setPlots(data);
      // Auto-increment order for next card
      if (data.length > 0) {
        const maxOrder = Math.max(...data.map(p => p.order_index));
        setNewPlot(prev => ({ ...prev, order_index: maxOrder + 1 }));
      }
    } catch (err) {
      setError('Failed to load plot beats.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlots();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await plotService.createPlot(newPlot);
      setIsModalOpen(false);
      setNewPlot({ title: '', description: '', conflict_intensity: 3, order_index: newPlot.order_index + 1 });
      fetchPlots();
    } catch (err) {
      alert('Failed to create plot beat');
    }
  };

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
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>+ New Beat</button>
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

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title="Add New Plot Beat"
      >
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label>Beat Title</label>
            <input 
              type="text" 
              required 
              value={newPlot.title}
              onChange={e => setNewPlot({...newPlot, title: e.target.value})}
              placeholder="e.g. The Call to Adventure"
            />
          </div>
          <div className={styles.formGroup}>
            <label>Sequence Order</label>
            <input 
              type="number" 
              required 
              value={newPlot.order_index}
              onChange={e => setNewPlot({...newPlot, order_index: parseInt(e.target.value)})}
            />
          </div>
          <div className={styles.formGroup}>
            <label>Conflict Intensity (1-5)</label>
            <input 
              type="range" 
              min="1" 
              max="5"
              value={newPlot.conflict_intensity}
              onChange={e => setNewPlot({...newPlot, conflict_intensity: parseInt(e.target.value)})}
            />
            <span style={{textAlign: 'center', fontWeight: 'bold'}}>{newPlot.conflict_intensity}</span>
          </div>
          <div className={styles.formGroup}>
            <label>Description</label>
            <textarea 
              value={newPlot.description}
              onChange={e => setNewPlot({...newPlot, description: e.target.value})}
              placeholder="What happens in this beat?"
            />
          </div>
          <button type="submit" className={styles.submitButton}>Create Beat</button>
        </form>
      </Modal>
    </div>
  );
};

export default PlotOutline;
