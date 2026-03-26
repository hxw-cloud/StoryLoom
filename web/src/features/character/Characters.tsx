import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { characterService } from '../../services/characterService';
import type { Character, CharacterInput, Relationship, CharacterArc } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './Characters.module.css';

const Characters: React.FC = () => {
  const { t } = useTranslation();
  const [characters, setCharacters] = useState<Character[]>([]);
  const [relationships, setRelationships] = useState<Relationship[]>([]);
  const [activeTab, setActiveTab] = useState<'list' | 'network' | 'arcs'>('list');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Search & Filter
  const [searchTerm, setSearchTerm] = useState('');
  const [campFilter, setCategoryFilter] = useState('');

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newChar, setNewChar] = useState<CharacterInput>({
    name: '',
    role: 'Protagonist',
    camp: '',
    age: 20,
    gender: 'Male',
    appearance: '',
    background: '',
    pov_type: 'Third Person Limited',
    want: '',
    need: '',
    persona_template: '',
  });

  const fetchData = async () => {
    try {
      const [charData, relData] = await Promise.all([
        characterService.getCharacters(),
        characterService.getRelationships(),
      ]);
      setCharacters(charData);
      setRelationships(relData);
    } catch (err) {
      setError(t('character.error'));
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
      await characterService.createCharacter(newChar);
      setIsModalOpen(false);
      setNewChar({
        name: '', role: 'Protagonist', camp: '', age: 20, gender: 'Male',
        appearance: '', background: '', pov_type: 'Third Person Limited',
        want: '', need: '', persona_template: ''
      });
      fetchData();
    } catch (err) {
      alert(t('character.addError'));
    }
  };

  const filteredCharacters = characters.filter(c => {
    const matchesSearch = c.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                         c.background?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCamp = campFilter === '' || c.camp === campFilter;
    return matchesSearch && matchesCamp;
  });

  if (loading) return <div className={styles.loading}>{t('common.loading')}</div>;
  if (error) return <div className={styles.error}>{error}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h2>{t('character.title')}</h2>
          <div className={styles.tabs}>
            <button className={`${styles.tab} ${activeTab === 'list' ? styles.activeTab : ''}`} onClick={() => setActiveTab('list')}>{t('character.tabs.list')}</button>
            <button className={`${styles.tab} ${activeTab === 'network' ? styles.activeTab : ''}`} onClick={() => setActiveTab('network')}>{t('character.tabs.network')}</button>
            <button className={`${styles.tab} ${activeTab === 'arcs' ? styles.activeTab : ''}`} onClick={() => setActiveTab('arcs')}>{t('character.tabs.arcs')}</button>
          </div>
        </div>
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>{t('character.newCharacter')}</button>
      </header>

      {activeTab === 'list' && (
        <>
          <div className={styles.searchBar}>
            <input type="text" placeholder={t('common.search')} value={searchTerm} onChange={e => setSearchTerm(e.target.value)} />
          </div>
          <div className={styles.grid}>
            {filteredCharacters.length === 0 ? (
              <div className={styles.empty}>{t('character.empty')}</div>
            ) : (
              filteredCharacters.map(char => (
                <div key={char.id} className={styles.card}>
                  <div className={styles.cardHeader}>
                    <h3>{char.name}</h3>
                    <span className={styles.roleBadge}>{char.role}</span>
                  </div>
                  {char.camp && <div className={styles.campLabel}>{char.camp}</div>}
                  <div className={styles.motivationBox}>
                    <div><strong>Want:</strong> {char.want || '???'}</div>
                    <div><strong>Need:</strong> {char.need || '???'}</div>
                  </div>
                  <p className={styles.bio}>{char.background?.substring(0, 100)}...</p>
                </div>
              ))
            )}
          </div>
        </>
      )}

      {activeTab === 'network' && (
        <div className={styles.networkView}>
          <div className={styles.relList}>
            {relationships.map((rel, i) => (
              <div key={i} className={styles.relItem}>
                <strong>{rel.source_id}</strong> is <span>{rel.type}</span> to <strong>{rel.target_id}</strong>
              </div>
            ))}
            {relationships.length === 0 && <p className={styles.empty}>No relationships defined. Connect your characters!</p>}
          </div>
        </div>
      )}

      {activeTab === 'arcs' && (
        <div className={styles.arcsView}>
          <p className={styles.empty}>Select a character to view their growth arc across the story beats.</p>
        </div>
      )}

      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title={t('character.newCharacter')}>
        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGrid}>
            <div className={styles.formGroup}><label>{t('character.fields.name')}</label><input type="text" required value={newChar.name} onChange={e => setNewChar({...newChar, name: e.target.value})}/></div>
            <div className={styles.formGroup}><label>{t('character.fields.role')}</label>
              <select value={newChar.role} onChange={e => setNewChar({...newChar, role: e.target.value})}>
                <option value="Protagonist">Protagonist</option>
                <option value="Antagonist">Antagonist</option>
                <option value="Supporting">Supporting</option>
                <option value="Mentor">Mentor</option>
              </select>
            </div>
            <div className={styles.formGroup}><label>{t('character.fields.camp')}</label><input type="text" value={newChar.camp} onChange={e => setNewChar({...newChar, camp: e.target.value})} placeholder="e.g. Rebel Alliance"/></div>
            <div className={styles.formGroup}><label>{t('character.fields.age')}</label><input type="number" value={newChar.age} onChange={e => setNewChar({...newChar, age: parseInt(e.target.value)})}/></div>
          </div>
          <div className={styles.formGroup}><label>{t('character.fields.want')}</label><input type="text" value={newChar.want} onChange={e => setNewChar({...newChar, want: e.target.value})} placeholder="External Goal"/></div>
          <div className={styles.formGroup}><label>{t('character.fields.need')}</label><input type="text" value={newChar.need} onChange={e => setNewChar({...newChar, need: e.target.value})} placeholder="Internal Need / Growth"/></div>
          <div className={styles.formGroup}><label>{t('character.fields.background')}</label><textarea value={newChar.background} onChange={e => setNewChar({...newChar, background: e.target.value})}/></div>
          <button type="submit" className={styles.submitButton}>{t('common.save')}</button>
        </form>
      </Modal>
    </div>
  );
};

export default Characters;
