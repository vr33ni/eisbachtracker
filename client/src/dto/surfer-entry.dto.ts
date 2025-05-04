export interface SurferEntryDto {
    timestamp: string
    count: number
    water_temperature: number // being passed to the server, cause takes longer than the other values to fetch; if null, it is being fetched again from the server 
  }
  