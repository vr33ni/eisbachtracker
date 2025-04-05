import axios from 'axios';
import fs from 'fs';
import unzipper from 'unzipper';
import { parse } from 'csv-parse/sync';

export async function downloadAndExtractCSV(downloadUrl) {
  const zipPath = './data.zip';
  const directUrl = `${downloadUrl}&dl=1`;

  console.log('⏳ Checking if download is ready...');

  // 🔁 Poll until the file is ready (status 200 with content)
  let isReady = false;
  let attempts = 0;
  while (!isReady && attempts < 10) {
    const head = await axios.head(directUrl, { validateStatus: () => true });
    if (head.status === 200 && +head.headers['content-length'] > 0) {
      isReady = true;
      break;
    }
    attempts++;
    console.log(`⏳ Waiting... (${attempts})`);
    await new Promise(resolve => setTimeout(resolve, 3000));
  }

  if (!isReady) throw new Error('❌ Download not ready after 10 attempts');

  // ⬇️ Download zip file
  console.log('⬇️ Downloading zip from:', directUrl);
  const response = await axios.get(directUrl, {
    responseType: 'stream',
  });

  const writer = fs.createWriteStream(zipPath);
  response.data.pipe(writer);
  await new Promise((resolve, reject) => {
    writer.on('finish', resolve);
    writer.on('error', reject);
  });

  // 📦 Unzip and extract CSV
  console.log('📦 Unzipping...');
  const directory = await unzipper.Open.file(zipPath);

  const csvFile = directory.files.find(
    file => typeof file.path === 'string' && file.path.endsWith('.csv')
  );

  if (!csvFile) {
    console.error('❌ No CSV file found. Available files:', directory.files.map(f => f.path));
    throw new Error('No CSV file found in the zip');
  }

  const content = await csvFile.buffer();
  const csvText = content.toString('utf-8');

  console.log('🔍 CSV Preview:\n' + csvText.split('\n').slice(0, 10).join('\n'));

  // 🔍 Find the real header starting point
  const lines = csvText.split('\n');
  const headerIndex = lines.findIndex(line => line.trim().startsWith('Datum'));

  if (headerIndex === -1) {
    throw new Error('CSV header not found!');
  }

  const csvData = lines.slice(headerIndex).join('\n');

  console.log('📄 Parsing CSV...');
  const records = parse(csvData, {
    delimiter: ';',
    columns: true,
    skip_empty_lines: true,
  });

  console.log('✅ Parsed rows:', records.length);
  console.log('🔍 Sample row:', records[0]);

  // 🧹 Clean up temp zip
  fs.unlink(zipPath, err => {
    if (err) {
      console.warn('⚠️ Failed to delete ZIP file:', err.message);
    } else {
      console.log('🧹 Deleted ZIP file');
    }
  });

  return records;
}
