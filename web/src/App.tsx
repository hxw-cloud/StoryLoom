import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Sidebar from './components/Sidebar';
import WorldSettings from './features/world/WorldSettings';
import Characters from './features/character/Characters';
import PlotOutline from './features/plot/PlotOutline';
import SceneEditor from './features/scene/SceneEditor';
import Timeline from './features/timeline/Timeline';
import styles from './App.module.css';

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
