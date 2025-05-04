export interface PredictionResponseDto {
    hour: number
    water_temperature: number
    air_temperature: number
    weather_condition: string
    prediction: number
    explanation: Record<string, number>  
}  