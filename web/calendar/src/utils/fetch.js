export default function fetchAPI(url, method = 'GET', body = '', signin = true) {
  const options = {
    method,
    headers: {
      'X-XSRF-TOKEN': 'csrf',
      'Content-Type': 'application/json; charset=UTF-8',
    },
  };

  if (signin) {
    options.credentials = 'same-origin';
  }

  if (body !== '') {
    options.body = body;
  }

  const prefix = '/api';

  return fetch(`${prefix}${url}`, options).catch((err) => {
    console.error(`Request failed: ${err}`);
    throw new Error(`Request failed: ${err}`);
  }).then((response) => {
    if (response.ok) {
      if (response.status === 204) {
        return null;
      }
      return response.json();
    }
    console.error(`Request failed: ${response.status}`);
    throw new Error(`Request failed: ${response.status}`);
  });
}
