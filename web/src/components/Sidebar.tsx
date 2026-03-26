import React from 'react';
import { NavLink } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import styles from './Sidebar.module.css';

const Sidebar: React.FC = () => {
  const { t, i18n } = useTranslation();

  const toggleLanguage = () => {
    const newLang = i18n.language === 'zh' ? 'en' : 'zh';
    i18n.changeLanguage(newLang);
  };

  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>
        <h1>StoryLoom</h1>
      </div>
      <nav className={styles.nav}>
        <NavLink to="/world" className={({ isActive }) => isActive ? styles.active : ''}>
          {t('nav.world')}
        </NavLink>
        <NavLink to="/character" className={({ isActive }) => isActive ? styles.active : ''}>
          {t('nav.character')}
        </NavLink>
        <NavLink to="/plot" className={({ isActive }) => isActive ? styles.active : ''}>
          {t('nav.plot')}
        </NavLink>
        <NavLink to="/scene" className={({ isActive }) => isActive ? styles.active : ''}>
          {t('nav.scene')}
        </NavLink>
        <NavLink to="/timeline" className={({ isActive }) => isActive ? styles.active : ''}>
          {t('nav.timeline')}
        </NavLink>
      </nav>
      <div className={styles.footer}>
        <button onClick={toggleLanguage} className={styles.langButton}>
          {i18n.language === 'zh' ? 'English' : '中文'}
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
