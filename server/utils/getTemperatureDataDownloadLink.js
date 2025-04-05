import axios from 'axios';
import { wrapper } from 'axios-cookiejar-support';
import { CookieJar } from 'tough-cookie';
import qs from 'qs';
 
const jar = new CookieJar();
const client = wrapper(axios.create({ jar }));

export async function getDownloadLink() {
  const downloadPageUrl = 'https://www.gkd.bayern.de/de/fluesse/wassertemperatur/kelheim/muenchen-himmelreichbruecke-16515005/download';
  const enqueueUrl = 'https://www.gkd.bayern.de/de/downloadcenter/enqueue_download';

  // 🍪 Step 1: Get session cookies by visiting the page
  console.log('🍪 Visiting download page...');
  await client.get(downloadPageUrl, {
    headers: {
      'User-Agent': 'Mozilla/5.0',
      'Referer': downloadPageUrl,
    },
  });

  // 📨 Step 2: Prepare request data
  const data = {
    zr: 'monat',
    beginn: '01.04.2025',
    ende: '05.04.2025',
    email: 'test@test.de',
    geprueft: '0',
    wertart: 'tmw',
    f: '',
    t: JSON.stringify({
      '16515005': ['fluesse.wassertemperatur'],
    }),
  };

  const headers = {
    'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
    'User-Agent': 'Mozilla/5.0',
    'Referer': downloadPageUrl,
    'X-Requested-With': 'XMLHttpRequest',
    'Origin': 'https://www.gkd.bayern.de',
    'Accept': 'application/json, text/javascript, */*; q=0.01',
  };

  // 📨 Step 3: POST request to enqueue download
  console.log('📨 Sending request to enqueue_download...');
  const res = await client.post(enqueueUrl, qs.stringify(data), { headers });

  // ✅ Step 4: Handle response
  if (res.data?.result === 'success' && res.data?.deeplink) {
    const tokenMatch = res.data.deeplink.match(/token=([a-zA-Z0-9]+)/);
    if (!tokenMatch) throw new Error('❌ Token not found in deeplink');

    const token = tokenMatch[1];
    console.log('✅ Token:', token);

    const downloadUrl = `https://www.gkd.bayern.de/de/downloadcenter/download?token=${token}`;
    console.log('⬇️ Download URL:', downloadUrl);

    return downloadUrl;
  } else {
    console.error('❌ Failed to enqueue download:', res.data);
    throw new Error('Failed to get token');
  }
}
