import React, { useEffect, useState } from 'react';
import { characterService } from '../../services/characterService';
import type { Character } from '../../services/types';
import styles from './Characters.module.css';

const Characters: React.FC = () => {
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCharacters = async () => {
      try {
        const data = await characterService.getCharacters();
        setCharacters(data);
      } catch (err) {
        setError('Failed to load characters.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCharacters();
  }, []);

  if (loading) return <div className={styles.loading}>Loading characters...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>Characters</h2>
        <button className={styles.addButton}>+ New Character</button>
      </header>
      
      <div className={styles.list}>
        {characters.length === 0 ? (
          <div className={styles.empty}>No characters tracked yet. Create your protagonist!</div>
        ) : (
          characters.map(char => (
            <div key={char.id} className={styles.card}>
              <div className={styles.cardHeader}>
                <h3>{char.name}</h3>
                <span className={styles.roleBadge}>{char.role}</span>
              </div>
              {char.pov_type && (
                <div className={styles.povInfo}>
                  <strong>POV Type:</strong> {char.pov_type}
                </div>
              )}
              <div className={styles.meta}>
                Added: {new Date(char.created_at).toLocaleDateString()}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Characters;
