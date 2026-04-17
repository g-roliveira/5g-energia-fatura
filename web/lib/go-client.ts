export class GoApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public path: string,
  ) {
    super(message)
    this.name = 'GoApiError'
  }
}

export async function goFetch<T>(
  path: string,
  options?: RequestInit & { timeoutMs?: number },
): Promise<T> {
  const baseUrl = process.env.BACKEND_GO_URL
  if (!baseUrl) {
    throw new Error('BACKEND_GO_URL environment variable is not set')
  }

  const { timeoutMs = 30_000, ...fetchOptions } = options ?? {}
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs)

  const method = (fetchOptions.method ?? 'GET').toUpperCase()
  const headers: HeadersInit = {
    ...(method === 'POST' || method === 'PATCH'
      ? { 'Content-Type': 'application/json' }
      : {}),
    ...(fetchOptions.headers ?? {}),
  }

  try {
    const response = await fetch(`${baseUrl}${path}`, {
      ...fetchOptions,
      headers,
      signal: controller.signal,
    })

    if (!response.ok) {
      let message = response.statusText
      try {
        const body = await response.json()
        if (body?.error) message = body.error
      } catch {}
      throw new GoApiError(response.status, message, path)
    }

    return response.json() as Promise<T>
  } catch (err) {
    if (err instanceof GoApiError) throw err
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new GoApiError(504, 'Request timed out', path)
    }
    throw err
  } finally {
    clearTimeout(timeoutId)
  }
}
