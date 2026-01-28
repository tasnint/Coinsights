import React from 'react';
import '../styles/Header.css';

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="header-logo">
        <img src="/logo.png" alt="Coinsights" className="logo-image" />
      </div>
    </header>
  );
};

export default Header;
