export type PositionType =
  | "mumbai"
  | "bengaluru"
  | "delhi"
  | "allcities"
  | string;

export interface WeatherSummary {
  avg_temperature: number; // Average temperature
  city_name: string; // Name of the city
  condition_counts: number | null; // Counts of conditions (nullable)
  date: string; // Date of the record in ISO format
  dominant_condition: string; // Most common weather condition
  max_temperature: number; // Maximum temperature
  min_temperature: number; // Minimum temperature
  total_measurements: number; // Total number of measurements
}

export type Units = "celsius" | "fahrenheit" | "kelvin" | string;
