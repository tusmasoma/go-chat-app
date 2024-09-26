// api.ts

export async function fetchInitialMessages({ sessionId }: { sessionId: string }) {
    const response = await fetch(`/api/messages?sessionId=${sessionId}`);
    const data = await response.json();
    return data.messages;
  }

  export async function fetchIsIndexed(url: string) {
    const response = await fetch(`/api/check-indexed?url=${encodeURIComponent(url)}`);
    const data = await response.json();
    return data.isIndexed;
  }

  export async function indexUrl(url: string) {
    await fetch(`/api/index-url`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url }),
    });
  }
