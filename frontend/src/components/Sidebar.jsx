import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../AuthContext";

export default function Sidebar() {
  const { email, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  const initials = email ? email[0].toUpperCase() : "?";

  return (
    <div className="sidebar">
      <div className="sidebar-logo">
        CodeFix <span>AI</span>
      </div>

      <nav className="sidebar-nav">
        <NavLink
          to="/dashboard"
          className={({ isActive }) => isActive ? "nav-item active" : "nav-item"}
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <rect x="2" y="2" width="12" height="12" rx="2" stroke="currentColor" strokeWidth="1.2"/>
            <path d="M5 8h6M5 5.5h3M5 10.5h4" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round"/>
          </svg>
          Analyser
        </NavLink>

        <NavLink
          to="/history"
          className={({ isActive }) => isActive ? "nav-item active" : "nav-item"}
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M3 4h10M3 8h7M3 12h5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round"/>
          </svg>
          History
        </NavLink>
      </nav>

      <div className="sidebar-user">
        <div className="avatar">{initials}</div>
        <div className="user-info">
          <p className="user-email">{email}</p>
          <button className="logout-btn" onClick={handleLogout}>Logout</button>
        </div>
      </div>
    </div>
  );
}