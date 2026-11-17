import { apiPath } from "../basePath.js";

export const resolveImageSrc = (src: string) => {
  if (src.startsWith("/api/")) {
    return apiPath(src.slice("/api/".length));
  }

  if (!src.startsWith("./")) {
    return src;
  }

  const parts = src.slice(2).split("/").filter((part) => part !== "");
  if (parts.length >= 3) {
    const [category, id, ...rest] = parts;
    return apiPath(`${category}/${id}/images/${rest.join("/")}`);
  }

  return apiPath(src.slice(2));
};
