const rawBaseUrl = import.meta.env.BASE_URL || "/";

const normalizeBasePath = (value: string) => {
  if (!value || value === "/" || value === "./" || value === ".") {
    return "";
  }

  let normalized = value.trim();
  if (!normalized.startsWith("/")) {
    normalized = `/${normalized}`;
  }
  if (normalized.length > 1 && normalized.endsWith("/")) {
    normalized = normalized.slice(0, -1);
  }
  return normalized;
};

export const routerBasePath = normalizeBasePath(rawBaseUrl);

export const appPath = (path: string) => {
  const normalizedPath = path.startsWith("./") ? path : `./${path}`;
  return routerBasePath ? `${routerBasePath}${normalizedPath}` : normalizedPath;
};

export const apiPath = (path: string) => appPath(`/api/${path.replace(/^\/+/, "")}`);

export const settingsPath = () => appPath("./settings");
export const apiSettingsPath = () => appPath("./api/settings");

export const stripBasePath = (path: string) => {
  if (!routerBasePath) {
    return path || "/";
  }
  if (path === routerBasePath) {
    return "/";
  }
  if (path.startsWith(`${routerBasePath}/`)) {
    return path.slice(routerBasePath.length) || "/";
  }
  return path || "/";
};
