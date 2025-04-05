import express from 'express'
import cors from 'cors';  
import { getDownloadLink } from './utils/getTemperatureDataDownloadLink.js'
import { downloadAndExtractCSV } from './utils/downloadAndUnzipTemperatureData.js';
 import dotenv from 'dotenv';

dotenv.config(); // Load env vars from .env

const app = express();
const PORT = process.env.PORT || 3000;

app.use(cors()); 

app.get('/', (req, res) => {
  res.send('🌊 Eisbach Tracker backend listening!');
});

app.listen(PORT, () => {
  console.log(`✅ Server is running at http://localhost:${PORT}`);
});


app.get('/api/temperature', async (req, res) => {
    try {
      // 1. Get the download link via token
      const downloadUrl = await getDownloadLink();
  
      // 2. Download, unzip and parse CSV rows
      const rows = await downloadAndExtractCSV(downloadUrl);
  
      // 3. Pick latest data (last row)
      const latest = rows.at(-1);
      if (!latest || !latest.Mittelwert) {
        throw new Error('No temperature data found');
      }
  
      // 4. Convert e.g. "7,1" to 7.1
      const temperature = parseFloat(latest.Mittelwert.replace(',', '.'));
  
      res.json({ temperature });
    } catch (err) {
      console.error('❌ Server error:', err.stack || err.message);
      res.status(500).send('Internal Server Error');
    }
  });
