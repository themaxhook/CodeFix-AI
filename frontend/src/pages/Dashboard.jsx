import { useState } from "react";
import Editor from "@monaco-editor/react";
import Sidebar from "../components/Sidebar";
import api from "../api/axios";
import ReactMarkdown from "react-markdown";

const LANGUAGES = ["Python", "JavaScript", "TypeScript", "Go", "Java", "C++", "Rust"];

export default function Dashboard() {
  const [code, setCode] = useState("");
  const [language, setLanguage] = useState("Python");
  const [review, setReview] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  async function handleAnalyse() {
    if (!code.trim()) {
      setError("Please paste some code first");
      return;
    }
    setError("");
    setLoading(true);
    setReview(null);

    try {
      const res = await api.post("/reviews", { code, language });
      setReview(res.data);
    } catch (err) {
      const data = err.response?.data;
      const detail = typeof data === "string" ? data : data?.message || err.message;
      setError(detail?.trim() ? detail : "Failed to analyse code. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  function copyFixedCode() {
    navigator.clipboard.writeText(review.fixed_code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  async function rateReview(isHelpful) {
    try {
      await api.patch(`/reviews/${review.id}/rate`, { is_helpful: isHelpful });
      setReview((prev) => ({ ...prev, is_helpful: isHelpful }));
    } catch {
      // not critical
    }
  }

  const editorLanguage = language.toLowerCase().replace("c++", "cpp");

  return (
    <div style={{ display: "flex", height: "100vh", overflow: "hidden" }}>
      <Sidebar />

      <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "hidden" }}>
        {/* topbar */}
        <div className="topbar">
          <h1 className="page-title">Code Analyser</h1>
          <select
            className="lang-select"
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
          >
            {LANGUAGES.map((l) => (
              <option key={l} value={l}>{l}</option>
            ))}
          </select>
        </div>

        {/* two panels */}
        <div style={{ flex: 1, display: "grid", gridTemplateColumns: "1fr 1fr", gap: "1rem", padding: "1rem 1.5rem", overflow: "hidden", minHeight: 0 }}>

          {/* left — code input */}
          <div style={{ display: "flex", flexDirection: "column", background: "#181825", borderRadius: "8px", border: "1px solid #313244", overflow: "hidden", minHeight: 0 }}>
            <div className="panel-header">
              <span className="panel-label">Your code</span>
            </div>
            <div style={{ flex: 1, minHeight: 0 }}>
              <Editor
                height="100%"
                language={editorLanguage}
                value={code}
                onChange={(val) => setCode(val || "")}
                theme="vs-dark"
                options={{
                  fontSize: 13,
                  minimap: { enabled: false },
                  scrollBeyondLastLine: false,
                  padding: { top: 12 },
                }}
              />
            </div>
            <div className="panel-footer">
              {error && <p className="error-msg">{error}</p>}
              <button className="btn-primary" onClick={handleAnalyse} disabled={loading}>
                {loading ? "Analysing..." : "Analyse code"}
              </button>
            </div>
          </div>

          {/* right — ai result */}
          <div style={{ display: "flex", flexDirection: "column", background: "#181825", borderRadius: "8px", border: "1px solid #313244", overflow: "hidden", minHeight: 0 }}>
            <div className="panel-header">
              <span className="panel-label">AI review</span>
            </div>

            <div style={{ flex: 1, overflowY: "auto", padding: "1rem", display: "flex", flexDirection: "column", gap: "0.875rem", minHeight: 0 }}>
              {!review && !loading && (
                <div className="empty-state">
                  <p>Paste your code and click Analyse</p>
                </div>
              )}

              {loading && (
                <div className="empty-state">
                  <div className="spinner"></div>
                  <p>Analysing your code...</p>
                </div>
              )}

              {review && (
                <>
                  <div className="review-section bug">
                    <div className="review-section-header">
                      <span className="dot red"></span>
                      Bug found
                    </div>
                    <div className="review-text">
                      <ReactMarkdown>{review.bug_explanation}</ReactMarkdown>
                    </div>
                  </div>

                  <div className="review-section fix">
                    <div className="review-section-header">
                      <span className="dot green"></span>
                      Fixed code
                      <button className="copy-btn" onClick={copyFixedCode}>
                        {copied ? "Copied!" : "Copy"}
                      </button>
                    </div>
                    <pre className="code-block">{review.fixed_code}</pre>
                  </div>

                  <div className="review-section suggestion">
                    <div className="review-section-header">
                      <span className="dot purple"></span>
                      Suggestions
                    </div>
                    <div className="review-text">
                      <ReactMarkdown>{review.suggestions}</ReactMarkdown>
                    </div>
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
                </>
              )}
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}