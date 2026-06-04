import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 50, // 50 Virtual Users
  duration: '30s', // Attack for 30 seconds
};

// Simple function to generate random string
function randomString(length = 10) {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < length; i++) {
    result += chars[Math.floor(Math.random() * chars.length)];
  }
  return result;
}

// Generate random URL
function randomUrl() {
  return `https://example.com/${randomString(8)}`;
}

export default function () {
  const originalUrl = randomUrl();
  
  // 1. Shorten the URL
  const shortenPayload = JSON.stringify({
    url: originalUrl,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const shortenRes = http.post('https://go-shorturl-mblo.onrender.com/shorten', shortenPayload, params);

  // Check if shorten worked
  const shortenSuccess = check(shortenRes, {
    'shorten: status is 200': (r) => r.status === 200,
    'shorten: has shortenedUrl': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.shortenedUrl !== undefined;
      } catch {
        return false;
      }
    }
  });

  // 2. If shorten was successful, try to retrieve the full URL
  if (shortenSuccess && shortenRes.status === 200) {
    try {
      const shortenBody = JSON.parse(shortenRes.body);
      const shortenedUrl = shortenBody.shortenedUrl;
      
      // Prepare payload for full URL retrieval
      const retrievePayload = JSON.stringify({
        url: shortenedUrl,
      });

      sleep(0.05); // Small delay between requests
      
      const retrieveRes = http.post('https://go-shorturl-mblo.onrender.com/full', retrievePayload, params);
      
      // Check if retrieve worked
      check(retrieveRes, {
        'retrieve: status is 200': (r) => r.status === 200,
        'retrieve: returns original URL': (r) => {
          try {
            const body = JSON.parse(r.body);
            return body.fullUrl === originalUrl;
          } catch {
            return false;
          }
        }
      });
      
    } catch (error) {
      console.log(`Error parsing shorten response: ${error}`);
    }
  }

  sleep(0.1); // Wait 100ms before next iteration
}