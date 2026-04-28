const DEFAULT_GITHUB_REPO = 'oOyH/sublinkPro';

const normalizedEnvRepo = import.meta.env.VITE_GITHUB_REPO?.trim().replace(/^\/+|\/+$/g, '') || '';

export const GITHUB_REPO = normalizedEnvRepo || DEFAULT_GITHUB_REPO;
export const GITHUB_URL = `https://github.com/${GITHUB_REPO}`;
export const GITHUB_ISSUES_URL = `${GITHUB_URL}/issues`;
export const GITHUB_API_REPO = `https://api.github.com/repos/${GITHUB_REPO}`;
export const GITHUB_API_LATEST_RELEASE = `${GITHUB_API_REPO}/releases/latest`;
export const GITHUB_API_RELEASES = `${GITHUB_API_REPO}/releases?per_page=5`;

