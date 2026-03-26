import React from 'react';
import { NavLink } from 'react-router-dom';
import styles from './Sidebar.module.css';

const Sidebar: React.FC = () => {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>
        <h1>StoryLoom</h1>
      </div>
      <nav className={styles.nav}>
        <NavLink to="/world" className={({ isActive }) => isActive ? styles.active : ''}>
          World Building
        </NavLink>
        <NavLink to="/character" className={({ isActive }) => isActive ? styles.active : ''}>
          Characters
        </NavLink>
        <NavLink to="/plot" className={({ isActive }) => isActive ? styles.active : ''}>
          Plot Outline
        </NavLink>
        <NavLink to="/scene" className={({ isActive }) => isActive ? styles.active : ''}>
          Scene Editor
        </NavLink>
        <NavLink to="/timeline" className={({ isActive }) => isActive ? styles.active : ''}>
          Timeline
        </NavLink>
      </nav>
    </aside>
  );
};

export default Sidebar;
