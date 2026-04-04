import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import Sidebar from "../components/Sidebar";
import api from "../api/axios";

export default function ReviewDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [review, setReview] = useState(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    fetchReview();
  }, [id]);

  async function fetchReview() {
    try {
      const res = await api.get(`/reviews/${id}`);
      setReview(res.data);
    } catch {
      navigate("/history");
    } finally {
      setLoading(false);
    }
  }

  async function rateReview(isHelpful) {
    try {
      await api.patch(`/reviews/${id}/rate`, { is_helpful: isHelpful });
      setReview((prev) => ({ ...prev, is_helpful: isHelpful }));
    } catch {
      // silently fail
    }
  }

  function copyCode() {
    navigator.clipboard.writeText(review.fixed_code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString("en-IN", {
      day: "numeric", month: "long", year: "numeric",
      hour: "2-digit", minute: "2-digit",
    });
  }

  if (loading) {
    return (
      <div style={{ display: "flex", height: "100vh" }}>
        <Sidebar />
        <div style={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center" }}>
          <p style={{ color: "var(--text3)" }}>Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div style={{ display: "flex", height: "100vh", overflow: "hidden" }}>
      <Sidebar />

      <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden", minHeight: 0 }}>

        {/* topbar */}
        <div className="topbar">
          <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
            <button className="back-btn" onClick={() => navigate("/history")}>← Back</button>
            <h1 className="page-title">Review detail</h1>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
            <span className={`lang-badge ${review.language.toLowerCase().replace("+", "p")}`}>
              {review.language}
            </span>
            <span className="card-date">{formatDate(review.created_at)}</span>
          </div>
        </div>

        {/* scrollable content */}
        <div style={{ flex: 1, overflowY: "auto", padding: "1.25rem 1.5rem", display: "flex", flexDirection: "column", gap: "1.25rem", maxWidth: "900px", width: "100%" }}>

          <div className="detail-section">
            <div className="detail-section-title">Original code</div>
            <pre className="code-block">{review.original_code}</pre>
          </div>

          <div className="detail-section">
            <div className="detail-section-title red">Bug found</div>
            <p className="detail-text">{review.bug_explanation}</p>
          </div>

          <div className="detail-section">
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <div className="detail-section-title green">Fixed code</div>
              <button className="copy-btn" style={{ margin: "0.75rem 1rem" }} onClick={copyCode}>
                {copied ? "Copied!" : "Copy"}
              </button>
            </div>
            <pre className="code-block">{review.fixed_code}</pre>
          </div>

          <div className="detail-section">
            <div className="detail-section-title purple">Suggestions</div>
            <p className="detail-text">{review.suggestions}</p>
          </div>

          <div className="rating-row">
            <span className="rating-label">Was this helpful?</span>
            <button
              className={`thumb-btn ${review.is_helpful === true ? "active" : ""}`}
              onClick={() => rateReview(true)}
            >👍</button>
            <button
              className={`thumb-btn ${review.is_helpful === false ? "active" : ""}`}
              onClick={() => rateReview(false)}
            >👎</button>
          </div>

        </div>
      </div>
    </div>
  );
}