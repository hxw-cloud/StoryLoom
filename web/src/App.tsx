import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Sidebar from './components/Sidebar';
import WorldSettings from './features/world/WorldSettings';
import styles from './App.module.css';

// Placeholder components
const Characters = () => <div className={styles.content}><h2>Characters</h2><p>Track your characters and their POVs.</p></div>;
const PlotOutline = () => <div className={styles.content}><h2>Plot Outline</h2><p>Structure your story beats and pacing.</p></div>;
const SceneEditor = () => <div className={styles.content}><h2>Scene Editor</h2><p>Write your scenes and get logic audits.</p></div>;
const Timeline = () => <div className={styles.content}><h2>Timeline</h2><p>Visualize the chronological flow of events.</p></div>;

const App: React.FC = () => {
  return (
    <Router>
      <div className={styles.app}>
        <Sidebar />
        <main className={styles.main}>
          <Routes>
            <Route path="/world" element={<WorldSettings />} />
            <Route path="/character" element={<Characters />} />
            <Route path="/plot" element={<PlotOutline />} />
            <Route path="/scene" element={<SceneEditor />} />
            <Route path="/timeline" element={<Timeline />} />
            <Route path="/" element={<Navigate to="/world" replace />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

export default App;
