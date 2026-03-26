import React, { useEffect, useState } from 'react';
import { timelineService } from '../../services/timelineService';
import type { TimelineEvent } from '../../services/types';
import styles from './Timeline.module.css';

const Timeline: React.FC = () => {
  const [events, setEvents] = useState<TimelineEvent[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const data = await timelineService.getEvents();
        setEvents(data);
      } catch (err) {
        console.error('Failed to load timeline events', err);
      } finally {
        setLoading(false);
      }
    };
    fetchEvents();
  }, []);

  if (loading) return <div className={styles.loading}>Chronological alignment in progress...</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>Master Timeline</h2>
        <button className={styles.addButton}>+ Add Historical Event</button>
      </header>

      <div className={styles.timeline}>
        {events.length === 0 ? (
          <div className={styles.empty}>The history of your world is a blank slate. Add your first event!</div>
        ) : (
          events.map((event, index) => (
            <div key={event.id} className={styles.timelineItem}>
              <div className={styles.connector}>
                <div className={styles.dot} />
                {index < events.length - 1 && <div className={styles.line} />}
              </div>
              <div className={styles.content}>
                <div className={styles.orderBadge}>Year/Order {event.chronological_order}</div>
                <h3>{event.title}</h3>
                <p>{event.description}</p>
                {event.scene_id && (
                  <div className={styles.sceneLink}>馃敆 Linked to Scene: {event.scene_id}</div>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Timeline;
