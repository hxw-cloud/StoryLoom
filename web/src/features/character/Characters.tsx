import React, { useEffect, useState } from 'react';
import { characterService } from '../../services/characterService';
import type { Character, CharacterInput } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './Characters.module.css';

const Characters: React.FC = () => {
  const [characters, setCharacters] = useState<Character[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newChar, setNewChar] = useState<CharacterInput>({
    name: '',
    role: 'Protagonist',
    pov_type: 'Third Person Limited',
  });

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

  useEffect(() => {
    fetchCharacters();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await characterService.createCharacter(newChar);
      setIsModalOpen(false);
      setNewChar({ name: '', role: 'Protagonist', pov_type: 'Third Person Limited' });
      fetchCharacters();
    } catch (err) {
      alert('Failed to create character');
    }
  };

  if (loading) return <div className={styles.loading}>Loading characters...</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2>Characters</h2>
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>+ New Character</button>
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

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title="Add New Character"
      >
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label>Full Name</label>
            <input 
              type="text" 
              required 
              value={newChar.name}
              onChange={e => setNewChar({...newChar, name: e.target.value})}
              placeholder="e.g. Elara Vance"
            />
          </div>
          <div className={styles.formGroup}>
            <label>Structural Role</label>
            <select 
              value={newChar.role} 
              onChange={e => setNewChar({...newChar, role: e.target.value})}
            >
              <option value="Protagonist">Protagonist</option>
              <option value="Antagonist">Antagonist</option>
              <option value="Supporting">Supporting</option>
              <option value="Mentor">Mentor</option>
              <option value="Love Interest">Love Interest</option>
            </select>
          </div>
          <div className={styles.formGroup}>
            <label>Default POV Type</label>
            <select 
              value={newChar.pov_type} 
              onChange={e => setNewChar({...newChar, pov_type: e.target.value})}
            >
              <option value="First Person">First Person</option>
              <option value="Third Person Limited">Third Person Limited</option>
              <option value="Third Person Omniscient">Third Person Omniscient</option>
              <option value="Third Person Objective">Third Person Objective</option>
            </select>
          </div>
          <button type="submit" className={styles.submitButton}>Create Character</button>
        </form>
      </Modal>
    </div>
  );
};

export default Characters;
