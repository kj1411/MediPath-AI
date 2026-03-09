const API_URL = 'http://localhost:8080';

export async function predictCondition(drugs) {
  const res = await fetch(`${API_URL}/predict`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ drugs }),
  });

  if (!res.ok) {
    let message = 'Something went wrong. Please try again.';
    try {
      const text = await res.text();
      if (text) message = text;
    } catch {
      // use default message
    }
    throw new Error(message);
  }

  return res.json();
}
