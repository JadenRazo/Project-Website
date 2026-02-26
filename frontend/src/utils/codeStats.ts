const GITHUB_README_URL = 'https://raw.githubusercontent.com/JadenRazo/JadenRazo/main/README.md';

function parseLOCFromReadme(text: string): number | null {
  const locMatch = text.match(/<!-- LOC_START -->\s*\n\*\*([0-9,]+)\*\* lines of code/);
  if (!locMatch) return null;
  return parseInt(locMatch[1].replace(/,/g, ''), 10) || null;
}

export async function fetchGitHubLOC(): Promise<{ totalLines: number }> {
  try {
    const res = await fetch(GITHUB_README_URL);
    if (res.ok) {
      const text = await res.text();
      const loc = parseLOCFromReadme(text);
      if (loc) return { totalLines: loc };
    }
  } catch {}

  try {
    const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
    const endpoint = apiUrl ? `${apiUrl}/api/v1/code/stats` : '/api/v1/code/stats';
    const res = await fetch(endpoint);
    if (res.ok) {
      const data = await res.json();
      if (typeof data.totalLines === 'number') return { totalLines: data.totalLines };
    }
  } catch {}

  try {
    const res = await fetch('/code_stats.json');
    if (res.ok) {
      const data = await res.json();
      if (typeof data.totalLines === 'number') return { totalLines: data.totalLines };
    }
  } catch {}

  throw new Error('Failed to fetch LOC from all sources');
}
