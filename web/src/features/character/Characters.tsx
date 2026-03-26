import React, { useEffect, useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import ForceGraph2D from 'react-force-graph-2d';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { characterService } from '../../services/characterService';
import type { Character, CharacterInput, Relationship } from '../../services/types';
import Modal from '../../components/Modal';
import styles from './Characters.module.css';

const Characters: React.FC = () => {
  const { t } = useTranslation();
  const [characters, setCharacters] = useState<Character[]>([]);
  const [relationships, setRelationships] = useState<Relationship[]>([]);
  const [activeTab, setActiveTab] = useState<'list' | 'network' | 'arcs' | 'blindTest'>('list');
  const [loading, setLoading] = useState(true);
  
  // Search & Filter
  const [searchTerm, setSearchTerm] = useState('');

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

  // Blind Test State
  const [revealName, setRevealName] = useState(false);
  const [currentQuoteIndex, setCurrentQuoteIndex] = useState(0);
  const mockQuotes = [
    { text: "I don't care about the rules. I care about results.", character: "Protagonist" },
    { text: "The stars align, but only for those who know how to read them.", character: "Mentor" },
    { text: "Everything has a price. You just haven't seen the bill yet.", character: "Antagonist" }
  ];

  const fetchData = async () => {
    try {
      const [charData, relData] = await Promise.all([
        characterService.getCharacters(),
        characterService.getRelationships(),
      ]);
      setCharacters(charData);
      setRelationships(relData);
    } catch (err) {
      console.error(t('character.error'), err);
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

  // REQ-2.4: Relationship Map Data Transformation
  const graphData = useMemo(() => {
    const nodes = characters.map(c => ({ id: c.id, name: c.name, role: c.role }));
    const links = relationships.map(r => ({ source: r.source_id, target: r.target_id, label: r.type }));
    return { nodes, links };
  }, [characters, relationships]);

  // REQ-2.3: Mock Growth Arc Data
  const arcData = [
    { name: 'Beat 1', growth: 10 },
    { name: 'Beat 2', growth: 25 },
    { name: 'Beat 3', growth: 20 },
    { name: 'Beat 4', growth: 45 },
    { name: 'Beat 5', growth: 80 },
  ];

  if (loading) return <div className={styles.loading}>{t('common.loading')}</div>;

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h2>{t('character.title')}</h2>
          <div className={styles.tabs}>
            <button className={`${styles.tab} ${activeTab === 'list' ? styles.activeTab : ''}`} onClick={() => setActiveTab('list')}>{t('character.tabs.list')}</button>
            <button className={`${styles.tab} ${activeTab === 'network' ? styles.activeTab : ''}`} onClick={() => setActiveTab('network')}>{t('character.tabs.network')}</button>
            <button className={`${styles.tab} ${activeTab === 'arcs' ? styles.activeTab : ''}`} onClick={() => setActiveTab('arcs')}>{t('character.tabs.arcs')}</button>
            <button className={`${styles.tab} ${activeTab === 'blindTest' ? styles.activeTab : ''}`} onClick={() => setActiveTab('blindTest')}>{t('character.tabs.blindTest')}</button>
          </div>
        </div>
        <button className={styles.addButton} onClick={() => setIsModalOpen(true)}>{t('character.newCharacter')}</button>
      </header>

      {activeTab === 'list' && (
        <div className={styles.listContent}>
          <div className={styles.searchBar}>
            <input type="text" placeholder={t('common.search')} value={searchTerm} onChange={e => setSearchTerm(e.target.value)} />
          </div>
          <div className={styles.grid}>
            {characters.map(char => (
              <div key={char.id} className={styles.card}>
                <div className={styles.cardHeader}>
                  <h3>{char.name}</h3>
                  <span className={styles.roleBadge}>{char.role}</span>
                </div>
                <div className={styles.motivationRow}>
                  <span><strong>Want:</strong> {char.want || 'None'}</span>
                  <span><strong>Need:</strong> {char.need || 'None'}</span>
                </div>
                <p className={styles.bio}>{char.background || 'No background story.'}</p>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'network' && (
        <div className={styles.graphContainer}>
          <p className={styles.hint}>{t('character.network.hint')}</p>
          <div className={styles.canvasWrapper}>
            <ForceGraph2D
              graphData={graphData}
              nodeAutoColorBy="role"
              nodeCanvasObject={(node: any, ctx: CanvasRenderingContext2D, globalScale: number) => {
                const label = node.name;
                const fontSize = 12/globalScale;
                ctx.font = `${fontSize}px Sans-Serif`;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillStyle = node.color;
                ctx.beginPath(); ctx.arc(node.x, node.y, 5, 0, 2 * Math.PI, false); ctx.fill();
                ctx.fillStyle = 'black';
                ctx.fillText(label, node.x, node.y + 10);
              }}
              linkDirectionalArrowLength={3.5}
              linkDirectionalArrowRelPos={1}
              linkLabel="label"
            />
          </div>
        </div>
      )}

      {activeTab === 'arcs' && (
        <div className={styles.chartContainer}>
          <h3>Character Growth Arc</h3>
          <div className={styles.chartWrapper}>
            <ResponsiveContainer width="100%" height={400}>
              <LineChart data={arcData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="growth" stroke="var(--color-primary)" strokeWidth={3} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}

      {activeTab === 'blindTest' && (
        <div className={styles.blindTestContainer}>
          <h3>{t('character.blindTest.title')}</h3>
          <div className={styles.testCard}>
            <p className={styles.quote}>"{mockQuotes[currentQuoteIndex].text}"</p>
            {revealName && <div className={styles.revealName}>{mockQuotes[currentQuoteIndex].character}</div>}
            <div className={styles.testActions}>
              <button onClick={() => setRevealName(true)}>{t('character.blindTest.reveal')}</button>
              <button onClick={() => { setCurrentQuoteIndex((currentQuoteIndex + 1) % mockQuotes.length); setRevealName(false); }}>{t('character.blindTest.next')}</button>
            </div>
          </div>
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
          </div>
          <div className={styles.formGroup}><label>{t('character.fields.want')}</label><input type="text" value={newChar.want} onChange={e => setNewChar({...newChar, want: e.target.value})}/></div>
          <div className={styles.formGroup}><label>{t('character.fields.need')}</label><input type="text" value={newChar.need} onChange={e => setNewChar({...newChar, need: e.target.value})}/></div>
          <button type="submit" className={styles.submitButton}>{t('common.save')}</button>
        </form>
      </Modal>
    </div>
  );
};

export default Characters;
