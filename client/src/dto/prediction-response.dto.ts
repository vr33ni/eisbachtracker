export interface PredictionResponseDto {
    hour: number
    water_temperature: number
    water_level: number
    air_temperature: number
    weather_condition: string
    prediction: number
    explanation: Record<string, number>  
}  