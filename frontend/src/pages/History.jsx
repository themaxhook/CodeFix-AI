import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "../components/Sidebar";
import api from "../api/axios";

const LANGUAGES = ["All", "Python", "JavaScript", "TypeScript", "Go", "Java", "C++"];

export default function History() {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("All");
  const navigate = useNavigate();

  useEffect(() => {
    fetchReviews();
  }, [filter]);

  async function fetchReviews() {
    setLoading(true);
    try {
      const params = filter !== "All" ? `?language=${filter}` : "";
      const res = await api.get(`/reviews${params}`);
      setReviews(res.data);
    } catch (err) {
      console.error("failed to fetch reviews", err);
    } finally {
      setLoading(false);
    }
  }

  async function deleteReview(id) {
    try {
      await api.delete(`/reviews/${id}`);
      setReviews((prev) => prev.filter((r) => r.id !== id));
    } catch {
      alert("Failed to delete");
    }
  }

  function truncate(str, n = 80) {
    return str.length > n ? str.slice(0, n) + "..." : str;
  }

  function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString("en-IN", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
  }

  return (
    <div style={{ display: "flex", height: "100vh", overflow: "hidden" }}>
      <Sidebar />

      <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden", minHeight: 0 }}>

        {/* topbar */}
        <div className="topbar">
          <h1 className="page-title">Review history</h1>
          <div className="filter-row">
            {LANGUAGES.map((l) => (
              <button
                key={l}
                className={`filter-btn ${filter === l ? "active" : ""}`}
                onClick={() => setFilter(l)}
              >
                {l}
              </button>
            ))}
          </div>
        </div>

        {/* scrollable list */}
        <div style={{ flex: 1, overflowY: "auto", padding: "1rem 1.5rem", display: "flex", flexDirection: "column", gap: "0.75rem" }}>

          {loading && <p className="loading-text">Loading...</p>}

          {!loading && reviews.length === 0 && (
            <div className="empty-state">
              <p>No reviews yet. Go analyse some code!</p>
            </div>
          )}

          {reviews.map((review) => (
            <div key={review.id} className="history-card">
              <div className="card-top">
                <span className={`lang-badge ${review.language.toLowerCase().replace("+", "p")}`}>
                  {review.language}
                </span>
                <span className="card-date">{formatDate(review.created_at)}</span>
              </div>

              <pre className="card-snippet">{truncate(review.original_code)}</pre>
              <p className="card-bug">{truncate(review.bug_explanation, 100)}</p>

              <div className="card-actions">
                <button className="btn-view" onClick={() => navigate(`/reviews/${review.id}`)}>
                  View review
                </button>
                <button className="btn-delete" onClick={() => deleteReview(review.id)}>
                  Delete
                </button>
                {review.is_helpful !== null && (
                  <span className="helpful-badge">
                    {review.is_helpful ? "👍 Helpful" : "👎 Not helpful"}
                  </span>
                )}
              </div>
            </div>
          ))}

        3</div>
      </div>
    </div>
  );
}