export type DaisyTheme = "catalog" | "catalog-dark";

export type ThemeOption = { label: string; description: string; value: DaisyTheme };

export const themes: ThemeOption[] = [
  { label: "Catalog Light", description: "Bright, calm, and focused", value: "catalog" },
  { label: "Catalog Dark", description: "Low-glare workspace", value: "catalog-dark" },
];
