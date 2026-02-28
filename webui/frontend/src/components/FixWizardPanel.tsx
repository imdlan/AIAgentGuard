import { useState } from 'react';
import type { FixRequest, FixResponse } from '../types';
import { apiClient } from '../api/client';
import './FixWizardPanel.css';

export default function FixWizardPanel() {
  const [dryRun, setDryRun] = useState(true);
  const [auto, setAuto] = useState(false);
  const [category, setCategory] = useState('');
  const [response, setResponse] = useState<FixResponse | null>(null);
  const [loading, setLoading] = useState(false);

  const handleFix = async () => {
    setLoading(true);
    setResponse(null);
    try {
      const request: FixRequest = {
        dry_run: dryRun,
        auto: auto,
        category: category || undefined,
      };
      const result = await apiClient.fix(request);
      setResponse(result);
    } catch (err) {
      setResponse({
        success: false,
        message: 'Failed to execute fix',
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fix-wizard-panel">
      <div className="panel-header">
        <h3>🔧 Security Fix Wizard</h3>
      </div>

      <div className="fix-options">
        <div className="option">
          <label>
            <input
              type="checkbox"
              checked={dryRun}
              onChange={(e) => setDryRun(e.target.checked)}
            />
            <span>Dry Run (preview changes)</span>
          </label>
          <small>Check this to see what would be fixed without making changes</small>
        </div>

        <div className="option">
          <label>
            <input
              type="checkbox"
              checked={auto}
              onChange={(e) => setAuto(e.target.checked)}
              disabled={dryRun}
            />
            <span>Auto-Fix</span>
          </label>
          <small>Automatically execute fix commands (use with caution)</small>
        </div>

        <div className="option">
          <label>
            Category (optional):
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
            >
              <option value="">All Categories</option>
              <option value="filesystem">Filesystem</option>
              <option value="shell">Shell</option>
              <option value="network">Network</option>
              <option value="secrets">Secrets</option>
            </select>
          </label>
        </div>

        <button
          onClick={handleFix}
          disabled={loading}
          className="fix-button"
        >
          {loading ? '⏳ Executing...' : '🚀 Execute Fixes'}
        </button>
      </div>

      {response && (
        <div className="fix-response">
          <div className={`response-header ${response.success ? 'success' : 'error'}`}>
            {response.success ? '✅' : '❌'} {response.message}
          </div>

          {response.fixed && response.fixed.length > 0 && (
            <div className="fix-section fixed">
              <h4>✅ Fixed ({response.fixed.length})</h4>
              <ul>
                {response.fixed.map((item, index) => (
                  <li key={index}>{item}</li>
                ))}
              </ul>
            </div>
          )}

          {response.failed && response.failed.length > 0 && (
            <div className="fix-section failed">
              <h4>❌ Failed ({response.failed.length})</h4>
              <ul>
                {response.failed.map((item, index) => (
                  <li key={index}>{item}</li>
                ))}
              </ul>
            </div>
          )}

          {response.skipped && response.skipped.length > 0 && (
            <div className="fix-section skipped">
              <h4>⏭️ Skipped ({response.skipped.length})</h4>
              <ul>
                {response.skipped.map((item, index) => (
                  <li key={index}>{item}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
